package integration

import (
	"context"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/rest"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestRkListsRest_ListsSearch(t *testing.T) {
	t.Run("when rk lists service responds successfully", func(t *testing.T) {

		expected := testdata.GetDefaultWhiteList(true)
		logger, _ := logs.New()
		cfg := config.NewConfig()

		listsSearch := testdata.GetDefaultRequestLists()
		restClient := rest.NewRkListsRestClient(cfg, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foundLists, err := restClient.ListsSearch(ctx, listsSearch)

		assert.Nil(t, err)
		assert.Equal(t, []entities.List{expected}, foundLists)
	})
}
