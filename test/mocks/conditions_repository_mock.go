package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type ConditionsRepositoryMock struct {
	mock.Mock
}

func (repository *ConditionsRepositoryMock) Add(ctx context.Context, condition *entities.Condition) error {
	args := repository.Mock.Called(ctx, condition)
	return args.Error(0)
}

func (repository *ConditionsRepositoryMock) FindByName(ctx context.Context, condition entities.Condition) (
	entities.Condition, error) {
	args := repository.Mock.Called(ctx, condition)
	return args.Get(0).(entities.Condition), args.Error(1)
}

func (repository *ConditionsRepositoryMock) GetAll(
	ctx context.Context,
	filter entities.ConditionsFilter,
	pagination entities.Pagination) (
	entities.PagedResponse, error) {
	args := repository.Mock.Called(ctx, filter, pagination)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}

func (repository *ConditionsRepositoryMock) Update(ctx context.Context, id string, condition entities.Condition) error {
	args := repository.Mock.Called(ctx, id, condition)
	return args.Error(0)
}

func (repository *ConditionsRepositoryMock) Delete(ctx context.Context, id string) error {
	args := repository.Mock.Called(ctx, id)
	return args.Error(0)
}
