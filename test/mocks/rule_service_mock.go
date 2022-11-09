package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type RuleServiceMock struct {
	mock.Mock
}

func (m *RuleServiceMock) BuildRule(ruleContent []entities.RuleContent) string {
	args := m.Mock.Called(ruleContent)
	return args.Get(0).(string)
}

func (m *RuleServiceMock) AddRule(ctx context.Context, rule entities.Rule) (entities.Rule, error) {
	args := m.Mock.Called(ctx, rule)
	return args.Get(0).(entities.Rule), args.Error(1)
}

func (m *RuleServiceMock) UpdateRule(ctx context.Context, ruleID string, ruleReq entities.Rule) error {
	args := m.Mock.Called(ctx, ruleID, ruleReq)
	return args.Error(0)
}

func (m *RuleServiceMock) RemoveRule(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}

func (m *RuleServiceMock) ListRules(ctx context.Context, ruleFilter entities.RuleFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	args := m.Mock.Called(ctx, ruleFilter, pagination)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}
