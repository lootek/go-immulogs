package immulogs

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) Start(ctx context.Context) error {
	s.Called("Start")
	return nil
}

func (s *serviceMock) Stop() error {
	s.Called("Stop")
	return nil
}

func TestService(t *testing.T) {
	t.Run("regular scenario", func(t *testing.T) {
		srvMock := serviceMock{}
		srvMock.On("Start", mock.Anything).Times(2).Return(nil)
		srvMock.On("Stop", mock.Anything).Times(2).Return(nil)

		service := NewService(&srvMock, &srvMock)
		ctx, cancelFn := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancelFn()

		go func() {
			time.Sleep(time.Millisecond * 250)
			err := service.Stop()
			require.NoError(t, err)
		}()

		err := service.Run(ctx)
		require.Error(t, err)
		require.Equal(t, context.Canceled, err)

		require.True(t, srvMock.AssertExpectations(t))
	})

	t.Run("timeout scenario", func(t *testing.T) {
		srvMock := serviceMock{}
		srvMock.On("Start", mock.Anything).Times(2).Return(nil)
		// srvMock.On("Stop", mock.Anything).Times(0).Return(nil)

		service := NewService(&srvMock, &srvMock)
		ctx, cancelFn := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancelFn()

		err := service.Run(ctx)
		require.Error(t, err)
		require.Equal(t, context.DeadlineExceeded, err)

		require.True(t, srvMock.AssertExpectations(t))
	})
}
