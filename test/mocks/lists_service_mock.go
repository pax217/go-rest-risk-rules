package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"

	"github.com/stretchr/testify/mock"
)

func NewListsServiceMock() ListsServiceMock {
	return ListsServiceMock{}
}

type ListsServiceMock struct {
	mock.Mock
}

func (m *ListsServiceMock) GetLists(ctx context.Context, listsSearch entities.ListsSearch) ([]entities.List, error) {
	args := m.Called(ctx, listsSearch)
	return args.Get(0).([]entities.List), args.Error(1)
}
