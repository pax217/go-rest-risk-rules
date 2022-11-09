package testdata

import (
	"time"

	"github.com/conekta/risk-rules/internal/entities"
)

func GetDefaultRequestLists() entities.ListsSearch {
	return entities.ListsSearch{
		Email:     "me@gmail.com",
		CardHash:  "27603623df3a43698ac6f791448ce4ebdb046fcfbec2ce9f76dd0ce903799444",
		Phone:     "+5215555555555",
		CompanyID: "2",
	}
}

func GetDefaultWhiteList(isTest bool) entities.List {
	allowEmail := "me@gmail.com"
	return entities.List{
		Description: "migration",
		CompanyID:   "57ed5affdba34df0630071dc",
		Decision:    entities.Accepted,
		CreatedAt:   time.Date(2021, 12, 02, 20, 04, 35, 00, time.UTC),
		CreatedBy:   "riesgo@conekta.com",
		IsTest:      isTest,
		IsGlobal:    false,
		Value:       allowEmail,
		Field:       entities.EmailField,
		Rule:        "email == me@gmail.com",
		Type:        entities.White.String(),
	}
}

func GetDefaultBlackList(isTest bool) entities.List {
	email := "me@gmail.com"
	return entities.List{
		Description: "empty",
		CompanyID:   "2",
		Decision:    entities.Declined,
		CreatedAt:   time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC),
		CreatedBy:   "me",
		IsTest:      isTest,
		IsGlobal:    false,
		Value:       email,
		Field:       entities.EmailField,
		Rule:        "",
		Type:        entities.Black.String(),
	}
}

func GetDefaultGrayList(isTest bool) entities.List {
	email := "mail@hotmail.com"
	expires := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC)
	return entities.List{
		Description: "empty",
		CompanyID:   "2",
		Decision:    entities.Undecided,
		CreatedAt:   time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC),
		CreatedBy:   "me",
		IsTest:      isTest,
		IsGlobal:    true,
		Value:       email,
		Field:       entities.EmailField,
		Rule:        "",
		Type:        entities.Gray.String(),
		TimeToLive:  1,
		Expires:     &expires,
	}
}
