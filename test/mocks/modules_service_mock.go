package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type ModuleServiceMock struct {
	mock.Mock
}

func (m *ModuleServiceMock) Add(ctx context.Context, module entities.Module) error {
	args := m.Mock.Called(ctx, module)
	return args.Error(0)
}

func (m *ModuleServiceMock) GetAll(ctx context.Context, pagination entities.Pagination,
	filter entities.ModuleFilter) (interface{}, error) {
	args := m.Mock.Called(ctx, pagination, filter)
	return args.Get(0), args.Error(1)
}

func (m *ModuleServiceMock) Update(ctx context.Context, id string, module entities.Module) error {
	args := m.Mock.Called(ctx, id, module)
	return args.Error(0)
}

func (m *ModuleServiceMock) Delete(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}
