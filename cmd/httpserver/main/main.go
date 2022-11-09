package main

import (
	"context"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/risk-rules/cmd/httpserver"
	"github.com/conekta/risk-rules/internal/container"
)

func main() {
	dependencies := container.Build()
	apm, err := datadog.NewAPM(
		datadog.WithEnv(dependencies.Config.Env),
		datadog.WithServiceVersion(dependencies.Config.ProjectVersion),
		datadog.WithService(dependencies.Config.ProjectName),
		datadog.WithRuntimeMetrics(),
	)
	if err != nil {
		dependencies.Logs.Fatal(context.Background(), err.Error())
	}
	if apm != nil {
		defer apm.Stop()
	}
	apm.Start()

	server := httpserver.NewServer(dependencies)
	server.Middlewares(httpserver.WithGzip(), httpserver.WithAPM(), httpserver.WithRecover(),
		httpserver.WithLogger(dependencies.Config), httpserver.WithRequestID(), httpserver.WithCORS(),
		httpserver.WithRequestHeaderValidator(dependencies.Config), httpserver.WithLoggerBody(dependencies.Logs))
	server.Validator()
	server.Routes()

	go dependencies.ChargebacksHandler.ListenChargebacks()

	server.SetErrorHandler(httpserver.HTTPErrorHandler)
	server.Start()
}
