package operators

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOperatorService_AddOperator_Fail(t *testing.T) {
	log, _ := logs.New()

	t.Run("When the operator already exist Fail", func(t *testing.T) {
		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		service := NewOperatorService(config.NewConfig(), log, operatorRepositoryMock, nil)
		operator := testdata.GetOperators()[0]
		filter := entities.OperatorFilter{
			Type: operator.Type,
			Name: operator.Name,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		operatorRepositoryMock.On("Get", ctx, filter).
			Once().
			Return(testdata.GetOperators(), nil)

		err := service.AddOperator(ctx, operator)

		assert.Error(t, err)
		assert.Errorf(t, err, fmt.Sprintf("the operator %s already exist", operator.Name))
		assert.True(t, operatorRepositoryMock.AssertNotCalled(t, "Save", ctx, &operator))
		assert.True(t, operatorRepositoryMock.AssertExpectations(t))
	})

	t.Run("when the database returns a error", func(t *testing.T) {
		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		service := NewOperatorService(config.NewConfig(), log, operatorRepositoryMock, nil)
		operator := testdata.GetOperators()[0]
		filter := entities.OperatorFilter{
			Type: operator.Type,
			Name: operator.Name,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		operatorRepositoryMock.On("Get", ctx, filter).
			Once().
			Return(make([]entities.Operator, 0), errors.New("unexpected error"))

		err := service.AddOperator(ctx, operator)

		assert.Error(t, err)
		assert.Errorf(t, err, "unexpected error")
		assert.True(t, operatorRepositoryMock.AssertNotCalled(t, "Save", ctx, &operator))
		assert.True(t, operatorRepositoryMock.AssertExpectations(t))

	})

	t.Run("save operator unsuccessfully ", func(t *testing.T) {
		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		service := NewOperatorService(config.NewConfig(), log, operatorRepositoryMock, new(datadog.MetricsDogMock))
		operator := testdata.GetOperators()[0]
		filter := entities.OperatorFilter{
			Type: operator.Type,
			Name: operator.Name,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		operatorRepositoryMock.On("Get", ctx, filter).
			Once().
			Return(make([]entities.Operator, 0), nil)
		operatorRepositoryMock.On("Save", ctx, mock.Anything).
			Once().
			Return(errors.New("unknown error"))
		err := service.AddOperator(ctx, operator)

		assert.Error(t, err)
		assert.Errorf(t, err, "unknown error")
		assert.True(t, operatorRepositoryMock.AssertExpectations(t))
	})

}
func TestOperatorService_AddOperator_Success(t *testing.T) {
	log, _ := logs.New()

	t.Run("save operator successfully", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		service := NewOperatorService(config.NewConfig(), log, operatorRepositoryMock, new(datadog.MetricsDogMock))
		operator := testdata.GetOperators()[0]
		filter := entities.OperatorFilter{
			Type: operator.Type,
			Name: operator.Name,
		}
		operatorRepositoryMock.On("Get", ctx, filter).
			Once().
			Return(make([]entities.Operator, 0), nil)
		operatorRepositoryMock.On("Save", ctx, mock.Anything).
			Once().
			Return(nil)

		err := service.AddOperator(ctx, operator)

		assert.NoError(t, err)
		assert.True(t, operatorRepositoryMock.AssertExpectations(t))
	})
}

func TestOperatorService_GetPaged(t *testing.T) {

	t.Run("test when repository passes directly data to service", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    nil,
		}
		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		operatorFilter := entities.OperatorFilter{
			Type:  "string",
			Paged: true,
		}

		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		operatorRepositoryMock.On("GetPaged",
			nil,
			operatorFilter,
			entities.NewDefaultPagination()).
			Return(serviceResponse, nil)
		service := NewOperatorService(config.Config{}, nil, operatorRepositoryMock, nil)

		pagedRules, err := service.Get(nil, operatorFilter, entities.NewDefaultPagination())

		assert.Nil(t, err)
		assert.Equal(t, serviceResponse, pagedRules)
		operatorRepositoryMock.AssertExpectations(t)
	})

	t.Run("test when repository return error for invalid id", func(t *testing.T) {
		invalidID := "invalidID"
		expectedError := exceptions.NewInvalidRequest(fmt.Sprintf("error: invalid id: %s", invalidID))

		operatorFilter := entities.OperatorFilter{
			ID:    invalidID,
			Paged: true,
		}

		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		operatorRepositoryMock.On("GetPaged",
			context.TODO(),
			operatorFilter,
			entities.NewDefaultPagination()).
			Return(entities.PagedResponse{}, expectedError)
		service := NewOperatorService(config.Config{}, nil, operatorRepositoryMock, nil)

		_, err := service.Get(context.TODO(), operatorFilter, entities.NewDefaultPagination())

		assert.NotNil(t, err)
		assert.Equal(t, err, expectedError)
		operatorRepositoryMock.AssertExpectations(t)
	})
}

func TestOperatorService_Get(t *testing.T) {

	t.Run("test when repository passes directly data to service", func(t *testing.T) {
		operatorFilter := entities.OperatorFilter{
			Type:  "string",
			Paged: false,
		}

		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		operatorRepositoryMock.On("Get",
			nil,
			operatorFilter).
			Return([]entities.Operator{}, nil)
		service := NewOperatorService(config.Config{}, nil, operatorRepositoryMock, nil)

		pagedRules, err := service.Get(nil, operatorFilter, entities.NewDefaultPagination())

		assert.Nil(t, err)
		assert.Equal(t, []entities.Operator{}, pagedRules)
		operatorRepositoryMock.AssertExpectations(t)
	})
}

func Test_Operator_Service_Update(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when update on repository fails then return error", func(t *testing.T) {
		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		service := NewOperatorService(config.Config{}, logger, operatorRepositoryMock, new(datadog.MetricsDogMock))
		expectedError := exceptions.NewDuplicatedException("the operator + of type  already exist")

		id := "61086352928b571237eab678"
		operatorRequest := entities.OperatorRequest{
			Author:      "carlos.maldonado@gmail.com",
			Name:        "+",
			Title:       "Más",
			Description: "indicates the sum",
		}

		idCondition, _ := primitive.ObjectIDFromHex(id)

		updatedAt := time.Now().Truncate(time.Millisecond)
		operator := entities.Operator{
			ID:          idCondition,
			Name:        operatorRequest.Name,
			Title:       operatorRequest.Title,
			Description: operatorRequest.Description,
			UpdatedAt:   &updatedAt,
			UpdatedBy:   &operatorRequest.Author,
		}
		filter := entities.OperatorFilter{
			Type: operator.Type,
			Name: operator.Name,
		}

		operatorRepositoryMock.On("Get", context.TODO(), filter).
			Once().
			Return(make([]entities.Operator, 0), nil)

		operatorRepositoryMock.Mock.On("Update",
			context.TODO(), id, operator).
			Once().
			Return(expectedError)

		err := service.Update(context.TODO(), id, operator)

		assert.NotNil(t, err)
		assert.Errorf(t, err, fmt.Sprintf("the operator %s of type %s already exist", operator.Name, operator.Type))
		operatorRepositoryMock.AssertExpectations(t)
	})

	t.Run("when update ok then return nil", func(t *testing.T) {
		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		service := NewOperatorService(config.Config{}, logger, operatorRepositoryMock, new(datadog.MetricsDogMock))

		id := "61086352928b571237eab678"
		operatorRequest := entities.OperatorRequest{
			Author:      "carlos.maldonado@conekta.com",
			Name:        ">",
			Title:       "Mayor a",
			Description: "indica la desigualdad matemática de 2 numeros",
		}

		idCondition, _ := primitive.ObjectIDFromHex(id)

		updatedAt := time.Now().Truncate(time.Millisecond)
		operator := entities.Operator{
			ID:          idCondition,
			Name:        operatorRequest.Name,
			Title:       operatorRequest.Title,
			Description: operatorRequest.Description,
			UpdatedAt:   &updatedAt,
			UpdatedBy:   &operatorRequest.Author,
		}

		operatorRepositoryMock.On("Get", context.TODO(), entities.OperatorFilter{
			Type: operator.Type,
			Name: operator.Name,
		}).
			Once().
			Return(make([]entities.Operator, 0), nil)

		operatorRepositoryMock.Mock.On("Update",
			context.TODO(),
			id,
			operator).
			Once().
			Return(nil)

		err := service.Update(context.TODO(), id, operator)

		assert.Nil(t, err)
		operatorRepositoryMock.AssertExpectations(t)
	})

	t.Run("when update with invalid ID return error", func(t *testing.T) {
		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		service := NewOperatorService(config.Config{}, logger, operatorRepositoryMock, new(datadog.MetricsDogMock))

		id := "xxx"
		operatorRequest := entities.OperatorRequest{
			Author:      "carlos.maldonado@gmail.com",
			Name:        "+",
			Title:       "Más",
			Description: "indicates the sum",
		}

		idCondition, _ := primitive.ObjectIDFromHex(id)
		updatedAt := time.Now().Truncate(time.Millisecond)
		operator := entities.Operator{
			ID:          idCondition,
			Name:        operatorRequest.Name,
			Title:       operatorRequest.Title,
			Description: operatorRequest.Description,
			UpdatedAt:   &updatedAt,
			UpdatedBy:   &operatorRequest.Author,
		}

		err := service.Update(context.TODO(), id, operator)

		assert.NotNil(t, err)
		assert.Errorf(t, err, fmt.Sprintf("the operator %s of type %s already exist", operator.Name, operator.Type))
		operatorRepositoryMock.AssertExpectations(t)
	})
}

func Test_Delete_Operator_Sevice(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when delete on repository fails then return error", func(t *testing.T) {
		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		service := NewOperatorService(config.Config{}, logger, operatorRepositoryMock, new(datadog.MetricsDogMock))
		expectedError := errors.New("connection lost")

		id := "60f6f32ba0f965ae8ae2c87e"

		operatorRepositoryMock.Mock.On("Delete", nil, id).
			Return(expectedError).Once()

		err := service.Delete(nil, id)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
		operatorRepositoryMock.AssertExpectations(t)
	})

	t.Run("when delete ok then return nil", func(t *testing.T) {
		operatorRepositoryMock := new(mocks.OperatorRepositoryMock)
		service := NewOperatorService(config.Config{}, logger, operatorRepositoryMock, new(datadog.MetricsDogMock))

		id := "60f6f32ba0f965ae8ae2c87e"

		operatorRepositoryMock.Mock.On("Delete", nil, id).
			Return(nil).Once()

		err := service.Delete(nil, id)

		assert.Nil(t, err)
		operatorRepositoryMock.AssertExpectations(t)
	})
}
