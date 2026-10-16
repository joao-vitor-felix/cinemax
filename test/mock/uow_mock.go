package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UnitOfWorkMock struct {
	mock.Mock
}

func (m *UnitOfWorkMock) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if err := args.Error(0); err != nil {
		return err
	}
	return fn(ctx)
}
