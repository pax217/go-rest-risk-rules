package acceptance

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conekta/risk-rules/internal/container"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {

	t.Run("when evaluation is successful", func(t *testing.T) {

		cfg := container.Build().Config
		mongoDB := mongodb.NewMongoDB(cfg)
		declinedRule := testdata.GetDefaultRule(false)
		approvedRule := testdata.GetDefaultRuleWithApprovedDecision(false)
		mongoDB.PrepareData(context.Background(), cfg.MongoDB.Collections.Rules, declinedRule, approvedRule)
		defer mongoDB.ClearCollection(context.Background(), cfg.MongoDB.Collections.Rules)
		defer mongoDB.ClearCollection(context.Background(), cfg.MongoDB.Collections.ChargeEvaluations)

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(testdata.GetDefaultCharge())
		if err != nil {
			assert.Fail(t, err.Error())
		}

		server := GetServer()
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/risk-rules/v1/charges/evaluate", &buf)
		req.Header.Set("Content-Type", "application/json")
		server.Server.ServeHTTP(w, req)

		assert.NoError(t, err, "err should be null when calling to charge evaluate endpoint")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Body.String())

		var response entities.EvaluationResponse

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "no error in Unmarshal body")

		assert.Equal(t, "A", response.Decision, "decision should be A because we have an accepted rule")
		assert.True(t, len(response.Modules.Rules.DecisionRules) > 0)
		assert.Containsf(t, getRulesID(response.Modules.Rules.DecisionRules), approvedRule.ID.String(),
			"%s should be included in decision rules", approvedRule.ID.String())

	})

}

func getRulesID(rules []entities.Rule) []string {

	var list []string
	for i := range rules {
		list = append(list, rules[i].ID.String())
	}
	return list

}
