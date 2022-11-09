package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type ChargeServiceMock struct {
	mock.Mock
}

func (m *ChargeServiceMock) EvaluateCharge(ctx context.Context,
	charge entities.ChargeRequest) (entities.EvaluationResponse, error) {
	args := m.Called(ctx, charge)
	return args.Get(0).(entities.EvaluationResponse), args.Error(1)
}

func (m *ChargeServiceMock) EvaluateChargeOnlyRules(ctx context.Context,
	charge entities.ChargeRequest) (entities.RulesEvaluationResponse, error) {
	args := m.Called(ctx, charge)
	return args.Get(0).(entities.RulesEvaluationResponse), args.Error(1)
}

func (m *ChargeServiceMock) Get(ctx context.Context, id string) (entities.EvaluationResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.EvaluationResponse), args.Error(1)
}

func (m *ChargeServiceMock) GetOnlyRules(ctx context.Context, id string) (entities.RulesEvaluationResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.RulesEvaluationResponse), args.Error(1)
}
