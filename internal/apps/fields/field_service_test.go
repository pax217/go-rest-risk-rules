package fields

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/test/testdata"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFieldService_Update(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when update on repository fails then return error", func(t *testing.T) {
		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		service := NewFieldsService(configs, fieldsRepositoryMock, logger, new(datadog.MetricsDogMock))
		expectedError := errors.New("connection lost")

		id := "60f6f32ba0f965ae8ae2c87e"
		fieldRequest := entities.FieldRequest{
			Name:        "email",
			Type:        "string",
			Description: "Representa el campo email perteneciente al charge",
			Author:      "carlos.maldonado@conekta.com",
		}

		idModule, _ := primitive.ObjectIDFromHex(id)
		updatedAt := time.Now().Truncate(time.Millisecond)
		updatedBy := "carlos.maldonado@conekta.com"
		field := entities.Field{
			ID:          idModule,
			Name:        fieldRequest.Name,
			Description: fieldRequest.Description,
			UpdatedAt:   &updatedAt,
			UpdatedBy:   &updatedBy,
			Type:        fieldRequest.Type,
		}

		fieldsRepositoryMock.Mock.On(
			"GetFields",
			nil,
			field.GetFieldsFilter(false)).
			Return([]entities.Field{}, nil).Once()

		fieldsRepositoryMock.Mock.On("Update", nil, id, &field).
			Return(expectedError).Once()

		err := service.Update(nil, id, field)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
		fieldsRepositoryMock.AssertExpectations(t)
	})

	t.Run("when update ok then return nil", func(t *testing.T) {
		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		service := NewFieldsService(configs, fieldsRepositoryMock, logger, new(datadog.MetricsDogMock))

		id := "60f6f32ba0f965ae8ae2c87e"
		fieldRequest := entities.FieldRequest{
			Name:        "email",
			Type:        "string",
			Description: "Representa el campo email perteneciente al charge",
			Author:      "carlos.maldonado@conekta.com",
		}

		idModule, _ := primitive.ObjectIDFromHex(id)
		updatedAt := time.Now().Truncate(time.Millisecond)
		updatedBy := "carlos.maldonado@conekta.com"
		field := entities.Field{
			ID:          idModule,
			Name:        fieldRequest.Name,
			Description: fieldRequest.Description,
			UpdatedAt:   &updatedAt,
			UpdatedBy:   &updatedBy,
			Type:        fieldRequest.Type,
		}
		fieldsRepositoryMock.Mock.On(
			"GetFields",
			nil,
			field.GetFieldsFilter(false)).
			Return([]entities.Field{}, nil).Once()

		fieldsRepositoryMock.Mock.On("Update", nil, id, &field).
			Return(nil).Once()

		err := service.Update(nil, id, field)

		assert.Nil(t, err)
		fieldsRepositoryMock.AssertExpectations(t)
	})
}

func TestFieldService_Delete(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when delete on repository fails then return error", func(t *testing.T) {
		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		service := NewFieldsService(configs, fieldsRepositoryMock, logger, new(datadog.MetricsDogMock))
		expectedError := errors.New("connection lost")

		id := "60f6f32ba0f965ae8ae2c87e"

		fieldsRepositoryMock.Mock.On("Delete", nil, id).
			Return(expectedError).Once()

		err := service.Delete(nil, id)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
		fieldsRepositoryMock.AssertExpectations(t)
	})

	t.Run("when delete ok then return nil", func(t *testing.T) {
		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		service := NewFieldsService(configs, fieldsRepositoryMock, logger, new(datadog.MetricsDogMock))

		id := "60f6f32ba0f965ae8ae2c87e"

		fieldsRepositoryMock.Mock.On("Delete", nil, id).
			Return(nil).Once()

		err := service.Delete(nil, id)

		assert.Nil(t, err)
		fieldsRepositoryMock.AssertExpectations(t)
	})
}

func TestFieldService_AddField(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("on database field exist connection lost then return error", func(t *testing.T) {
		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		service := NewFieldsService(configs, fieldsRepositoryMock, logger, new(datadog.MetricsDogMock))
		field := testdata.GetDefaultField()
		expectedError := errors.New("database connection lost")

		fieldsRepositoryMock.Mock.On(
			"GetFields",
			nil,
			field.GetFieldsFilter(false)).
			Return([]entities.Field{}, expectedError)

		err := service.AddField(nil, field)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
		fieldsRepositoryMock.AssertExpectations(t)
	})

	t.Run("on database add field fails if field exist", func(t *testing.T) {
		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		service := NewFieldsService(configs, fieldsRepositoryMock, logger, new(datadog.MetricsDogMock))
		field := testdata.GetDefaultField()
		expectedError := exceptions.NewDuplicatedException(fmt.Sprintf("the field '%s' of type '%s' already exist", field.Name, field.Type))

		fieldsRepositoryMock.Mock.On(
			"GetFields",
			context.TODO(),
			field.GetFieldsFilter(false)).
			Return(testdata.GetDefaultArrayFields(), nil)

		err := service.AddField(context.TODO(), field)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
		fieldsRepositoryMock.AssertExpectations(t)
	})

	t.Run("on database add field connection lost then return error", func(t *testing.T) {
		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		service := NewFieldsService(configs, fieldsRepositoryMock, logger, new(datadog.MetricsDogMock))
		field := testdata.GetDefaultField()
		expectedError := errors.New("database connection lost")

		fieldsRepositoryMock.Mock.On(
			"GetFields",
			context.TODO(),
			field.GetFieldsFilter(false)).
			Return([]entities.Field{}, nil).Once()

		fieldsRepositoryMock.Mock.On(
			"AddField",
			context.TODO(),
			&field).
			Return(expectedError).Once()

		err := service.AddField(context.TODO(), field)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
		fieldsRepositoryMock.AssertExpectations(t)
	})

	t.Run("on database add field ok then return nil", func(t *testing.T) {
		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		service := NewFieldsService(configs, fieldsRepositoryMock, logger, new(datadog.MetricsDogMock))
		field := testdata.GetDefaultField()

		fieldsRepositoryMock.Mock.On(
			"GetFields",
			context.TODO(),
			field.GetFieldsFilter(false)).
			Return([]entities.Field{}, nil).Once()

		fieldsRepositoryMock.Mock.On(
			"AddField",
			context.TODO(),
			&field).
			Return(nil).Once()

		err := service.AddField(context.TODO(), field)

		assert.Nil(t, err)
		fieldsRepositoryMock.AssertExpectations(t)
	})
}

func TestFieldService_GetFieldsPaged(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("test when repository passes directly data to service", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    []entities.Field{},
		}
		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		q.Set("paged", "true")
		fieldsFilter := entities.FieldsFilter{Paged: true}

		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		fieldsRepositoryMock.On("GetFieldsPaged",
			nil,
			fieldsFilter,
			entities.NewDefaultPagination()).
			Return(serviceResponse, nil)

		service := NewFieldsService(
			configs,
			fieldsRepositoryMock,
			logger,
			nil)

		pagedRules, err := service.
			GetFields(
				nil,
				fieldsFilter,
				entities.NewDefaultPagination())

		assert.Nil(t, err)
		assert.Equal(t, serviceResponse, pagedRules)
	})
}

func TestFieldService_GetFields(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("test when repository passes directly data to service", func(t *testing.T) {
		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		q.Set("paged", "false")
		fieldsFilter := entities.FieldsFilter{
			Paged: false,
		}

		fieldsRepositoryMock := new(mocks.FieldsRepositoryMock)
		fieldsRepositoryMock.On("GetFields",
			nil,
			fieldsFilter).
			Return([]entities.Field{}, nil)

		service := NewFieldsService(
			configs,
			fieldsRepositoryMock,
			logger,
			nil)

		rules, err := service.
			GetFields(
				nil,
				fieldsFilter,
				entities.NewDefaultPagination())

		assert.Nil(t, err)
		assert.Equal(t, []entities.Field{}, rules)
	})
}
