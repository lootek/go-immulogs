package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
	"github.com/stretchr/testify/require"
)

type storageMock struct {
	entries []log.Entry
}

func (s *storageMock) Start(ctx context.Context) error {
	return nil
}

func (s *storageMock) Stop() error {
	return nil
}

func (s *storageMock) WriteOne(b bucket.Bucket, e log.Entry) (map[string]any, error) {
	s.entries = append(s.entries, e)
	return map[string]any{"written": 1}, nil
}

func (s *storageMock) WriteBatch(b bucket.Bucket, e []log.Entry) (map[string]any, error) {
	s.entries = append(s.entries, e...)
	return map[string]any{"written": len(e)}, nil
}

func (s *storageMock) All(b bucket.Bucket) ([]log.Entry, error) {
	return s.entries, nil
}

func (s *storageMock) Last(b bucket.Bucket, n uint64) ([]log.Entry, error) {
	return s.entries[uint64(len(s.entries))-n:], nil
}

func (s *storageMock) Count(b bucket.Bucket) (uint64, error) {
	return uint64(len(s.entries)), nil
}

func TestREST(t *testing.T) {
	for testCase, bucketName := range map[string]string{
		"globally":   "",
		"per bucket": "/my-bucket-name",
	} {
		r := NewREST(&storageMock{}, "localhost:8000", 10*time.Second)
		ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
		go func() {
			err := r.Start(ctx)
			require.NoError(t, err)
		}()
		defer r.Stop()
		defer cancelFn()

		t.Run(testCase, func(t *testing.T) {
			t.Run("count empty", func(t *testing.T) {
				req, _ := http.NewRequest("GET", fmt.Sprintf("%s/count", bucketName), nil)
				w := httptest.NewRecorder()
				r.srv.Handler.ServeHTTP(w, req)

				gotResponse, _ := ioutil.ReadAll(w.Body)
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, `{"count":0}`, string(gotResponse))
			})

			t.Run("add one", func(t *testing.T) {
				req, _ := http.NewRequest("POST", fmt.Sprintf("%s/add", bucketName), bytes.NewBufferString(`a sample log entry`))
				w := httptest.NewRecorder()
				r.srv.Handler.ServeHTTP(w, req)

				gotResponse, _ := ioutil.ReadAll(w.Body)
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, `{"written":1}`, string(gotResponse))
			})

			t.Run("count one", func(t *testing.T) {
				req, _ := http.NewRequest("GET", fmt.Sprintf("%s/count", bucketName), nil)
				w := httptest.NewRecorder()
				r.srv.Handler.ServeHTTP(w, req)

				gotResponse, _ := ioutil.ReadAll(w.Body)
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, `{"count":1}`, string(gotResponse))
			})

			t.Run("add batch", func(t *testing.T) {
				req, _ := http.NewRequest("POST", fmt.Sprintf("%s/batch", bucketName), bytes.NewBufferString(strings.Join([]string{
					`a sample log entry #1`,
					`a sample log entry #2`,
					`a sample log entry #3`,
				}, "\n")))
				w := httptest.NewRecorder()
				r.srv.Handler.ServeHTTP(w, req)

				gotResponse, _ := ioutil.ReadAll(w.Body)
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, `{"written":3}`, string(gotResponse))
			})

			t.Run("count all by now", func(t *testing.T) {
				req, _ := http.NewRequest("GET", fmt.Sprintf("%s/count", bucketName), nil)
				w := httptest.NewRecorder()
				r.srv.Handler.ServeHTTP(w, req)

				gotResponse, _ := ioutil.ReadAll(w.Body)
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, `{"count":4}`, string(gotResponse))
			})

			t.Run("get all", func(t *testing.T) {
				req, _ := http.NewRequest("GET", fmt.Sprintf("%s/last/-1", bucketName), nil)
				w := httptest.NewRecorder()
				r.srv.Handler.ServeHTTP(w, req)

				gotResponse, _ := ioutil.ReadAll(w.Body)
				require.Equal(t, http.StatusOK, w.Code)

				expectedJSON, _ := json.Marshal(map[string]any{
					"entries": []string{
						"a sample log entry",
						"a sample log entry #1",
						"a sample log entry #2",
						"a sample log entry #3",
					},
				})
				require.Equal(t, expectedJSON, gotResponse)
			})

			t.Run("get last 2", func(t *testing.T) {
				req, _ := http.NewRequest("GET", fmt.Sprintf("%s/last/2", bucketName), nil)
				w := httptest.NewRecorder()
				r.srv.Handler.ServeHTTP(w, req)

				gotResponse, _ := ioutil.ReadAll(w.Body)
				require.Equal(t, http.StatusOK, w.Code)

				expectedJSON, _ := json.Marshal(map[string]any{
					"entries": []string{
						"a sample log entry #2",
						"a sample log entry #3",
					},
				})
				require.Equal(t, expectedJSON, gotResponse)
			})
		})
	}
}
