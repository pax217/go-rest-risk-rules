package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type RkListsRestClient struct {
	mock.Mock
}

func (m *RkListsRestClient) ListsSearch(ctx context.Context, listsSearch entities.ListsSearch) ([]entities.List, error) {
	args := m.Mock.Called(ctx, listsSearch)
	return args.Get(0).([]entities.List), args.Error(1)
}
