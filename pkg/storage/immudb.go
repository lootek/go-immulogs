package storage

import (
	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
)

type ImmuDB struct {
}

func NewImmuDB() *ImmuDB {
	return &ImmuDB{}
}

func (m ImmuDB) WriteOne(b bucket.Bucket, e log.Entry) error {
	// TODO implement me
	return nil
}

func (m ImmuDB) WriteBatch(b bucket.Bucket, e []log.Entry) error {
	// TODO implement me
	return nil
}

func (m ImmuDB) All(b bucket.Bucket) ([]log.Entry, error) {
	// TODO implement me
	return nil, nil
}

func (m ImmuDB) Last(b bucket.Bucket, n uint64) ([]log.Entry, error) {
	// TODO implement me
	return nil, nil
}

func (m ImmuDB) Count(b bucket.Bucket) (uint64, error) {
	// TODO implement me
	return 0, nil
}
