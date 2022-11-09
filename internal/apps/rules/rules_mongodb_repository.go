package rules

import (
	"context"
	"errors"
	"fmt"

	"github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const repositoryName = "rules.repository.mongo.%s"

type RuleRepository interface {
	AddRule(ctx context.Context, request entities.Rule) (entities.Rule, error)
	UpdateRule(ctx context.Context, ruleID string, ruleReq entities.Rule) error
	RemoveRule(ctx context.Context, ruleID string) error
	FindRulesPaged(ctx context.Context, filter entities.RuleFilter, pagination entities.Pagination) (entities.PagedResponse, error)
	GetRulesByFilters(ctx context.Context, filter entities.RuleFilter, component entities.ConsoleComponent) ([]entities.Rule, error)
}

type RuleMongoDBRepository struct {
	config  config.Config
	mongodb mongodb.MongoDBier
	log     logs.Logger
}

func NewRuleMongoDBRepository(cfg config.Config, mongoDBier mongodb.MongoDBier, logger logs.Logger) RuleRepository {
	return &RuleMongoDBRepository{
		config:  cfg,
		mongodb: mongoDBier,
		log:     logger,
	}
}

func (r *RuleMongoDBRepository) AddRule(ctx context.Context, rule entities.Rule) (entities.Rule, error) {
	rulesCollection := r.mongodb.Collection(r.config.MongoDB.Collections.Rules)
	result, err := rulesCollection.InsertOne(ctx, rule, nil)
	if err != nil {
		r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "AddRule"))
		return rule, err
	}

	rule.ID = result.InsertedID.(primitive.ObjectID)
	return rule, nil
}

func (r *RuleMongoDBRepository) UpdateRule(ctx context.Context, ruleID string, rule entities.Rule) error {
	_id, err := primitive.ObjectIDFromHex(ruleID)
	if err != nil {
		r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "UpdateRule"))
		return err
	}

	rulesCollection := r.mongodb.Collection(r.config.MongoDB.Collections.Rules)
	filter := bson.M{"_id": _id}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: "updated_by", Value: rule.UpdatedBy},
				primitive.E{Key: "updated_at", Value: rule.UpdatedAt},
				primitive.E{Key: "is_test", Value: rule.IsTest},
				primitive.E{Key: "module", Value: rule.Module},
				primitive.E{Key: "is_global", Value: rule.IsGlobal},
				primitive.E{Key: "description", Value: rule.Description},
				primitive.E{Key: "rules", Value: rule.Rules},
				primitive.E{Key: "rule", Value: rule.Rule},
				primitive.E{Key: "company_id", Value: rule.CompanyID},
				primitive.E{Key: "family_company_id", Value: rule.FamilyCompanyID},
				primitive.E{Key: "family_id", Value: rule.FamilyMccID},
				primitive.E{Key: "decision", Value: rule.Decision},
				primitive.E{Key: "is_yellow_flag", Value: rule.IsYellowFlag},
			},
		},
	}

	result, err := rulesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "UpdateRule"))
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("error, document not found")
	}

	return nil
}

func (r *RuleMongoDBRepository) RemoveRule(ctx context.Context, ruleID string) error {
	ID, err := primitive.ObjectIDFromHex(ruleID)
	if err != nil {
		r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "RemoveRule"))
		return err
	}

	rulesCollection := r.mongodb.Collection(r.config.MongoDB.Collections.Rules)
	filter := bson.M{"_id": ID}
	result, err := rulesCollection.DeleteOne(ctx, filter, nil)
	if err != nil {
		r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "RemoveRule"))
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("error: record not found")
	}

	return nil
}

func (r *RuleMongoDBRepository) GetRulesByFilters(ctx context.Context, filter entities.RuleFilter,
	component entities.ConsoleComponent) ([]entities.Rule, error) {
	collection := r.mongodb.Collection(r.config.MongoDB.Collections.Rules)

	query := buildRulesFilter(component, filter)

	rules := make([]entities.Rule, 0)
	cur, err := collection.Find(ctx, query)
	if err != nil {
		r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetRulesByFilters"),
			text.CompanyID, filter.CompanyID, text.FamilyID, filter.FamilyID)
		return nil, err
	}

	err = cur.All(ctx, &rules)
	if err != nil {
		r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetRulesByFilters"),
			text.CompanyID, filter.CompanyID, text.FamilyID, filter.FamilyID)
		return nil, err
	}

	return rules, nil
}

func (r *RuleMongoDBRepository) FindRulesPaged(ctx context.Context, filter entities.RuleFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	collection := r.mongodb.Collection(r.config.MongoDB.Collections.Rules)
	query := bson.D{}

	if !strings.IsEmpty(filter.ID) {
		ID, err := primitive.ObjectIDFromHex(filter.ID)
		if err != nil {
			r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "FindRulesPaged"))
			return entities.PagedResponse{}, err
		}

		query = append(query, primitive.E{Key: "_id", Value: ID})
	}

	if !filter.IsEmptyCompanyID() {
		query = append(query, bson.E{Key: "company_id", Value: filter.CompanyID})
	}

	if !filter.IsEmptyFamilyID() {
		query = append(query, bson.E{Key: "family_id", Value: filter.FamilyID})
	}

	if !strings.IsEmpty(filter.Rule) {
		query = append(query, bson.E{Key: "rule", Value: filter.Rule})
	}

	if !filter.IsEmptyCompanyID() && !filter.IsEmptyFamilyID() {
		query = append(query, bson.E{Key: "is_global", Value: true})
	}

	if !filter.IsEmptyFamilyCompaniesID() {
		query = append(query, bson.E{Key: "family_company_id", Value: filter.FamilyCompanyID})
	}

	total, _ := collection.CountDocuments(ctx, query)
	hasMore := pagination.HasMorePages(total)

	opts := options.FindOptions{}
	opts.SetLimit(pagination.PageSize)
	opts.SetSkip(pagination.GetPageStartIndex())
	opts.SetSort(bson.D{primitive.E{Key: "created_at", Value: -1}})

	rules := make([]entities.Rule, 0)
	cur, err := collection.Find(ctx, query, &opts)
	if err != nil {
		r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "FindRulesPaged"),
			text.List, filter)
		return entities.PagedResponse{}, err
	}

	err = cur.All(ctx, &rules)
	if err != nil {
		r.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "FindRulesPaged"),
			text.List, filter)
		return entities.PagedResponse{}, err
	}

	return entities.NewPagedResponse(rules, hasMore, total), nil
}

func buildRulesFilter(component entities.ConsoleComponent, filter entities.RuleFilter) bson.M {
	var query []bson.M
	findQuery := bson.M{}

	if component == entities.CompanyRulesType {
		query = append(query,
			bson.M{"company_id": filter.CompanyID})
	}

	if component == entities.FamilyCompanyRulesType {
		query = append(query, bson.M{"family_id": filter.FamilyID})
	}

	if component == entities.FamilyMccRulesType {
		query = append(query, bson.M{"family_company_id": bson.M{"$in": filter.FamilyCompaniesIDs}})
	}

	if component == entities.GlobalRulesType {
		query = append(query, bson.M{"is_global": true})
	}

	if component == entities.YellowFlagType {
		query = append(query, bson.M{"is_yellow_flag": true})
	}

	if component == entities.IdentityModuleType {
		return buildRulesIdentityModule(filter)
	}

	query = append(query, bson.M{"rules.field": bson.M{"$not": bson.M{"$regex": `email_proximity.*`}}})

	findQuery["$and"] = query

	return findQuery
}

func buildRulesIdentityModule(filter entities.RuleFilter) bson.M {
	var query []bson.M
	findQuery := bson.M{}

	if !strings.IsEmpty(filter.CompanyID) {
		query = append(query, bson.M{"company_id": filter.CompanyID})
	}

	if !strings.IsEmpty(filter.FamilyID) {
		query = append(query, bson.M{"family_id": filter.FamilyID})
	}

	if len(filter.FamilyCompaniesIDs) > 0 {
		query = append(query, bson.M{"family_company_id": bson.M{"$in": filter.FamilyCompaniesIDs}})
	}

	if filter.IsGlobal {
		query = append(query, bson.M{"is_global": true})
	}

	query = append(query, bson.M{"rules.field": bson.M{"$regex": `email_proximity.*`}})

	findQuery["$and"] = query

	return findQuery
}
