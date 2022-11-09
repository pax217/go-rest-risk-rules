package familycom

import (
	"context"
	"errors"
	"fmt"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"testing"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestFamilyCompaniesService_Create(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when familyCompaniesRepository fails, then return error", func(t *testing.T) {
		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)

		ctx := context.TODO()
		request := testdata.GetFamilyCompaniesRequest()
		familyCompanies := request.NewFamilyCompaniesFromPostRequest()

		expErr := errors.New("database connection lost")
		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyCompaniesRepositoryMock.On("AddFamilyCompanies",
			ctx,
			&familyCompanies).
			Return(expErr).Once()

		familyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx,
			entities.FamilyCompaniesFilter{
				Name: familyCompanies.Name,
			}).Return([]entities.FamilyCompanies{}, nil).Once()

		err := service.Create(ctx, familyCompanies)

		assert.Error(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})

	t.Run("when familyCompaniesRepository found family company, then return already created", func(t *testing.T) {
		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		ctx := context.TODO()
		request := testdata.GetFamilyCompaniesRequest()
		familyComanies := entities.FamilyCompanies{
			ID:   request.NewFamilyCompaniesFromPostRequest().ID,
			Name: request.NewFamilyCompaniesFromPostRequest().Name,
		}
		expErr := errors.New(fmt.Sprintf(
			"family companies with name [%s] is already created",
			familyComanies.Name))
		service := NewFamilyCompaniesService(configs,
			familyCompaniesRepositoryMock,
			nil,
			logger,
			new(datadog.MetricsDogMock))

		familyCompaniesFound := testdata.GetDefaultFamilyCompanies()

		familyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx,
			entities.FamilyCompaniesFilter{
				Name: familyComanies.Name,
			}).Return([]entities.FamilyCompanies{familyCompaniesFound}, nil)

		err := service.Create(ctx, familyComanies)

		assert.Error(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		assert.Equal(t, exceptions.FamilyCompaniesNameDuplicated, err.(exceptions.DuplicatedException).Causes().Code)
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})

	t.Run("when familyCompanies create success", func(t *testing.T) {
		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		ctx := context.TODO()
		request := testdata.GetFamilyCompaniesRequest()
		familyCompanies := request.NewFamilyCompaniesFromPostRequest()
		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx,
			entities.FamilyCompaniesFilter{
				Name: familyCompanies.Name,
			}).
			Return([]entities.FamilyCompanies{}, nil).Once()

		familyCompaniesRepositoryMock.On("AddFamilyCompanies",
			ctx,
			&familyCompanies).Return(nil)

		err := service.Create(ctx, familyCompanies)

		assert.Nil(t, err)
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})
}

func TestFamilyCompaniesService_Get(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("test when repository paged is false", func(t *testing.T) {

		filter := entities.FamilyCompaniesFilter{Paged: false}

		repositoryMock := new(mocks.FamilyCompaniesRepositoryMock)

		repositoryMock.On("Search", context.TODO(), entities.FamilyCompaniesFilter{
			CompanyIDs: nil,
			Name:       "",
			Paged:      false,
		}).Return([]entities.FamilyCompanies{}, nil)

		service := NewFamilyCompaniesService(configs, repositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyCompanies, err := service.Get(context.TODO(), entities.NewDefaultPagination(), filter)

		assert.Nil(t, err)
		assert.Equal(t, []entities.FamilyCompanies{}, familyCompanies)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("test when repository paged is true", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    []entities.FamilyCompanies{},
		}
		filter := entities.FamilyCompaniesFilter{Paged: true}

		repositoryMock := new(mocks.FamilyCompaniesRepositoryMock)

		repositoryMock.On("SearchPaged", context.TODO(),
			entities.NewDefaultPagination(),
			entities.FamilyCompaniesFilter{
				CompanyIDs: nil,
				Name:       "",
				Paged:      true,
			}).Return(serviceResponse, nil)

		service := NewFamilyCompaniesService(configs, repositoryMock, nil, logger, new(datadog.MetricsDogMock))

		pagedFamilyCompanies, err := service.Get(context.TODO(), entities.NewDefaultPagination(), filter)

		assert.Nil(t, err)
		assert.Equal(t, serviceResponse, pagedFamilyCompanies)
		repositoryMock.AssertExpectations(t)
	})

}

func TestFamilyCompaniesService_GetFamilyCompaniesFromFilter(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when familyCompaniesRepository.GetFamilyCompanies fails, then return error", func(t *testing.T) {
		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		ctx := context.TODO()
		familyCompaniesFilter := testdata.GetDefaultFamilyCompaniesFilter()

		expErr := errors.New("database connection lost")
		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx, familyCompaniesFilter).Return([]entities.FamilyCompanies{}, expErr).Once()

		_, err := service.GetFamiliesCompaniesFromFilter(ctx, familyCompaniesFilter)

		assert.Error(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})

	t.Run("when there are 0 family companies found, it should return an empty family", func(t *testing.T) {
		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		ctx := context.TODO()
		FamilyCompaniesFilter := testdata.GetDefaultFamilyCompaniesFilter()

		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx, FamilyCompaniesFilter).Return([]entities.FamilyCompanies{}, nil).Once()

		familiesCompaniesFound, err := service.GetFamiliesCompaniesFromFilter(ctx, FamilyCompaniesFilter)

		assert.NoError(t, err)
		assert.Equal(t, []entities.FamilyCompanies{}, familiesCompaniesFound)
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})

	t.Run("when a family companies is found, it should return it", func(t *testing.T) {
		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		ctx := context.TODO()
		FamilyCompaniesFilter := testdata.GetDefaultFamilyCompaniesFilter()
		FamilyCompaniesMock := []entities.FamilyCompanies{testdata.GetDefaultFamilyCompanies()}
		expectedFamilyCompanies := FamilyCompaniesMock

		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx, FamilyCompaniesFilter).Return(FamilyCompaniesMock, nil).Once()

		familiesCompaniesFound, err := service.GetFamiliesCompaniesFromFilter(ctx, FamilyCompaniesFilter)

		assert.NoError(t, err)
		assert.Equal(t, expectedFamilyCompanies, familiesCompaniesFound)
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})
}

func TestFamilyCompaniesService_Delete(t *testing.T) {
	configs := config.NewConfig()
	logger, _ := logs.New()
	familyCompanyID := "61e991ad1214eac062ada43d"

	t.Run("test when delete family companies fail", func(t *testing.T) {
		expectedError := errors.New("error")

		rulesRepositoryMock := new(mocks.RulesRepositoryMock)
		rulesRepositoryMock.On("FindRulesPaged",
			context.TODO(),
			entities.RuleFilter{FamilyCompanyID: familyCompanyID},
			entities.Pagination{},
		).Return(entities.PagedResponse{
			Data: []entities.Rule{},
		}, nil)

		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		familyCompaniesRepositoryMock.On("Delete", context.TODO(), familyCompanyID).Return(expectedError)
		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, rulesRepositoryMock, logger, new(datadog.MetricsDogMock))

		err := service.Delete(context.TODO(), familyCompanyID)

		assert.NotNil(t, err)
		assert.Equal(t, err, expectedError)
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})

	t.Run("test when delete family companies successful", func(t *testing.T) {
		rulesRepositoryMock := new(mocks.RulesRepositoryMock)
		rulesRepositoryMock.On("FindRulesPaged",
			context.TODO(),
			entities.RuleFilter{FamilyCompanyID: familyCompanyID},
			entities.Pagination{},
		).Return(entities.PagedResponse{
			Data: []entities.Rule{},
		}, nil)

		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		familyCompaniesRepositoryMock.On("Delete", context.TODO(), familyCompanyID).Return(nil)
		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, rulesRepositoryMock, logger, new(datadog.MetricsDogMock))

		err := service.Delete(context.TODO(), familyCompanyID)

		assert.Nil(t, err)
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})

	t.Run("test when delete family companies is associated with an exiting rule", func(t *testing.T) {
		ruleAssociated := testdata.GetDefaultRuleWithFamilyCompanyID(false)
		expectedError := fmt.Sprintf(
			"family companies with id [%s], is associated with the rule [%s]",
			*ruleAssociated.FamilyCompanyID,
			ruleAssociated.Description)

		// family companies with id [61e991ad1214eac062ada43d], is associated with the rule [empty]

		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		rulesRepositoryMock := new(mocks.RulesRepositoryMock)
		rulesRepositoryMock.
			On("FindRulesPaged",
				context.TODO(),
				entities.RuleFilter{FamilyCompanyID: *ruleAssociated.FamilyCompanyID},
				entities.Pagination{},
			).
			Return(entities.PagedResponse{
				Data: []entities.Rule{ruleAssociated},
			}, nil)

		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, rulesRepositoryMock, logger, new(datadog.MetricsDogMock))

		err := service.Delete(context.TODO(), *ruleAssociated.FamilyCompanyID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedError)

		familyCompaniesRepositoryMock.AssertExpectations(t)
	})
}

func TestFamilyCompaniesService_Update(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when familyCompaniesRepository fails, then return error", func(t *testing.T) {
		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		ctx := context.TODO()
		expErr := errors.New("database connection lost")
		request := testdata.GetFamilyCompaniesRequest()
		familyCompanies := request.NewFamilyCompaniesFromPutRequest()
		id := "61685179378d2ad5c3405bc5"

		repository := new(mocks.FamilyCompaniesRepositoryMock)
		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyCompaniesRepositoryMock.On("Update", ctx, id, &familyCompanies).
			Return(expErr).Once()

		familyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx,
			entities.FamilyCompaniesFilter{
				Name: familyCompanies.Name,
			}).Return([]entities.FamilyCompanies{}, nil)

		err := service.Update(ctx, id, familyCompanies)

		assert.NotNil(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		repository.AssertExpectations(t)
	})

	t.Run("when familyCompaniesRepository on search fails, then return error", func(t *testing.T) {
		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		ctx := context.TODO()

		request := testdata.GetFamilyCompaniesRequest()
		family := request.NewFamilyCompaniesFromPutRequest()
		id := "61685179378d2ad5c3405bc5-x"
		expErr := errors.New(fmt.Sprintf("invalid format for family_company_id: [%s]", id))

		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx,
			entities.FamilyCompaniesFilter{
				Name: family.Name,
			}).Return([]entities.FamilyCompanies{}, expErr)

		err := service.Update(ctx, id, family)

		assert.NotNil(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})

	t.Run("when familyCompaniesRepository search found family companies, then return name duplicated", func(t *testing.T) {
		damilyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		ctx := context.TODO()

		request := testdata.GetFamilyRequest()
		family := entities.FamilyCompanies{
			ID:   request.NewFamilyFromPostRequest().ID,
			Name: request.NewFamilyFromPostRequest().Name,
		}
		id := "61685179378d2ad5c3405bc5-x"
		expErr := errors.New(fmt.Sprintf(
			"family companies name: [%s] is duplicated",
			family.Name))
		repository := new(mocks.FamilyRepositoryMock)
		service := NewFamilyCompaniesService(configs, damilyCompaniesRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyFound := testdata.GetDefaultFamilyCompanies()
		familyCompaniesFound := []entities.FamilyCompanies{
			{
				ID:         familyFound.ID,
				Name:       familyFound.Name,
				CompanyIDs: familyFound.CompanyIDs,
			},
		}
		damilyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx,
			entities.FamilyCompaniesFilter{
				CompanyIDs: family.CompanyIDs,
				Name:       family.Name,
			}).Return(familyCompaniesFound, nil)

		err := service.Update(ctx, id, family)

		assert.NotNil(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		assert.Equal(t, exceptions.FamilyCompaniesNameDuplicated, err.(exceptions.DuplicatedException).Causes().Code)
		repository.AssertExpectations(t)
	})

	t.Run("when update is success", func(t *testing.T) {
		familyCompaniesRepositoryMock := new(mocks.FamilyCompaniesRepositoryMock)
		ctx := context.TODO()
		request := testdata.GetFamilyCompaniesRequest()
		familyCompanies := request.NewFamilyCompaniesFromPutRequest()
		id := "61685179378d2ad5c3405bc5"

		service := NewFamilyCompaniesService(configs, familyCompaniesRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyCompaniesRepositoryMock.On("Update",
			ctx,
			id,
			&familyCompanies).Return(nil).Once()

		familyCompaniesRepositoryMock.On("GetFamilyCompanies",
			ctx,
			entities.FamilyCompaniesFilter{
				Name: familyCompanies.Name,
			}).Return([]entities.FamilyCompanies{}, nil)

		err := service.Update(ctx, id, familyCompanies)

		assert.Nil(t, err)
		familyCompaniesRepositoryMock.AssertExpectations(t)
	})
}
