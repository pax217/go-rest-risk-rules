package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type FamilyCompaniesServiceMock struct {
	mock.Mock
}

func NewFamilyCompaniesServiceMock() FamilyCompaniesServiceMock {
	return FamilyCompaniesServiceMock{}
}

func (m *FamilyCompaniesServiceMock) Create(ctx context.Context, familyCompanies entities.FamilyCompanies) error {
	args := m.Mock.Called(ctx, familyCompanies)
	return args.Error(0)
}

func (m *FamilyCompaniesServiceMock) Update(ctx context.Context, id string,
	familyCompanies entities.FamilyCompanies) error {
	args := m.Mock.Called(ctx, id, familyCompanies)
	return args.Error(0)
}

func (m *FamilyCompaniesServiceMock) Get(ctx context.Context, pagination entities.Pagination,
	filter entities.FamilyCompaniesFilter) (interface{}, error) {
	args := m.Mock.Called(ctx, pagination, filter)
	return args.Get(0), args.Error(1)
}

func (m *FamilyCompaniesServiceMock) GetFamiliesCompaniesFromFilter(ctx context.Context,
	filter entities.FamilyCompaniesFilter) ([]entities.FamilyCompanies, error) {
	args := m.Mock.Called(ctx, filter)
	return args.Get(0).([]entities.FamilyCompanies), args.Error(1)
}

func (m *FamilyCompaniesServiceMock) Delete(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}
