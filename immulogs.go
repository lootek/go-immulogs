package immulogs

import (
	"context"
)

type Service interface {
	Start(ctx context.Context) error
	Stop() error

	AddLog(e Entry) error
	AddLogsBatch(e []Entry) error

	LastN(b Bucket) ([]Entry, error)
	Count(b Bucket) (int64, error)
}

type Storage interface {
	WriteOne(b Bucket, e Entry) error
	WriteBatch(b Bucket, e []Entry) error

	All(b Bucket) ([]Entry, error)
	Last(b Bucket) ([]Entry, error)
	Count(b Bucket) (int64, error)
}

type Entry interface {
	String() string
}

type Bucket interface {
	String() string
}

type service struct {
	service Service
	storage Storage
}

func NewService(io Service, s Storage) *service {
	return &service{
		service: io,
		storage: s,
	}
}

func (s service) Run(ctx context.Context) error {
	errCh := make(chan error)
	go func() {
		if err := s.service.Start(ctx); err != nil {
			errCh <- err
			return
		}
	}()

	select {
	case <-ctx.Done():
	case <-errCh:
	}

	return nil
}

func (s service) Stop() error {
	if err := s.service.Stop(); err != nil {
		return err
	}

	return nil
}
