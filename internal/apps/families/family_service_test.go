package families_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/families"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFamilyService_Create(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when familyRepository fails, then return error", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()
		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPostRequest()
		expErr := errors.New("database connection lost")
		service := families.
			NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyRepositoryMock.On("AddFamily", ctx, &family).
			Return(expErr).Once()

		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(entities.PagedResponse{Data: []entities.Family{}}, nil).Once()

		err := service.Create(ctx, family)

		assert.Error(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		familyRepositoryMock.AssertExpectations(t)
	})

	t.Run("when familyRepository search found families, then return mccs duplicated", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()
		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPostRequest()
		expErr := errors.New(fmt.Sprintf("family: [%s] has MCCs [%s] duplicated",
			family.Name, family.Mccs[0]))

		service := families.
			NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyFound := testdata.GetDefaultFamily()
		pagedResponseFamilies := entities.PagedResponse{
			Data: []entities.Family{
				{
					ID:   familyFound.ID,
					Name: familyFound.Name,
					Mccs: []string{family.Mccs[0]},
				},
			},
		}
		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(pagedResponseFamilies, nil).Once()

		err := service.Create(ctx, family)

		assert.Error(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		assert.Equal(t, exceptions.FamiliesMccsDuplicated, err.(exceptions.DuplicatedException).Causes().Code)

		familyRepositoryMock.AssertExpectations(t)
	})

	t.Run("when familyRepository search found families, then return name duplicated", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()
		request := testdata.GetFamilyRequest()
		family := entities.Family{
			ID:   request.NewFamilyFromPostRequest().ID,
			Name: request.NewFamilyFromPostRequest().Name,
		}
		expErr := errors.New(fmt.Sprintf("family name: [%s] is duplicated", family.Name))
		service := families.
			NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyFound := testdata.GetDefaultFamily()
		pagedResponseFamilies := entities.PagedResponse{
			Data: []entities.Family{
				{
					ID:   familyFound.ID,
					Name: familyFound.Name,
					Mccs: familyFound.Mccs,
				},
			},
		}
		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(pagedResponseFamilies, nil)

		err := service.Create(ctx, family)

		assert.Error(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		assert.Equal(t, exceptions.FamiliesNameDuplicated, err.(exceptions.DuplicatedException).Causes().Code)
		familyRepositoryMock.AssertExpectations(t)
	})

	t.Run("when family exist, then return error ", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()
		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPostRequest()
		expErr := errors.New(fmt.Sprintf("error, family %s. already exist", family.Name))
		service := families.
			NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(entities.PagedResponse{Data: []entities.Family{family}}, expErr)

		err := service.Create(ctx, family)

		assert.Error(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		familyRepositoryMock.AssertExpectations(t)
	})

	t.Run("when family create success", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()
		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPostRequest()
		service := families.
			NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyRepositoryMock.On("AddFamily", ctx, &family).
			Return(nil).Once()

		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(entities.PagedResponse{
			Data: []entities.Family{},
		}, nil)

		err := service.Create(ctx, family)

		assert.Nil(t, err)
		familyRepositoryMock.AssertExpectations(t)
	})
}

func TestFamilyService_SearchPaged(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("test when familyRepository passes directly data to service", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    []entities.Family{},
		}
		filter := entities.FamilyFilter{Paged: true}

		repositoryMock := new(mocks.FamilyRepositoryMock)

		repositoryMock.On("SearchPaged", context.TODO(), entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs:  nil,
			Name:  "",
			Paged: true,
		}).Return(serviceResponse, nil)

		service := families.NewFamilyService(configs, repositoryMock, nil, logger, new(datadog.MetricsDogMock))

		pagedFamilies, err := service.Get(context.TODO(), entities.NewDefaultPagination(), filter)

		assert.Nil(t, err)
		assert.Equal(t, serviceResponse, pagedFamilies)
		repositoryMock.AssertExpectations(t)
	})
}

func TestFamilyService_Search(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("test when repository passes directly data to service", func(t *testing.T) {

		filter := entities.FamilyFilter{Paged: false}

		repositoryMock := new(mocks.FamilyRepositoryMock)

		repositoryMock.On("Search", context.TODO(), entities.FamilyFilter{
			Mccs:  nil,
			Name:  "",
			Paged: false,
		}).Return([]entities.Family{}, nil)

		service := families.NewFamilyService(configs, repositoryMock, nil, logger, new(datadog.MetricsDogMock))

		pagedFamilies, err := service.Get(context.TODO(), entities.NewDefaultPagination(), filter)

		assert.Nil(t, err)
		assert.Equal(t, []entities.Family{}, pagedFamilies)
		repositoryMock.AssertExpectations(t)
	})
}

func TestFamilyService_Update(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when familyRepository fails, then return error", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()
		expErr := errors.New("database connection lost")
		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPutRequest()
		id := "61685179378d2ad5c3405bc5"

		repository := new(mocks.FamilyRepositoryMock)
		service := families.NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyRepositoryMock.On("Update", ctx, id, &family).
			Return(expErr).Once()

		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(entities.PagedResponse{Data: []entities.Family{}}, nil)

		err := service.Update(ctx, "61685179378d2ad5c3405bc5", family)

		assert.NotNil(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		repository.AssertExpectations(t)
	})

	t.Run("when familyRepository on search fails, then return error", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()

		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPutRequest()
		id := "61685179378d2ad5c3405bc5-x"
		expErr := errors.New(fmt.Sprintf("invalid format for family_id: [%s]", id))

		service := families.NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(entities.PagedResponse{Data: []entities.Family{}}, expErr)

		err := service.Update(ctx, id, family)

		assert.NotNil(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		familyRepositoryMock.AssertExpectations(t)
	})

	t.Run("when familyRepository search found families, then return mccs duplicated", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()

		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPutRequest()
		id := "61685179378d2ad5c3405bc5-x"
		expErr := errors.New(fmt.Sprintf(
			"family: [%s] has MCCs [%s] duplicated",
			family.Name,
			family.Mccs[0]))
		repository := new(mocks.FamilyRepositoryMock)
		service := families.NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyFound := testdata.GetDefaultFamily()
		pagedResponseFamilies := entities.PagedResponse{
			Data: []entities.Family{
				{
					ID:   familyFound.ID,
					Name: familyFound.Name,
					Mccs: []string{family.Mccs[0]},
				},
			},
		}
		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(pagedResponseFamilies, nil)

		err := service.Update(ctx, id, family)

		assert.NotNil(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		assert.Equal(t, exceptions.FamiliesMccsDuplicated, err.(exceptions.DuplicatedException).Causes().Code)
		repository.AssertExpectations(t)
	})

	t.Run("when familyRepository search found families, then return name duplicated", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()

		request := testdata.GetFamilyRequest()
		family := entities.Family{
			ID:   request.NewFamilyFromPostRequest().ID,
			Name: request.NewFamilyFromPostRequest().Name,
		}
		id := "61685179378d2ad5c3405bc5-x"
		expErr := errors.New(fmt.Sprintf(
			"family name: [%s] is duplicated",
			family.Name))
		repository := new(mocks.FamilyRepositoryMock)
		service := families.NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyFound := testdata.GetDefaultFamily()
		pagedResponseFamilies := entities.PagedResponse{
			Data: []entities.Family{
				{
					ID:   familyFound.ID,
					Name: familyFound.Name,
					Mccs: familyFound.Mccs,
				},
			},
		}
		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(pagedResponseFamilies, nil)

		err := service.Update(ctx, id, family)

		assert.NotNil(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		assert.Equal(t, exceptions.FamiliesNameDuplicated, err.(exceptions.DuplicatedException).Causes().Code)
		repository.AssertExpectations(t)
	})

	t.Run("when familyRepository on search found families, then try to get error existing family", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()

		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPutRequest()
		id := "61685179378d2ad5c3405bc5-x"
		expErr := errors.New(fmt.Sprintf(
			"there is more than one family with the name[%s] or the same mcc[%s]",
			family.Name,
			family.Mccs))
		repository := new(mocks.FamilyRepositoryMock)
		service := families.NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyRepositoryMock.On("Update", ctx, id, &family).
			Return(nil).Once()

		pagedResponseFamilies := entities.PagedResponse{
			Data: testdata.GetFamilies(),
		}
		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(pagedResponseFamilies, nil)

		err := service.Update(ctx, id, family)

		assert.NotNil(t, err)
		assert.Equal(t, expErr.Error(), err.Error())
		repository.AssertExpectations(t)
	})

	t.Run("when update is success", func(t *testing.T) {
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		ctx := context.TODO()
		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPutRequest()
		id := "61685179378d2ad5c3405bc5"

		service := families.NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))

		familyRepositoryMock.On("Update", ctx, id, &family).Return(nil).Once()

		familyRepositoryMock.On("SearchPaged", ctx, entities.NewDefaultPagination(), entities.FamilyFilter{
			Mccs: family.Mccs,
			Name: family.Name,
		}).Return(entities.PagedResponse{Data: []entities.Family{}}, nil)

		err := service.Update(ctx, id, family)

		assert.Nil(t, err)
		familyRepositoryMock.AssertExpectations(t)
	})
}

func Test_FamilyService_Delete(t *testing.T) {
	configs := config.NewConfig()
	logger, _ := logs.New()
	familyID := "611709bb70cbe3606baa3f8d"

	t.Run("test when delete family fail", func(t *testing.T) {
		expectedError := errors.New("error")

		rulesRepositoryMock := new(mocks.RulesRepositoryMock)
		rulesRepositoryMock.On("FindRulesPaged",
			context.TODO(),
			entities.RuleFilter{FamilyID: familyID},
			entities.Pagination{},
		).Return(entities.PagedResponse{
			Data: []entities.Rule{},
		}, nil)

		repository := new(mocks.FamilyRepositoryMock)
		repository.On("Delete", context.TODO(), familyID).Return(expectedError)
		service := families.NewFamilyService(configs, repository, rulesRepositoryMock, logger, new(datadog.MetricsDogMock))

		err := service.Delete(context.TODO(), familyID)

		assert.NotNil(t, err)
		assert.Equal(t, err, expectedError)
		repository.AssertExpectations(t)
	})

	t.Run("test when delete family fail looking associated rules", func(t *testing.T) {
		expectedError := fmt.Errorf("error loocking associated rules for familyID %s", familyID)

		rulesRepositoryMock := new(mocks.RulesRepositoryMock)
		rulesRepositoryMock.On("FindRulesPaged",
			context.TODO(),
			mock.AnythingOfType("entities.RuleFilter"),
			mock.AnythingOfType("entities.Pagination"),
		).Return(entities.PagedResponse{}, expectedError)

		service := families.NewFamilyService(configs, nil, rulesRepositoryMock, logger, new(datadog.MetricsDogMock))

		err := service.Delete(context.TODO(), familyID)

		assert.NotNil(t, err)
		assert.Equal(t, err, expectedError)
		rulesRepositoryMock.AssertExpectations(t)
	})

	t.Run("test when delete family successful", func(t *testing.T) {
		rulesRepositoryMock := new(mocks.RulesRepositoryMock)
		rulesRepositoryMock.On("FindRulesPaged",
			context.TODO(),
			entities.RuleFilter{FamilyID: familyID},
			entities.Pagination{},
		).Return(entities.PagedResponse{
			Data: []entities.Rule{},
		}, nil)

		repository := new(mocks.FamilyRepositoryMock)
		repository.On("Delete", context.TODO(), familyID).Return(nil)
		service := families.NewFamilyService(configs, repository, rulesRepositoryMock, logger, new(datadog.MetricsDogMock))

		err := service.Delete(context.TODO(), familyID)

		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("test when delete family is associated with an exiting rule", func(t *testing.T) {
		ruleAssociated := testdata.GetDefaultRuleWithFamilyMccID(false)

		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		rulesRepositoryMock := new(mocks.RulesRepositoryMock)
		rulesRepositoryMock.
			On("FindRulesPaged",
				context.TODO(),
				entities.RuleFilter{FamilyID: *ruleAssociated.FamilyMccID},
				entities.Pagination{},
			).
			Return(entities.PagedResponse{
				Data: []entities.Rule{ruleAssociated},
			}, nil)

		service := families.NewFamilyService(configs, familyRepositoryMock, rulesRepositoryMock, logger, new(datadog.MetricsDogMock))

		err := service.Delete(context.TODO(), *ruleAssociated.FamilyMccID)

		assert.Error(t, err)
		familyRepositoryMock.AssertExpectations(t)
	})
}

func TestFamilyServiceGetFamily(t *testing.T) {
	configs := config.NewConfig()
	logger, _ := logs.New()

	t.Run("GetFamily successfully", func(t *testing.T) {
		expectedFamilies := testdata.GetFamilies()
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		familiesFilter := entities.FamilyFilter{
			Mccs: expectedFamilies[0].Mccs,
		}
		familyRepositoryMock.
			On("SearchEvaluate",
				context.TODO(),
				familiesFilter,
			).
			Once().
			Return([]entities.Family{expectedFamilies[0]}, nil)

		service := families.NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))
		family, err := service.GetFamily(context.TODO(), familiesFilter)

		assert.NoError(t, err)
		assert.Equal(t, expectedFamilies[0].Name, family.Name)
		assert.Equal(t, expectedFamilies[0].ID, family.ID)
		assert.Equal(t, expectedFamilies[0].ExcludedCompanies, family.ExcludedCompanies)
		familyRepositoryMock.AssertExpectations(t)

	})
	t.Run("GetFamily with empty result", func(t *testing.T) {
		expectedFamilies := testdata.GetFamilies()
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		familiesFilter := entities.FamilyFilter{
			Mccs: expectedFamilies[0].Mccs,
		}
		familyRepositoryMock.
			On("SearchEvaluate",
				context.TODO(),
				familiesFilter,
			).
			Once().
			Return([]entities.Family{}, nil)

		service := families.NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))
		family, err := service.GetFamily(context.TODO(), familiesFilter)

		assert.NoError(t, err)
		assert.Equal(t, strings.Empty, family.Name)
		assert.Equal(t, 0, len(family.ExcludedCompanies))
		assert.Equal(t, 0, len(family.Mccs))
		familyRepositoryMock.AssertExpectations(t)
	})
	t.Run("GetFamily with an error", func(t *testing.T) {
		expectedFamilies := testdata.GetFamilies()
		familyRepositoryMock := new(mocks.FamilyRepositoryMock)
		familiesFilter := entities.FamilyFilter{
			Mccs: expectedFamilies[0].Mccs,
		}
		familyRepositoryMock.
			On("SearchEvaluate",
				context.TODO(),
				familiesFilter,
			).
			Once().
			Return([]entities.Family{}, fmt.Errorf("error getting family with mcc %s", expectedFamilies[0].Mccs))

		service := families.NewFamilyService(configs, familyRepositoryMock, nil, logger, new(datadog.MetricsDogMock))
		family, err := service.GetFamily(context.TODO(), familiesFilter)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("error getting family with mcc %s", expectedFamilies[0].Mccs))
		assert.Equal(t, strings.Empty, family.Name)
		assert.Equal(t, 0, len(family.ExcludedCompanies))
		assert.Equal(t, 0, len(family.Mccs))
		familyRepositoryMock.AssertExpectations(t)
	})
}
