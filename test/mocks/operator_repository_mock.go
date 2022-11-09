package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type OperatorRepositoryMock struct {
	mock.Mock
}

func (m *OperatorRepositoryMock) Save(ctx context.Context, operator *entities.Operator) error {
	args := m.Called(ctx, operator)
	return args.Error(0)
}
func (m *OperatorRepositoryMock) Get(ctx context.Context, operatorFilter entities.OperatorFilter) (
	[]entities.Operator, error) {
	args := m.Called(ctx, operatorFilter)
	return args.Get(0).([]entities.Operator), args.Error(1)
}
func (m *OperatorRepositoryMock) GetPaged(ctx context.Context, filter entities.OperatorFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	args := m.Called(ctx, filter, pagination)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}

func (m *OperatorRepositoryMock) Update(ctx context.Context, id string, operator entities.Operator) error {
	args := m.Mock.Called(ctx, id, operator)
	return args.Error(0)
}

func (m *OperatorRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}
