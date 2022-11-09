package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MerchantsScoreServiceMock struct {
	mock.Mock
}

func (m *MerchantsScoreServiceMock) MerchantScoreProcessing(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
