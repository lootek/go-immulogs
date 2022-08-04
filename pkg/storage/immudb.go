package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
)

const (
	defaultTimeout = 10 * time.Second
)

type ImmuDB struct {
	ctx      context.Context
	cancelFn context.CancelFunc

	client immudb.ImmuClient
	opts   *immudb.Options
}

func NewImmuDB(opts *immudb.Options) *ImmuDB {
	return &ImmuDB{
		client: immudb.NewClient().WithOptions(opts),
		opts:   opts,
	}
}

func (i *ImmuDB) Start(ctx context.Context) error {
	ctx, cancelFn := context.WithCancel(context.Background())
	i.ctx = ctx
	i.cancelFn = cancelFn

	err := i.client.OpenSession(ctx, []byte(i.opts.Username), []byte(i.opts.Password), i.opts.Database)
	if err != nil {
		return err
	}

	return nil
}

func (i *ImmuDB) Stop() error {
	i.cancelFn()
	return nil
}

func (i *ImmuDB) WriteOne(b bucket.Bucket, e log.Entry) (map[string]any, error) {
	ctx, cancelFn := context.WithTimeout(i.ctx, defaultTimeout)
	defer cancelFn()

	tx, err := i.client.Set(ctx, b.Bytes(), e.Bytes())
	if err != nil {
		return nil, err
	}

	var resp map[string]any
	txJSON, _ := json.Marshal(tx)
	_ = json.Unmarshal(txJSON, &resp)

	return resp, nil
}

func (i *ImmuDB) WriteBatch(b bucket.Bucket, e []log.Entry) (map[string]any, error) {
	ctx, cancelFn := context.WithTimeout(i.ctx, defaultTimeout)
	defer cancelFn()

	var KVs []*schema.KeyValue
	for _, entry := range e {
		KVs = append(KVs, &schema.KeyValue{Key: b.Bytes(), Value: entry.Bytes()})
	}

	tx, err := i.client.SetAll(ctx, &schema.SetRequest{KVs: KVs})
	if err != nil {
		return nil, err
	}

	var resp map[string]any
	txJSON, _ := json.Marshal(tx)
	_ = json.Unmarshal(txJSON, &resp)

	return resp, nil
}

func (i *ImmuDB) All(b bucket.Bucket) ([]log.Entry, error) {
	return i.Last(b, 0)
}

func (i *ImmuDB) Last(b bucket.Bucket, n uint64) ([]log.Entry, error) {
	ctx, cancelFn := context.WithTimeout(i.ctx, defaultTimeout)
	defer cancelFn()

	history, err := i.client.History(ctx, &schema.HistoryRequest{
		Key:   []byte(b.String()),
		Limit: int32(n),
		Desc:  false,
	})
	if err != nil {
		return nil, err
	}

	var entries []log.Entry
	for _, e := range history.Entries {
		entries = append(entries, log.FromBytes(e.Value))
	}

	return entries, nil
}

func (i *ImmuDB) Count(b bucket.Bucket) (uint64, error) {
	ctx, cancelFn := context.WithTimeout(i.ctx, defaultTimeout)
	defer cancelFn()

	// TODO: Add buckets support
	count, err := i.client.CountAll(ctx)

	return count.GetCount(), err
}
