package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type FamilyRepositoryMock struct {
	mock.Mock
}

func (m *FamilyRepositoryMock) AddFamily(ctx context.Context, family *entities.Family) error {
	args := m.Called(ctx, family)
	return args.Error(0)
}

func (m *FamilyRepositoryMock) ValidateExist(ctx context.Context, family entities.Family) (bool, error) {
	args := m.Called(ctx, family)
	return args.Get(0).(bool), args.Error(1)
}

func (m *FamilyRepositoryMock) Update(ctx context.Context, id string, family *entities.Family) error {
	args := m.Called(ctx, id, family)
	return args.Error(0)
}

func (m *FamilyRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *FamilyRepositoryMock) SearchPaged(ctx context.Context, pagination entities.Pagination,
	filter entities.FamilyFilter) (entities.PagedResponse, error) {
	args := m.Called(ctx, pagination, filter)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}

func (m *FamilyRepositoryMock) Search(ctx context.Context, filter entities.FamilyFilter) ([]entities.Family, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entities.Family), args.Error(1)
}
func (m *FamilyRepositoryMock) SearchEvaluate(ctx context.Context, filter entities.FamilyFilter) ([]entities.Family, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entities.Family), args.Error(1)
}
