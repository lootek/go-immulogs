package service

import (
	"context"

	"github.com/lootek/go-immulogs"
)

type REST struct {
}

func NewREST() *REST {
	return &REST{}
}

func (r REST) Start(context.Context) error {
	// TODO implement me
	panic("implement me")
	return nil
}

func (r REST) Stop() error {
	// TODO implement me
	panic("implement me")
	return nil
}

func (r REST) AddLog(immulogs.Entry) error {
	// TODO implement me
	panic("implement me")
	return nil
}

func (r REST) AddLogsBatch([]immulogs.Entry) error {
	// TODO implement me
	panic("implement me")
	return nil
}

func (r REST) LastN(immulogs.Bucket) ([]immulogs.Entry, error) {
	// TODO implement me
	panic("implement me")
	return nil, nil
}

func (r REST) Count(immulogs.Bucket) (int64, error) {
	// TODO implement me
	panic("implement me")
	return 0, nil
}
