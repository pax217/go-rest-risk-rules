package httpserver

import (
	"github.com/conekta/go_common/http/resterror"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/conekta/risk-rules/internal/container"

	"github.com/conekta/risk-rules/internal/config"
)

func TestWithRequestHeaderValidator_WhenTokensAreNotEquals_ShouldBeFail(t *testing.T) {
	dependencies := container.Dependencies{}
	server := NewServer(dependencies)
	os.Setenv("REQUEST_HEADER_TOKEN", "123456789")
	defer os.Clearenv()
	cfg := config.NewConfig()
	handlerToTest := WithRequestHeaderValidator(cfg)
	server.Middlewares(handlerToTest)
	req := httptest.NewRequest("GET", "/risk-rules/v1/rules", nil)
	req.Header.Set("X-Request-Risk", "12345678")
	resp := httptest.NewRecorder()
	server.ServerHTTP(resp, req)

	assert.Equal(t, 403, resp.Code)
}
func TestWithRequestHeaderValidator2(t *testing.T) {
	dependencies := container.Dependencies{}
	server := NewServer(dependencies)
	os.Setenv("REQUEST_HEADER_TOKEN", "123456789")
	defer os.Clearenv()
	cfg := config.NewConfig()
	handlerToTest := WithRequestHeaderValidator(cfg)
	server.Middlewares(handlerToTest)
	req := httptest.NewRequest("GET", "/risk-rules/v1/rules", nil)
	resp := httptest.NewRecorder()
	server.ServerHTTP(resp, req)

	assert.Equal(t, 403, resp.Code)
}

func TestWithRequestHeaderValidator_TokenIsEmpty_ShouldFail(t *testing.T) {
	dependencies := container.Dependencies{}
	server := NewServer(dependencies)
	os.Setenv("REQUEST_HEADER_TOKEN", "")
	defer os.Clearenv()
	cfg := config.NewConfig()
	handlerToTest := WithRequestHeaderValidator(cfg)
	server.Middlewares(handlerToTest)
	req := httptest.NewRequest("GET", "/risk-rules/v1/rules", nil)
	resp := httptest.NewRecorder()
	server.ServerHTTP(resp, req)

	assert.Equal(t, 403, resp.Code)
}

func TestWithRequestHeaderValidator_TokenIsNotPresent_ShouldFail(t *testing.T) {
	dependencies := container.Dependencies{}
	server := NewServer(dependencies)
	defer os.Clearenv()
	cfg := config.NewConfig()
	handlerToTest := WithRequestHeaderValidator(cfg)
	server.Middlewares(handlerToTest)
	req := httptest.NewRequest("GET", "/risk-rules/v1/rules", nil)
	resp := httptest.NewRecorder()
	server.ServerHTTP(resp, req)

	assert.Equal(t, 403, resp.Code)
}

func TestHTTPErrorHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/risk-rules/v1/rules", nil)
	resp := httptest.NewRecorder()

	t.Run("exception InvalidRequest", func(t *testing.T) {
		err := exceptions.NewInvalidRequest("reason")
		ctx := echo.New().NewContext(req, resp)

		HTTPErrorHandler(err, ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	})
	t.Run("exception AssociatedException", func(t *testing.T) {
		err := exceptions.NewAssociatedExceptionWithCause("reason", exceptions.Causes{
			Code:    "002",
			Message: "causes",
		})
		ctx := echo.New().NewContext(req, resp)

		HTTPErrorHandler(err, ctx)

		assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	})

	t.Run("exception DuplicatedException", func(t *testing.T) {
		err := exceptions.NewDuplicatedException("reason")
		ctx := echo.New().NewContext(req, resp)

		HTTPErrorHandler(err, ctx)

		assert.Equal(t, http.StatusConflict, ctx.Response().Status)
	})
	t.Run("exception NotFoundException", func(t *testing.T) {
		err := exceptions.NewNotFoundException("reason")
		ctx := echo.New().NewContext(req, resp)

		HTTPErrorHandler(err, ctx)

		assert.Equal(t, http.StatusNotFound, ctx.Response().Status)
	})

	t.Run("exception RestErr", func(t *testing.T) {
		err := resterror.NewUnauthorizedError("unauthorized")
		ctx := echo.New().NewContext(req, resp)

		HTTPErrorHandler(err, ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Response().Status)
	})

}
