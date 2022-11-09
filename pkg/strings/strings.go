package strings

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

const (
	HexLen  = 24
	Empty   = ""
	bitSize = 64
)

var ErrInvalidHex = errors.New("the provided hex string is not a valid")

func ToInt64(value string, defaultValue int64) (int64, error) {
	if value == "" {
		return defaultValue, nil
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return int64(n), nil
}

func IsEmpty(value string) bool {
	return value == ""
}

func IsStringPointerEmpty(value *string) bool {
	return value == nil || IsEmpty(*value)
}

func StringToStringPointer(value string) *string {
	if IsEmpty(value) {
		return nil
	}
	return &value
}

func StringPointerToString(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}

func StringArrPointerToStringArr(value *[]string) []string {
	res := make([]string, 0)
	if value == nil || len(*value) == 0 {
		return res
	}
	return *value
}

func IsBoolean(value string) bool {
	if value == "true" || value == "false" {
		return true
	}
	return false
}

func ContainsAnyString(value string, lists ...string) bool {
	for _, list := range lists {
		if strings.Contains(value, list) {
			return true
		}
	}
	return false
}

func IsHex(s string) error {
	if len(s) != HexLen {
		return ErrInvalidHex
	}

	_, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	return nil
}

func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func ToFloat64(value string, defaultValue float64) (float64, error) {
	if IsEmpty(value) {
		return defaultValue, nil
	}

	n, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func ToReader(s string) *strings.Reader {
	return strings.NewReader(s)
}
