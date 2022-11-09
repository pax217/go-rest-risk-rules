package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"

	"github.com/stretchr/testify/mock"
)

func NewFieldsServiceMock() FieldsServiceMock {
	return FieldsServiceMock{}
}

type FieldsServiceMock struct {
	mock.Mock
}

func (m *FieldsServiceMock) AddField(ctx context.Context, field entities.Field) error {
	args := m.Called(ctx, field)
	return args.Error(0)
}

func (m *FieldsServiceMock) Update(ctx context.Context, id string, field entities.Field) error {
	args := m.Mock.Called(ctx, id, field)
	return args.Error(0)
}

func (m *FieldsServiceMock) Delete(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}

func (m *FieldsServiceMock) GetFields(ctx context.Context, filter entities.FieldsFilter,
	pagination entities.Pagination) (interface{}, error) {
	args := m.Called(ctx, filter, pagination)
	return args.Get(0), args.Error(1)
}
