package storage

import (
	"context"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
)

type ImmuDB struct {
	ctx      context.Context
	cancelFn context.CancelFunc

	client immudb.ImmuClient
	opts   *immudb.Options
}

func NewImmuDB() *ImmuDB {
	ctx, cancelFn := context.WithCancel(context.Background())

	opts := immudb.DefaultOptions().
		WithAddress("localhost").
		WithPort(3322).
		WithUsername("user").
		WithPassword("pass").
		WithDatabase("logs")

	client := immudb.NewClient().WithOptions(opts)

	return &ImmuDB{
		ctx:      ctx,
		cancelFn: cancelFn,

		client: client,
		opts:   opts,
	}
}

func (i *ImmuDB) Start() error {
	err := i.client.OpenSession(i.ctx, []byte(i.opts.Username), []byte(i.opts.Password), i.opts.Database)
	if err != nil {
		return err
	}

	return nil
}

func (i *ImmuDB) Stop() error {
	i.cancelFn()
	return nil
}

func (i *ImmuDB) WriteOne(b bucket.Bucket, e log.Entry) error {
	// TODO implement me
	return nil
}

func (i *ImmuDB) WriteBatch(b bucket.Bucket, e []log.Entry) error {
	// TODO implement me
	return nil
}

func (i *ImmuDB) All(b bucket.Bucket) ([]log.Entry, error) {
	// TODO implement me
	return nil, nil
}

func (i *ImmuDB) Last(b bucket.Bucket, n uint64) ([]log.Entry, error) {
	// TODO implement me
	return nil, nil
}

func (i *ImmuDB) Count(b bucket.Bucket) (uint64, error) {
	// TODO implement me
	return 0, nil
}
