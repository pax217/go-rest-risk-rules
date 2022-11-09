package modules

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getModule() entities.Module {
	now := time.Now().Truncate(time.Millisecond)

	return entities.Module{
		ID:          primitive.NewObjectID(),
		CreatedAt:   now,
		UpdatedAt:   &now,
		CreatedBy:   "santiago.ceron@conekta.com",
		Name:        "policy_compliance",
		Description: "Regla para validar contratos con OXXO",
	}
}

func TestServiceAdd_WhenSearchOnFindRepositoryFailsThenReturnError(t *testing.T) {

	modulesRepository := new(mocks.ModulesRepositoryMock)
	service := NewModuleService(config.Config{}, modulesRepository, nil, nil)
	expectedError := errors.New("connection lost")
	module := getModule()

	modulesRepository.Mock.On("Get", nil, module.GetModuleFilter(false)).
		Return([]entities.Module{}, expectedError).Once()

	err := service.Add(nil, module)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError.Error(), err.Error())
	modulesRepository.AssertExpectations(t)
}

func TestAdd_WhenFoundADuplicatedModuleThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	modulesRepository := new(mocks.ModulesRepositoryMock)
	service := NewModuleService(config.Config{}, modulesRepository, logger, nil)
	module := getModule()
	modulesFound := make([]entities.Module, 0)
	modulesFound = append(modulesFound, module)
	modulesRepository.Mock.On("Get", nil, module.GetModuleFilter(false)).
		Return(modulesFound, nil).Once()

	err := service.Add(nil, module)

	assert.NotNil(t, err)
	assert.IsType(t, exceptions.NewDuplicatedException(""), err)
	assert.Equal(t, "module with name 'policy_compliance' is already created", err.Error())
	modulesRepository.AssertExpectations(t)
}

func TestAdd_WhenSaveOnRepositoryFailsThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	modulesRepository := new(mocks.ModulesRepositoryMock)
	service := NewModuleService(config.Config{}, modulesRepository, logger, new(datadog.MetricsDogMock))
	module := getModule()
	expectedError := errors.New("connection lost")

	modulesRepository.Mock.On("Get", nil, module.GetModuleFilter(false)).
		Return([]entities.Module{}, nil).Once()
	modulesRepository.Mock.On("Save", nil, &module).
		Return(expectedError).Once()

	err := service.Add(nil, module)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError.Error(), err.Error())
	modulesRepository.AssertExpectations(t)
}

func TestAdd_WhenIsSavedOkThenReturnNil(t *testing.T) {
	logger, _ := logs.New()
	modulesRepository := new(mocks.ModulesRepositoryMock)
	service := NewModuleService(config.Config{}, modulesRepository, logger, new(datadog.MetricsDogMock))
	module := getModule()

	modulesRepository.Mock.On("Get", nil, module.GetModuleFilter(false)).
		Return([]entities.Module{}, nil).Once()
	modulesRepository.Mock.On("Save", nil, &module).
		Return(nil).Once()

	err := service.Add(nil, module)

	assert.Nil(t, err)
	modulesRepository.AssertExpectations(t)
}

func TestUpdate_WhenUpdateOnRepositoryFailsThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	modulesRepository := new(mocks.ModulesRepositoryMock)
	service := NewModuleService(config.Config{}, modulesRepository, logger, new(datadog.MetricsDogMock))
	expectedError := errors.New("connection lost")

	id := "60f6f32ba0f965ae8ae2c87e"
	moduleRequest := entities.ModuleRequest{
		Name:        "Name Module",
		Description: "Description Test",
	}

	idModule, _ := primitive.ObjectIDFromHex(id)
	updatedAt := time.Now().Truncate(time.Millisecond)
	module := entities.Module{
		ID:          idModule,
		Name:        moduleRequest.Name,
		Description: moduleRequest.Description,
		UpdatedAt:   &updatedAt,
	}
	modulesRepository.Mock.On("Get", nil, module.GetModuleFilter(false)).
		Return([]entities.Module{}, nil)

	modulesRepository.Mock.On("Update", nil, id, module).
		Return(expectedError).Once()

	err := service.Update(nil, id, module)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	modulesRepository.AssertExpectations(t)
}

func TestUpdate_WhenUpdateOkThenReturnNil(t *testing.T) {
	logger, _ := logs.New()
	modulesRepository := new(mocks.ModulesRepositoryMock)
	service := NewModuleService(config.Config{}, modulesRepository, logger, new(datadog.MetricsDogMock))

	id := "60f6f32ba0f965ae8ae2c87e"
	moduleRequest := entities.ModuleRequest{
		Name:        "Name Module",
		Description: "Description Test",
	}

	idModule, _ := primitive.ObjectIDFromHex(id)
	updatedAt := time.Now().Truncate(time.Millisecond)
	module := entities.Module{
		ID:          idModule,
		Name:        moduleRequest.Name,
		Description: moduleRequest.Description,
		UpdatedAt:   &updatedAt,
	}
	modulesRepository.Mock.On("Get", nil, module.GetModuleFilter(false)).
		Return([]entities.Module{}, nil)

	modulesRepository.Mock.On("Update", nil, id, module).
		Return(nil).Once()

	err := service.Update(nil, id, module)

	assert.Nil(t, err)
	modulesRepository.AssertExpectations(t)
}

func TestUpdate_WhenUpdateFindByNameReturnError(t *testing.T) {
	logger, _ := logs.New()
	expectedError := errors.New("error")
	modulesRepository := new(mocks.ModulesRepositoryMock)
	id := "61086352928b571237eab678"
	idModule, _ := primitive.ObjectIDFromHex(id)

	moduleRequest := entities.ModuleRequest{
		Name:        "Name Module",
		Description: "Description Test",
	}
	updatedAt := time.Now().Truncate(time.Millisecond)
	module := entities.Module{
		ID:          idModule,
		Name:        moduleRequest.Name,
		Description: moduleRequest.Description,
		UpdatedAt:   &updatedAt,
	}
	modulesRepository.On("Get", nil, module.GetModuleFilter(false)).
		Return([]entities.Module{}, expectedError).Once()

	service := NewModuleService(config.Config{}, modulesRepository, logger, nil)

	err := service.Update(nil, id, module)

	assert.NotNil(t, err)
	assert.IsType(t, expectedError.Error(), err.Error())
	modulesRepository.AssertExpectations(t)
	assert.True(t, modulesRepository.AssertExpectations(t))
}

func TestUpdate_WhenFoundADuplicatedModuleThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	modulesRepository := new(mocks.ModulesRepositoryMock)
	service := NewModuleService(config.Config{}, modulesRepository, logger, nil)
	module := getModule()
	id := "60f6f32ba0f965ae8ae2c87e"
	modulesFound := make([]entities.Module, 0)
	modulesFound = append(modulesFound, module)

	modulesRepository.Mock.On("Get", nil, module.GetModuleFilter(false)).
		Return(modulesFound, nil)

	err := service.Update(nil, id, module)

	assert.NotNil(t, err)
	assert.IsType(t, exceptions.NewDuplicatedException(""), err)
	assert.Equal(t, "module with name 'policy_compliance' is already created", err.Error())
	modulesRepository.AssertExpectations(t)
}

func TestDelete_WhenDeleteOnRepositoryFailsThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	modulesRepository := new(mocks.ModulesRepositoryMock)
	service := NewModuleService(config.Config{}, modulesRepository, logger, new(datadog.MetricsDogMock))
	expectedError := errors.New("connection lost")

	id := "60f6f32ba0f965ae8ae2c87e"

	modulesRepository.Mock.On("Delete", nil, id).
		Return(expectedError).Once()

	err := service.Delete(nil, id)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	modulesRepository.AssertExpectations(t)
}

func TestDelete_WhenDeleteOkThenReturnNil(t *testing.T) {
	logger, _ := logs.New()
	modulesRepository := new(mocks.ModulesRepositoryMock)
	service := NewModuleService(config.Config{}, modulesRepository, logger, new(datadog.MetricsDogMock))

	id := "60f6f32ba0f965ae8ae2c87e"

	modulesRepository.Mock.On("Delete", nil, id).
		Return(nil).Once()

	err := service.Delete(nil, id)

	assert.Nil(t, err)
	modulesRepository.AssertExpectations(t)
}

func TestModuleService_GetAll(t *testing.T) {

	t.Run("test when repository passes directly data to service", func(t *testing.T) {
		pagination := entities.NewDefaultPagination()
		notOperatorsFound := entities.PagedResponse{
			HasMore: false,
			Total:   0,
			Object:  "",
			Data:    nil,
		}
		modulesRepositoryMock := new(mocks.ModulesRepositoryMock)
		modulesRepositoryMock.On("GetPaged", nil, pagination, entities.ModuleFilter{Paged: true}).Return(notOperatorsFound, nil)
		service := NewModuleService(config.Config{}, modulesRepositoryMock, nil, nil)

		pagedRules, err := service.GetAll(nil, pagination, entities.ModuleFilter{Paged: true})

		assert.Nil(t, err)
		assert.Equal(t, notOperatorsFound, pagedRules)
	})
}

func TestModuleService_Get(t *testing.T) {

	t.Run("test when repository passes directly data to service", func(t *testing.T) {
		moduleFilter := entities.ModuleFilter{
			Name:  "",
			Paged: false,
		}

		repositoryMock := new(mocks.ModulesRepositoryMock)
		repositoryMock.On("Get",
			nil,
			moduleFilter).
			Return([]entities.Module{}, nil)
		service := NewModuleService(config.Config{}, repositoryMock, nil, nil)

		modules, err := service.GetAll(nil, entities.NewDefaultPagination(), entities.ModuleFilter{Paged: false})

		assert.Nil(t, err)
		assert.Equal(t, []entities.Module{}, modules)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("test when repository return error for invalid id", func(t *testing.T) {
		invalidID := "invalidID"
		expectedError := exceptions.NewInvalidRequest(fmt.Sprintf("error: invalid id: %s", invalidID))

		moduleFilter := entities.ModuleFilter{
			ID:    invalidID,
			Paged: false,
		}

		repositoryMock := new(mocks.ModulesRepositoryMock)
		repositoryMock.On("Get",
			context.TODO(),
			moduleFilter).
			Return([]entities.Module{}, expectedError)
		service := NewModuleService(config.Config{}, repositoryMock, nil, nil)

		_, err := service.GetAll(context.TODO(), entities.NewDefaultPagination(), moduleFilter)

		assert.NotNil(t, err)
		assert.Equal(t, err, expectedError)
		repositoryMock.AssertExpectations(t)
	})
}
