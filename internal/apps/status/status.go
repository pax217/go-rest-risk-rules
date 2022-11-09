package status

import (
	"net/http"
	"time"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/labstack/echo/v4"
)

type StatusHandler interface {
	Ping(c echo.Context) error
}

type Response struct {
	Version string    `json:"version"`
	Name    string    `json:"name"`
	Uptime  time.Time `json:"uptime"`
}

type statusHandler struct {
	config config.Config
	metric datadog.Metricer
}

func NewStatusHandler(cfg config.Config, metric datadog.Metricer) StatusHandler {
	return &statusHandler{
		config: cfg,
		metric: metric,
	}
}

func (c *statusHandler) Ping(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, Response{
		Version: c.config.ProjectVersion,
		Name:    c.config.ProjectName,
		Uptime:  time.Now().UTC(),
	})
}
