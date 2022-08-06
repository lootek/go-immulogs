package storage

import (
	"context"
	"testing"
	"time"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
	"github.com/stretchr/testify/require"
)

type item struct {
	k []byte
	v []byte
}

type immuMock struct {
	storage []item
}

func (i *immuMock) OpenSession(ctx context.Context, user []byte, pass []byte, database string) (err error) {
	return nil
}

func (i *immuMock) CloseSession(ctx context.Context) error {
	i.storage = nil
	return nil
}

func (i *immuMock) WithOptions(options *immudb.Options) ImmuClient {
	return i
}

func (i *immuMock) Set(ctx context.Context, key []byte, value []byte) (*schema.TxHeader, error) {
	i.storage = append(i.storage, item{key, value})
	return &schema.TxHeader{Nentries: 1}, nil
}

func (i *immuMock) Scan(ctx context.Context, req *schema.ScanRequest) (*schema.Entries, error) {
	var entries []*schema.Entry
	for _, i := range i.storage {
		entries = append(entries, &schema.Entry{
			Key:   i.k,
			Value: i.v,
		})
	}

	return &schema.Entries{Entries: entries}, nil
}

func (i *immuMock) SetAll(ctx context.Context, kvList *schema.SetRequest) (*schema.TxHeader, error) {
	for _, kv := range kvList.KVs {
		i.storage = append(i.storage, item{kv.Key, kv.Value})
	}

	return &schema.TxHeader{Nentries: int32(len(kvList.KVs))}, nil
}

func TestImmuDB(t *testing.T) {
	for testCase, bucketName := range map[string]string{
		"globally":   "",
		"per bucket": "/my-bucket-name",
	} {
		r := NewImmuDB(&immudb.Options{
			Username: "user",
			Password: "pass",
			Database: "db",
		})
		r.client = &immuMock{}

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
				require.Equal(t, map[string]any{"nentries": 1.}, got)
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
				require.Equal(t, map[string]any{"nentries": 3.}, got)
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

				require.NoError(t, err)
				require.Equal(t, []log.Entry{
					log.FromString("a sample log entry"),
					log.FromString("a sample log entry #1"),
					log.FromString("a sample log entry #2"),
					log.FromString("a sample log entry #3"),
				}, got)
			})
		})
	}
}
