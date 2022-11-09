package conditions_test

import (
	"errors"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/conditions"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestConditionService_Add_fails(t *testing.T) {
	log, _ := logs.New()
	condition := entities.Condition{
		CreatedBy:   "santiago.ceron@conekta.com",
		Name:        "or",
		Description: "debe cumplir al menos una de las condiciones",
	}

	t.Run("when search on repository fail", func(t *testing.T) {
		expectedError := errors.New("database connection lost")
		mockedRepository := new(mocks.ConditionsRepositoryMock)
		mockedRepository.On("FindByName", nil, condition).
			Return(entities.Condition{}, expectedError).Once()

		service := conditions.NewConditionsService(config.Config{}, mockedRepository, log, nil)

		err := service.Add(nil, condition)

		assert.NotNil(t, err)
		mockedRepository.AssertExpectations(t)
		assert.True(t, mockedRepository.AssertExpectations(t))
	})

	t.Run("when condition is duplicated", func(t *testing.T) {
		mockedRepository := new(mocks.ConditionsRepositoryMock)
		mockedRepository.On("FindByName", nil, condition).
			Return(condition, nil).Once()

		service := conditions.NewConditionsService(config.Config{}, mockedRepository, log, nil)

		err := service.Add(nil, condition)

		assert.NotNil(t, err)
		assert.IsType(t, exceptions.NewDuplicatedException(""), err)
		mockedRepository.AssertExpectations(t)
		assert.True(t, mockedRepository.AssertExpectations(t))
	})

	t.Run("when save on repository fail", func(t *testing.T) {
		expectedError := errors.New("database connection lost")
		mockedRepository := new(mocks.ConditionsRepositoryMock)
		mockedRepository.On("FindByName", nil, condition).
			Return(entities.Condition{}, nil).Once()
		mockedRepository.On("Add", nil, &condition).
			Return(expectedError).Once()

		service := conditions.NewConditionsService(config.Config{},
			mockedRepository, log, new(datadog.MetricsDogMock))

		err := service.Add(nil, condition)

		assert.NotNil(t, err)
		mockedRepository.AssertExpectations(t)
		assert.True(t, mockedRepository.AssertExpectations(t))
	})
}

func TestConditionService_Update_fails(t *testing.T) {
	log, _ := logs.New()
	condition := entities.Condition{
		CreatedBy:   "santiago.ceron@conekta.com",
		Name:        "or",
		Description: "debe cumplir al menos una de las condiciones",
	}

	t.Run("when condition is duplicated", func(t *testing.T) {
		mockedRepository := new(mocks.ConditionsRepositoryMock)
		mockedRepository.On("FindByName", nil, condition).
			Return(condition, nil).Once()
		id := "61086352928b571237eab678"

		service := conditions.NewConditionsService(config.Config{}, mockedRepository, log, nil)

		err := service.Update(nil, id, condition)

		assert.NotNil(t, err)
		assert.IsType(t, exceptions.NewDuplicatedException(""), err)
		mockedRepository.AssertExpectations(t)
		assert.True(t, mockedRepository.AssertExpectations(t))
	})

	t.Run("when FindByName return error ", func(t *testing.T) {
		expectedError := errors.New("error")
		mockedRepository := new(mocks.ConditionsRepositoryMock)
		mockedRepository.On("FindByName", nil, condition).
			Return(entities.Condition{}, expectedError).Once()
		id := "61086352928b571237eab678"

		service := conditions.NewConditionsService(config.Config{}, mockedRepository, log, nil)

		err := service.Update(nil, id, condition)

		assert.NotNil(t, err)
		assert.IsType(t, expectedError.Error(), err.Error())
		mockedRepository.AssertExpectations(t)
		assert.True(t, mockedRepository.AssertExpectations(t))
	})
}

func TestConditionService_Add_success(t *testing.T) {
	log, _ := logs.New()
	condition := entities.Condition{
		CreatedBy:   "santiago.ceron@conekta.com",
		Name:        "or",
		Description: "debe cumplir al menos una de las condiciones",
	}

	mockedRepository := new(mocks.ConditionsRepositoryMock)
	mockedRepository.On("FindByName", nil, condition).
		Return(entities.Condition{}, nil).Once()
	mockedRepository.On("Add", nil, &condition).
		Return(nil).Once()

	service := conditions.NewConditionsService(config.Config{}, mockedRepository, log, new(datadog.MetricsDogMock))

	err := service.Add(nil, condition)

	assert.Nil(t, err)
	assert.True(t, mockedRepository.AssertExpectations(t))
}

func TestUpdate_WhenUpdateOnRepositoryFailsThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	conditionsRepositoryMock := new(mocks.ConditionsRepositoryMock)
	service := conditions.NewConditionsService(config.Config{}, conditionsRepositoryMock, logger, new(datadog.MetricsDogMock))
	expectedError := errors.New("connection lost")

	id := "61086352928b571237eab678"
	conditionRequest := entities.ConditionRequest{
		Author:      "carlos.maldonado@gmail.com",
		Name:        "or",
		Description: "debe cumplir al menos una de las condiciones",
	}

	idCondition, _ := primitive.ObjectIDFromHex(id)
	updatedAt := time.Now().Truncate(time.Millisecond)
	condition := entities.Condition{
		ID:          idCondition,
		Name:        conditionRequest.Name,
		Description: conditionRequest.Description,
		UpdatedAt:   &updatedAt,
		UpdatedBy:   &conditionRequest.Author,
	}
	conditionsRepositoryMock.On("FindByName", nil, condition).
		Return(entities.Condition{}, nil).Once()

	conditionsRepositoryMock.Mock.On("Update", nil, id, condition).
		Return(expectedError).Once()

	err := service.Update(nil, id, condition)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	conditionsRepositoryMock.AssertExpectations(t)
}

func TestUpdate_WhenUpdateOkThenReturnNil(t *testing.T) {
	logger, _ := logs.New()
	conditionsRepositoryMock := new(mocks.ConditionsRepositoryMock)
	service := conditions.NewConditionsService(config.Config{}, conditionsRepositoryMock, logger, new(datadog.MetricsDogMock))

	id := "61086352928b571237eab678"
	conditionRequest := entities.ConditionRequest{
		Author:      "carlos.maldonado@conekta.com",
		Name:        "Name Module",
		Description: "Description Test",
	}

	idCondition, _ := primitive.ObjectIDFromHex(id)
	updatedAt := time.Now().Truncate(time.Millisecond)
	condition := entities.Condition{
		ID:          idCondition,
		Name:        conditionRequest.Name,
		Description: conditionRequest.Description,
		UpdatedAt:   &updatedAt,
		UpdatedBy:   &conditionRequest.Author,
	}

	conditionsRepositoryMock.On("FindByName", nil, condition).
		Return(entities.Condition{}, nil).Once()

	conditionsRepositoryMock.Mock.On("Update", nil, id, condition).
		Return(nil).Once()

	err := service.Update(nil, id, condition)

	assert.Nil(t, err)
	conditionsRepositoryMock.AssertExpectations(t)
}

func Test_Delete_Condition_Sevice(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when delete on repository fails then return error", func(t *testing.T) {
		conditionsRepositoryMock := new(mocks.ConditionsRepositoryMock)
		service := conditions.NewConditionsService(config.Config{}, conditionsRepositoryMock, logger, new(datadog.MetricsDogMock))
		expectedError := errors.New("connection lost")

		id := "60f6f32ba0f965ae8ae2c87e"

		conditionsRepositoryMock.Mock.On("Delete", nil, id).
			Return(expectedError).Once()

		err := service.Delete(nil, id)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
		conditionsRepositoryMock.AssertExpectations(t)
	})

	t.Run("when delete o then return nil", func(t *testing.T) {
		conditionsRepositoryMock := new(mocks.ConditionsRepositoryMock)
		service := conditions.NewConditionsService(config.Config{}, conditionsRepositoryMock, logger, new(datadog.MetricsDogMock))

		id := "60f6f32ba0f965ae8ae2c87e"

		conditionsRepositoryMock.Mock.On("Delete", nil, id).
			Return(nil).Once()

		err := service.Delete(nil, id)

		assert.Nil(t, err)
		conditionsRepositoryMock.AssertExpectations(t)
	})
}

func TestConditionsService_GetAll(t *testing.T) {
	logger, _ := logs.New()

	t.Run("test when repository passes directly data to service", func(t *testing.T) {
		pagination := entities.NewDefaultPagination()
		conditionID := "some_id"
		notConditionsFound := entities.PagedResponse{
			HasMore: false,
			Total:   0,
			Object:  "",
			Data:    nil,
		}
		conditionsRepositoryMock := new(mocks.ConditionsRepositoryMock)
		conditionsRepositoryMock.On("GetAll",
			nil,
			entities.ConditionsFilter{
				ID: conditionID,
			},
			pagination).
			Return(notConditionsFound, nil)

		service := conditions.NewConditionsService(config.Config{}, conditionsRepositoryMock, logger, new(datadog.MetricsDogMock))

		pagedRules, err := service.GetAll(nil,
			entities.ConditionsFilter{
				ID: conditionID,
			}, pagination)

		assert.Nil(t, err)
		assert.Equal(t, notConditionsFound, pagedRules)
	})
}
