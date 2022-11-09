package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type ConditionServiceMock struct {
	mock.Mock
}

func (m *ConditionServiceMock) Add(ctx context.Context, condition entities.Condition) error {
	args := m.Mock.Called(ctx, condition)
	return args.Error(0)
}
func (m *ConditionServiceMock) FindByName(ctx context.Context, condition entities.Condition) (
	entities.Condition, error) {
	args := m.Mock.Called(condition, ctx)
	return args.Get(0).(entities.Condition), args.Error(1)
}

func (m *ConditionServiceMock) GetAll(ctx context.Context, filter entities.ConditionsFilter, pagination entities.Pagination) (
	entities.PagedResponse, error) {
	args := m.Mock.Called(ctx, filter, pagination)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}

func (m *ConditionServiceMock) Update(ctx context.Context, id string, condition entities.Condition) error {
	args := m.Mock.Called(ctx, id, condition)
	return args.Error(0)
}

func (m *ConditionServiceMock) Delete(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}
