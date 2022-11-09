package acceptance

import (
	"os"
	"testing"

	"github.com/conekta/go_common/strings"
	"github.com/conekta/risk-rules/cmd/httpserver"
	"github.com/conekta/risk-rules/internal/container"
)

func TestMain(m *testing.M) {
	if !strings.IsEmpty(os.Getenv("ENV")) {
		m.Run()
	}
}
func GetServer() *httpserver.Server {
	dependencies := container.Build()

	server := httpserver.NewServer(dependencies)
	server.Validator()
	server.Routes()
	server.SetErrorHandler(httpserver.HTTPErrorHandler)
	return server
}
