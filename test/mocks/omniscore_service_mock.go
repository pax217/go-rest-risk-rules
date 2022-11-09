package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type OmniscoreServiceMock struct {
	mock.Mock
}

func (m *OmniscoreServiceMock) GetScore(ctx context.Context, charge entities.ChargeRequest) float64 {
	args := m.Mock.Called(ctx, charge)
	return args.Get(0).(float64)
}
