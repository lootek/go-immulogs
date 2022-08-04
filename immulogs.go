package immulogs

import (
	"context"
)

type Service interface {
	Start(ctx context.Context) error
	Stop() error
}

type service struct {
	service Service
}

func NewService(io Service) *service {
	return &service{
		service: io,
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
