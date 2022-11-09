package entities_test

import (
	"testing"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestList_IsEmpty(t *testing.T) {
	t.Run("Check list is empty", func(t *testing.T) {
		list := entities.List{}

		assert.True(t, list.IsEmpty())
	})
}

func TestList_IsWhitelist(t *testing.T) {
	t.Run("Check if is whitelist", func(t *testing.T) {
		list := testdata.GetDefaultWhiteList(true)

		assert.True(t, list.IsWhitelist())
	})
}

func TestList_IsBlacklist(t *testing.T) {
	t.Run("Check if is blacklist", func(t *testing.T) {
		list := testdata.GetDefaultBlackList(true)

		assert.True(t, list.IsBlacklist())
	})
}

func TestList_IsGraylist(t *testing.T) {
	t.Run("Check if is graylist", func(t *testing.T) {
		list := testdata.GetDefaultGrayList(true)

		assert.True(t, list.IsGraylist())
		assert.True(t, list.IsValidListType())
	})
}

func TestList_IsListResponseEmpty(t *testing.T) {
	t.Run("Check if is empty ListResponse", func(t *testing.T) {
		listResponse := entities.ListResponse{}
		assert.True(t, listResponse.IsListResponseEmpty())
	})
}

func TestList_NewListResponse(t *testing.T) {
	t.Run("Check if is empty List Response", func(t *testing.T) {
		listResponse := entities.NewListResponse()
		decision := listResponse.GetDecision()

		assert.Empty(t, listResponse.GetEntityName())
		assert.Equal(t, entities.Undecided, decision)
		assert.True(t, listResponse.IsListResponseEmpty())
	})
}
