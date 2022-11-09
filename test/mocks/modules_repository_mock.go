package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type ModulesRepositoryMock struct {
	mock.Mock
}

func (repositoryMock *ModulesRepositoryMock) Get(ctx context.Context, filter entities.ModuleFilter) ([]entities.Module, error) {
	args := repositoryMock.Mock.Called(ctx, filter)
	return args.Get(0).([]entities.Module), args.Error(1)
}

func (repositoryMock *ModulesRepositoryMock) Save(ctx context.Context, module *entities.Module) error {
	args := repositoryMock.Mock.Called(ctx, module)
	return args.Error(0)
}

func (repositoryMock *ModulesRepositoryMock) GetPaged(ctx context.Context, filter entities.ModuleFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	args := repositoryMock.Called(ctx, pagination, filter)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}

func (repositoryMock *ModulesRepositoryMock) Update(ctx context.Context, id string, module entities.Module) error {
	args := repositoryMock.Mock.Called(ctx, id, module)
	return args.Error(0)
}

func (repositoryMock *ModulesRepositoryMock) Delete(ctx context.Context, id string) error {
	args := repositoryMock.Mock.Called(ctx, id)
	return args.Error(0)
}
