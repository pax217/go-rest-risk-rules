package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type FamilyCompaniesRepositoryMock struct {
	mock.Mock
}

func (m *FamilyCompaniesRepositoryMock) AddFamilyCompanies(ctx context.Context, familyCompanies *entities.FamilyCompanies) error {
	args := m.Called(ctx, familyCompanies)
	return args.Error(0)
}

func (m *FamilyCompaniesRepositoryMock) GetFamilyCompanies(
	ctx context.Context,
	familyCompaniesFilter entities.FamilyCompaniesFilter) ([]entities.FamilyCompanies, error) {
	args := m.Called(ctx, familyCompaniesFilter)
	return args.Get(0).([]entities.FamilyCompanies), args.Error(1)
}

func (m *FamilyCompaniesRepositoryMock) SearchPaged(ctx context.Context,
	pag entities.Pagination,
	fil entities.FamilyCompaniesFilter) (entities.PagedResponse, error) {
	args := m.Called(ctx, pag, fil)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}

func (m *FamilyCompaniesRepositoryMock) Search(ctx context.Context,
	filter entities.FamilyCompaniesFilter) ([]entities.FamilyCompanies, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entities.FamilyCompanies), args.Error(1)
}

func (m *FamilyCompaniesRepositoryMock) Update(
	ctx context.Context,
	id string,
	familyCompanies *entities.FamilyCompanies) error {
	args := m.Called(ctx, id, familyCompanies)
	return args.Error(0)
}

func (m *FamilyCompaniesRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
