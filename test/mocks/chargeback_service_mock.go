package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type ChargebackServiceMock struct {
	mock.Mock
}

func NewChargebackServiceMock() ChargebackServiceMock {
	return ChargebackServiceMock{}
}

func (m *ChargebackServiceMock) Save(ctx context.Context, payer entities.Payer) error {
	args := m.Mock.Called(ctx, payer)
	return args.Error(0)
}
