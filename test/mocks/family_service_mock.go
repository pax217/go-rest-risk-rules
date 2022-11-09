package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type FamilyServiceMock struct {
	mock.Mock
}

func NewFamilyServiceMock() FamilyServiceMock {
	return FamilyServiceMock{}
}

func (m *FamilyServiceMock) Create(ctx context.Context, family entities.Family) error {
	args := m.Mock.Called(ctx, family)
	return args.Error(0)
}

func (m *FamilyServiceMock) Update(ctx context.Context, id string, family entities.Family) error {
	args := m.Mock.Called(ctx, id, family)
	return args.Error(0)
}

func (m *FamilyServiceMock) GetFamilyFromFilter(ctx context.Context,
	filter entities.FamilyFilter) (entities.Family, error) {
	args := m.Mock.Called(ctx, filter)
	return args.Get(0).(entities.Family), args.Error(1)
}
func (m *FamilyServiceMock) GetFamily(ctx context.Context,
	filter entities.FamilyFilter) (entities.Family, error) {
	args := m.Mock.Called(ctx, filter)
	return args.Get(0).(entities.Family), args.Error(1)
}
func (m *FamilyServiceMock) Get(ctx context.Context, pagination entities.Pagination,
	filter entities.FamilyFilter) (interface{}, error) {
	args := m.Mock.Called(ctx, pagination, filter)
	return args.Get(0), args.Error(1)
}

func (m *FamilyServiceMock) Delete(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}
