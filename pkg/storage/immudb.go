package storage

import (
	"github.com/lootek/go-immulogs"
)

type ImmuDB struct {
}

func NewImmuDB() *ImmuDB {
	return &ImmuDB{}
}

func (i ImmuDB) WriteOne(immulogs.Bucket, immulogs.Entry) error {
	// TODO implement me
	panic("implement me")
	return nil
}

func (i ImmuDB) WriteBatch(immulogs.Bucket, []immulogs.Entry) error {
	// TODO implement me
	panic("implement me")
	return nil
}

func (i ImmuDB) All(immulogs.Bucket) ([]immulogs.Entry, error) {
	// TODO implement me
	panic("implement me")
	return nil, nil
}

func (i ImmuDB) Last(immulogs.Bucket) ([]immulogs.Entry, error) {
	// TODO implement me
	panic("implement me")
	return nil, nil
}

func (i ImmuDB) Count(immulogs.Bucket) (int64, error) {
	// TODO implement me
	panic("implement me")
	return 0, nil
}
