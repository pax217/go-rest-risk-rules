package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type OperatorServiceMock struct {
	mock.Mock
}

func (m *OperatorServiceMock) AddOperator(ctx context.Context, operator entities.Operator) error {
	args := m.Mock.Called(ctx, operator)
	return args.Error(0)
}

func (m *OperatorServiceMock) Get(ctx context.Context, filter entities.OperatorFilter,
	pagination entities.Pagination) (interface{}, error) {
	args := m.Mock.Called(ctx, filter, pagination)
	return args.Get(0), args.Error(1)
}

func (m *OperatorServiceMock) Update(ctx context.Context, id string, operator entities.Operator) error {
	args := m.Mock.Called(ctx, id, operator)
	return args.Error(0)
}

func (m *OperatorServiceMock) Delete(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}
