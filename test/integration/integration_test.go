package integration

import (
	"github.com/conekta/risk-rules/pkg/strings"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if !strings.IsEmpty(os.Getenv("ENV")) {
		m.Run()
	}
}
