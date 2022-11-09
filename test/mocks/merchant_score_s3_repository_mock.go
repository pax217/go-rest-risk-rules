package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type MerchantScoreS3RepositoryMock struct {
	mock.Mock
}

func (m *MerchantScoreS3RepositoryMock) GetFileContent(ctx context.Context, fileName string) ([]entities.MerchantScore, error) {
	args := m.Mock.Called(ctx, fileName)
	return args.Get(0).([]entities.MerchantScore), args.Error(1)
}
