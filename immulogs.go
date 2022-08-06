package immulogs

import (
	"context"
)

type Service interface {
	Start(ctx context.Context) error
	Stop() error
}

type service struct {
	storageService Service
	ioService      Service

	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewService(storage Service, io Service) *service {
	return &service{
		storageService: storage,
		ioService:      io,
	}
}

func (s *service) Run(ctx context.Context) error {
	ctx, cancelFn := context.WithCancel(ctx)
	s.ctx = ctx
	s.cancelFn = cancelFn

	errCh := make(chan error)
	go func() {
		if err := s.storageService.Start(ctx); err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		if err := s.ioService.Start(ctx); err != nil {
			errCh <- err
			return
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (s *service) Stop() error {
	if err := s.ioService.Stop(); err != nil {
		return err
	}

	if err := s.storageService.Stop(); err != nil {
		return err
	}

	s.cancelFn()

	return nil
}
