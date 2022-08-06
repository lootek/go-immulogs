package storage

import (
	"context"
	"testing"
	"time"

	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	for testCase, bucketName := range map[string]string{
		"globally":   "",
		"per bucket": "/my-bucket-name",
	} {
		r := NewMemory()
		ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
		go func() {
			err := r.Start(ctx)
			require.NoError(t, err)
		}()
		defer r.Stop()
		defer cancelFn()

		t.Run(testCase, func(t *testing.T) {
			t.Run("count empty", func(t *testing.T) {
				got, err := r.Count(bucket.NewBucket(bucketName))
				require.NoError(t, err)
				require.Equal(t, uint64(0), got)
			})

			t.Run("add one", func(t *testing.T) {
				got, err := r.WriteOne(bucket.NewBucket(bucketName), log.FromString(`a sample log entry`))
				require.NoError(t, err)
				require.Equal(t, map[string]any{"written": 1}, got)
			})

			t.Run("count one", func(t *testing.T) {
				got, err := r.Count(bucket.NewBucket(bucketName))
				require.NoError(t, err)
				require.Equal(t, uint64(1), got)
			})

			t.Run("add batch", func(t *testing.T) {
				got, err := r.WriteBatch(bucket.NewBucket(bucketName), []log.Entry{
					log.FromString(`a sample log entry #1`),
					log.FromString(`a sample log entry #2`),
					log.FromString(`a sample log entry #3`),
				})
				require.NoError(t, err)
				require.Equal(t, map[string]any{"written": 3}, got)
			})

			t.Run("count all by now", func(t *testing.T) {
				got, err := r.Count(bucket.NewBucket(bucketName))
				require.NoError(t, err)
				require.Equal(t, uint64(4), got)
			})

			t.Run("get all", func(t *testing.T) {
				got, err := r.All(bucket.NewBucket(bucketName))
				require.NoError(t, err)
				require.Equal(t, []log.Entry{
					log.FromString("a sample log entry"),
					log.FromString("a sample log entry #1"),
					log.FromString("a sample log entry #2"),
					log.FromString("a sample log entry #3"),
				}, got)
			})

			t.Run("get last 2", func(t *testing.T) {
				got, err := r.Last(bucket.NewBucket(bucketName), 2)

				if bucketName != "" {
					require.NoError(t, err)
					require.Equal(t, []log.Entry{
						log.FromString("a sample log entry #2"),
						log.FromString("a sample log entry #3"),
					}, got)
				} else {
					require.Error(t, err)
					require.Equal(t, "empty bucket not supported in this context with in-memory storage", err.Error())
					require.Equal(t, []log.Entry(nil), got)
				}
			})
		})
	}
}
