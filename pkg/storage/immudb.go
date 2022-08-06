package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/google/uuid"
	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
)

const (
	defaultTimeout = 3 * time.Second
)

type ImmuDB struct {
	ctx      context.Context
	cancelFn context.CancelFunc

	client ImmuClient
	opts   *immudb.Options
}

func NewImmuDB(opts *immudb.Options) *ImmuDB {
	return &ImmuDB{
		client: immuClientWrapper{immudb.NewClient().WithOptions(opts)},
		opts:   opts,
	}
}

// immuClientWrapper is a hack to make ImmuClient possible to implement locally (and mockable for tests)
// immudb.ImmuClient is impossible to implement here directly as it relies on an unexported type immudb.*immuClient
type immuClientWrapper struct {
	immudb.ImmuClient
}

func (i immuClientWrapper) WithOptions(options *immudb.Options) ImmuClient {
	i.ImmuClient.WithOptions(options)
	return i
}

// ImmuClient is a subset of insanely huge immudb.ImmuClient
// it contains only the functions we really need
type ImmuClient interface {
	OpenSession(ctx context.Context, user []byte, pass []byte, database string) (err error)
	CloseSession(ctx context.Context) error
	WithOptions(options *immudb.Options) ImmuClient
	Set(ctx context.Context, key []byte, value []byte) (*schema.TxHeader, error)
	Scan(ctx context.Context, req *schema.ScanRequest) (*schema.Entries, error)
	SetAll(ctx context.Context, kvList *schema.SetRequest) (*schema.TxHeader, error)
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

	tx, err := i.client.Set(ctx, i.key(b), e.Bytes())
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
		KVs = append(KVs, &schema.KeyValue{Key: i.key(b), Value: entry.Bytes()})
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

	scanned, err := i.client.Scan(ctx, &schema.ScanRequest{
		Prefix: b.Bytes(),
		Desc:   false,
		Limit:  n,
	})
	if err != nil {
		return nil, err
	}

	var entries []log.Entry
	for _, e := range scanned.Entries {
		entries = append(entries, log.FromBytes(e.Value))
	}

	return entries, nil
}

func (i *ImmuDB) Count(b bucket.Bucket) (uint64, error) {
	entries, err := i.Last(b, 0)
	return uint64(len(entries)), err
}

func (i *ImmuDB) key(b bucket.Bucket) []byte {
	return []byte(fmt.Sprintf("%s_%s", b.String(), uuid.NewString()))
}
