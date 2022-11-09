package echo

import (
	"fmt"
	"strings"

	"net/http/httptest"

	"github.com/conekta/risk-rules/cmd/httpserver"
	"github.com/conekta/risk-rules/internal/container"
	str "github.com/conekta/risk-rules/pkg/strings"
	"github.com/labstack/echo/v4"
)

func SetupAsRecorder(method, target, id, body string) (echo.Context, *httptest.ResponseRecorder) {
	mockServer := httpserver.NewServer(container.Dependencies{})
	mockServer.SetErrorHandler(httpserver.HTTPErrorHandler)
	mockServer.Validator()

	request := httptest.NewRequest(method, fmt.Sprintf("%s%s", target, id), strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	context := mockServer.NewServerContext(request, recorder)
	context.SetPath(target)

	if !str.IsEmpty(id) {
		context.SetPath(fmt.Sprintf("%s%s", target, "/:id"))
		context.SetParamNames("id")
		context.SetParamValues(id)
	}

	return context, recorder
}
