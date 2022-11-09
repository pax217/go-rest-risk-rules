package chargebacks_test

import (
	"context"
	"errors"
	"testing"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/chargebacks"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	connectionError        = "connection lost"
	emailEmptyWarning      = "warning, email is empty"
	existChargebackWarning = "warning, a chargeback already exists with this ID"
)

func Test_SavePayerChargeback(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("When save is success", func(t *testing.T) {
		request := testdata.GetDefaultPayer().NewPayerFromPostRequest()

		repositoryMock := new(mocks.ChargebackRepositoryMock)
		service := chargebacks.NewChargebacksService(configs, repositoryMock, logger, new(datadog.MetricsDogMock))

		repositoryMock.On("Find", context.Background(), request).Return(entities.Payer{}, nil)
		repositoryMock.On("Save", context.Background(), request).Return(nil)
		err := service.Save(context.Background(), request)

		assert.Nil(t, err)
	})

	t.Run("When exist Payer then update only chargebacks", func(t *testing.T) {
		request := testdata.GetDefaultPayer().NewPayerFromPostRequest()
		payerFound := request

		repositoryMock := new(mocks.ChargebackRepositoryMock)
		service := chargebacks.NewChargebacksService(configs, repositoryMock, logger, new(datadog.MetricsDogMock))

		found := testdata.GetPayerWithDistinctChargebackID().NewPayerFromPostRequest()
		repositoryMock.On("Find", context.Background(), request).Return(payerFound, nil)
		payerFound.ID = found.ID
		payerFound.Chargebacks = append(payerFound.Chargebacks, found.Chargebacks...)
		repositoryMock.On("Update", context.Background(), payerFound).Return(nil)

		err := service.Save(context.Background(), request)

		assert.Nil(t, err)
	})

	t.Run("When find return error", func(t *testing.T) {
		request := testdata.GetDefaultPayer().NewPayerFromPostRequest()
		expectedError := errors.New(connectionError)

		repositoryMock := new(mocks.ChargebackRepositoryMock)
		service := chargebacks.NewChargebacksService(configs, repositoryMock, logger, new(datadog.MetricsDogMock))

		repositoryMock.On("Find", context.Background(), request).Return(entities.Payer{}, expectedError).Once()

		err := service.Save(context.Background(), request)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("When save return error", func(t *testing.T) {
		request := testdata.GetDefaultPayer().NewPayerFromPostRequest()
		expectedError := errors.New(connectionError)

		repositoryMock := new(mocks.ChargebackRepositoryMock)
		service := chargebacks.NewChargebacksService(configs, repositoryMock, logger, new(datadog.MetricsDogMock))

		repositoryMock.On("Find", context.Background(), request).Return(entities.Payer{}, nil)
		repositoryMock.On("Save", context.Background(), request).Return(expectedError)

		err := service.Save(context.Background(), request)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("When email is empty", func(t *testing.T) {
		request := testdata.GetDefaultPayer().NewPayerFromPostRequest()
		request.Email = ""

		repositoryMock := new(mocks.ChargebackRepositoryMock)
		service := chargebacks.NewChargebacksService(configs, repositoryMock, logger, new(datadog.MetricsDogMock))

		err := service.Save(context.Background(), request)

		assert.Nil(t, err)
	})

	t.Run("When already exist chargeback", func(t *testing.T) {
		request := testdata.GetDefaultPayer().NewPayerFromPostRequest()

		repositoryMock := new(mocks.ChargebackRepositoryMock)
		service := chargebacks.NewChargebacksService(configs, repositoryMock, logger, new(datadog.MetricsDogMock))
		repositoryMock.On("Find", context.Background(), request).Return(request, nil)

		err := service.Save(context.Background(), request)

		assert.Nil(t, err)
	})

	t.Run("When update return error", func(t *testing.T) {
		request := testdata.GetDefaultPayer().NewPayerFromPostRequest()
		payerFound := request
		expectedError := errors.New(connectionError)

		repositoryMock := new(mocks.ChargebackRepositoryMock)
		service := chargebacks.NewChargebacksService(configs, repositoryMock, logger, new(datadog.MetricsDogMock))

		found := testdata.GetPayerWithDistinctChargebackID().NewPayerFromPostRequest()

		repositoryMock.On("Find", context.Background(), request).Return(found, nil).Once()
		payerFound.ID = found.ID
		payerFound.Chargebacks = append(payerFound.Chargebacks, found.Chargebacks...)
		repositoryMock.On("Update", context.Background(), payerFound).Return(expectedError)

		err := service.Save(context.Background(), request)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
	})
}
