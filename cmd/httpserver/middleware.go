package httpserver

import (
	"fmt"
	"strings"

	"github.com/conekta/go_common/logs"

	customString "github.com/conekta/risk-rules/pkg/strings"

	"net/http"

	"github.com/conekta/go_common/http/resterror"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

var (
	excludedGzipPaths = []string{"docs", "metrics"}

	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = echoMiddleware.RequestIDConfig{
		Skipper: echoMiddleware.DefaultSkipper,
		Generator: func() string {
			return uuid.New().String()
		},
	}
	XRequestRisk = "X-Request-Risk"
)

type Middleware func(*Server)

// Middlewares build the middlewares of the server
func (s *Server) Middlewares(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		middleware(s)
	}
}

func WithLogger(cfg config.Config) Middleware {
	return func(s *Server) {
		// TODO: implement conekta standard
		s.Server.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Skipper: func(e echo.Context) bool {
				return strings.Contains(e.Path(), "ping")
			},
			CustomTimeFormat: "2006-01-02T15:04:05.1483386-00:00",
			Format: `{"time":"${time_custom}","level":"Info","error":"${error}" ,"method":"${method}","uri":"${uri}",` +
				fmt.Sprintf(`"status":"${status}", "origin":"${header:X-Application-ID}" , "service": %q }`,
					cfg.ProjectName) + "\n",
		}))
	}
}

func WithGzip() Middleware {
	return func(s *Server) {
		s.Server.Use(echoMiddleware.GzipWithConfig(echoMiddleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				for _, path := range excludedGzipPaths {
					if strings.Contains(c.Request().URL.Path, path) {
						return true
					}
				}
				return false
			},
		}))
	}
}

func WithRecover() Middleware {
	return func(s *Server) {
		s.Server.Use(echoMiddleware.Recover())
	}
}

func WithAPM() Middleware {
	return func(s *Server) {
		s.Server.Use(
			echotrace.Middleware(
				echotrace.WithServiceName(s.dependencies.Config.ProjectName),
			),
		)
	}
}

func WithRequestID() Middleware {
	return func(s *Server) {
		s.Server.Use(
			echoMiddleware.RequestIDWithConfig(DefaultRequestIDConfig),
		)
	}
}

func WithRequestHeaderValidator(cfg config.Config) Middleware {
	return func(s *Server) {
		s.Server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if customString.ContainsAnyString(c.Path(), "ping", "health") {
					return next(c)
				}
				token := c.Request().Header.Get(XRequestRisk)
				if customString.IsEmpty(token) || !strings.EqualFold(token, cfg.RequestHeaderToken) {
					err := resterror.NewRestError(http.StatusText(http.StatusForbidden), http.StatusForbidden,
						http.StatusText(http.StatusForbidden), []interface{}{http.StatusText(http.StatusForbidden)})
					return c.JSON(err.Status(), err)
				}
				return next(c)
			}
		})
	}
}

func WithCORS() Middleware {
	return func(s *Server) {
		s.Server.Use(
			echoMiddleware.CORS(),
		)
	}
}

func WithLoggerBody(logger logs.Logger) Middleware {
	return func(s *Server) {
		s.Server.Use(
			echoMiddleware.BodyDumpWithConfig(echoMiddleware.BodyDumpConfig{
				Skipper: func(ctx echo.Context) bool {
					if customString.ContainsAnyString(ctx.Path(), "ping", "health") {
						return true
					}
					if !strings.EqualFold(ctx.Request().Method, http.MethodPost) {
						return true
					}

					return false
				},
				Handler: func(context echo.Context, reqBody, resBody []byte) {
					logger.Info(context.Request().Context(), fmt.Sprintf("[logger-middleware] path:%s body:%s",
						context.Path(), string(reqBody)))
				},
			}))
	}
}

func HTTPErrorHandler(err error, ctx echo.Context) {
	var apiError resterror.RestErr

	switch value := err.(type) {
	case exceptions.DuplicatedException:
		apiError = resterror.NewRestError(err.Error(), http.StatusConflict, "conflict",
			[]interface{}{value.Causes()})
	case resterror.RestErr:
		apiError = value
	case exceptions.NotFoundException:
		apiError = resterror.NewNotFoundError(err.Error())
	case exceptions.InvalidRequestException:
		apiError = resterror.NewRestError(err.Error(), http.StatusBadRequest, http.StatusText(http.StatusBadRequest),
			[]interface{}{value.Causes()})
	case exceptions.AssociatedException:
		apiError = resterror.NewRestError(err.Error(), http.StatusBadRequest, "Bad Request",
			[]interface{}{value.Causes()})
	default:
		apiError = resterror.NewInternalServerError(err.Error(), err)
	}

	ctx.JSON(apiError.Status(), apiError)
}
