package lists_test

import (
	"context"
	"errors"
	"testing"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/lists"

	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_ServiceList_Add(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("When add list success", func(t *testing.T) {
		listServiceMock := new(mocks.RkListsRestClient)
		ctx := context.TODO()
		request := testdata.GetDefaultRequestLists()

		service := lists.NewListsService(configs, logger, new(datadog.MetricsDogMock), listServiceMock)
		expectedList := testdata.GetDefaultWhiteList(true)

		listServiceMock.Mock.On("ListsSearch", ctx, request).
			Return([]entities.List{expectedList}, nil)

		foundLists, err := service.GetLists(ctx, request)

		assert.Nil(t, err)
		assert.NotEmpty(t, foundLists)
		assert.True(t, listServiceMock.AssertExpectations(t))
	})

	t.Run("on service client fails", func(t *testing.T) {
		listServiceMock := new(mocks.RkListsRestClient)
		ctx := context.TODO()
		expErr := errors.New("service connection lost")
		request := testdata.GetDefaultRequestLists()

		service := lists.NewListsService(configs, logger, new(datadog.MetricsDogMock), listServiceMock)

		listServiceMock.On("ListsSearch", ctx, request).
			Return([]entities.List{}, expErr).Once()

		foundLists, err := service.GetLists(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		assert.Empty(t, foundLists)
		listServiceMock.AssertExpectations(t)
	})
}
