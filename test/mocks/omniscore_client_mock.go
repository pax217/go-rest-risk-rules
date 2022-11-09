package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type OmniscoreClientMock struct {
	mock.Mock
}

func (m *OmniscoreClientMock) GetScore(ctx context.Context, charge entities.ChargeRequest) (float64, error) {
	args := m.Mock.Called(ctx, charge)
	return args.Get(0).(float64), args.Error(1)
}
