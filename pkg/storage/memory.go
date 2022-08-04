package storage

import (
	"github.com/lootek/go-immulogs"
)

type Memory struct {
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m Memory) WriteOne(immulogs.Bucket, immulogs.Entry) error {
	// TODO implement me
	panic("implement me")
	return nil
}

func (m Memory) WriteBatch(immulogs.Bucket, []immulogs.Entry) error {
	// TODO implement me
	panic("implement me")
	return nil
}

func (m Memory) All(immulogs.Bucket) ([]immulogs.Entry, error) {
	// TODO implement me
	panic("implement me")
	return nil, nil
}

func (m Memory) Last(immulogs.Bucket) ([]immulogs.Entry, error) {
	// TODO implement me
	panic("implement me")
	return nil, nil
}

func (m Memory) Count(immulogs.Bucket) (int64, error) {
	// TODO implement me
	panic("implement me")
	return 0, nil
}
