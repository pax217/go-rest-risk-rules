package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type FieldsRepositoryMock struct {
	mock.Mock
}

func (m *FieldsRepositoryMock) AddField(ctx context.Context, field *entities.Field) error {
	args := m.Called(ctx, field)
	return args.Error(0)
}

func (m *FieldsRepositoryMock) Update(ctx context.Context, id string, field *entities.Field) error {
	args := m.Mock.Called(ctx, id, field)
	return args.Error(0)
}

func (m *FieldsRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}

func (m *FieldsRepositoryMock) GetFields(ctx context.Context, filter entities.FieldsFilter) ([]entities.Field, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entities.Field), args.Error(1)
}

func (m *FieldsRepositoryMock) GetFieldsPaged(ctx context.Context, filter entities.FieldsFilter,
	page entities.Pagination) (entities.PagedResponse, error) {
	args := m.Called(ctx, filter, page)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}
