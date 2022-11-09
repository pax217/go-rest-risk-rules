package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type RulesRepositoryMock struct {
	mock.Mock
}

func (m *RulesRepositoryMock) GetRulesByFilters(ctx context.Context, filter entities.RuleFilter,
	component entities.ConsoleComponent) ([]entities.Rule, error) {
	args := m.Called(ctx, filter, component)
	return args.Get(0).([]entities.Rule), args.Error(1)
}

func (m *RulesRepositoryMock) FindRulesPaged(ctx context.Context, filter entities.RuleFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	args := m.Called(ctx, filter, pagination)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}

func (m *RulesRepositoryMock) AddRule(ctx context.Context, request entities.Rule) (entities.Rule, error) {
	args := m.Called(request, ctx)
	return args.Get(0).(entities.Rule), args.Error(1)
}

func (m *RulesRepositoryMock) UpdateRule(ctx context.Context, ruleID string, ruleReq entities.Rule) error {
	args := m.Called(ctx, ruleID, ruleReq)
	return args.Error(0)
}

func (m *RulesRepositoryMock) RemoveRule(ctx context.Context, ruleID string) error {
	args := m.Called(ctx, ruleID)
	return args.Error(0)
}
