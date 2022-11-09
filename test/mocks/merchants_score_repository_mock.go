package mocks

import (
	"context"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/mock"
)

type MerchantsScoreRepositoryMock struct {
	mock.Mock
}

func (r *MerchantsScoreRepositoryMock) WriteMerchantsScore(ctx context.Context, merchant []entities.MerchantScore) error {
	args := r.Mock.Called(ctx, merchant)
	return args.Error(0)
}

func (r *MerchantsScoreRepositoryMock) FindByMerchantID(ctx context.Context, companyID string) (entities.MerchantScore, error) {
	args := r.Called(ctx, companyID)
	return args.Get(0).(entities.MerchantScore), args.Error(1)
}
