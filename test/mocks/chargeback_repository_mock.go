package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type ChargebackRepositoryMock struct {
	mock.Mock
}

func (m *ChargebackRepositoryMock) Save(ctx context.Context, payer entities.Payer) error {
	args := m.Called(ctx, payer)
	return args.Error(0)
}

func (m *ChargebackRepositoryMock) Update(ctx context.Context, payer entities.Payer) error {
	args := m.Called(ctx, payer)
	return args.Error(0)
}

func (m *ChargebackRepositoryMock) Find(ctx context.Context, payer entities.Payer) (entities.Payer, error) {
	args := m.Called(ctx, payer)
	return args.Get(0).(entities.Payer), args.Error(1)
}
