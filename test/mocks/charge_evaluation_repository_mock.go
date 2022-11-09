package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type ChargeEvaluationRepositoryMock struct {
	mock.Mock
}

func (m *ChargeEvaluationRepositoryMock) Save(ctx context.Context, evaluation entities.EvaluationResponse) error {
	args := m.Called(ctx, evaluation)
	return args.Error(0)
}

func (m *ChargeEvaluationRepositoryMock) SaveOnlyRules(ctx context.Context,
	evaluation entities.RulesEvaluationResponse) error {
	args := m.Called(ctx, evaluation)
	return args.Error(0)
}

func (m *ChargeEvaluationRepositoryMock) Get(ctx context.Context, id string) (entities.EvaluationResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.EvaluationResponse), args.Error(1)
}

func (m *ChargeEvaluationRepositoryMock) GetOnlyRules(ctx context.Context, id string) (entities.RulesEvaluationResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.RulesEvaluationResponse), args.Error(1)
}
