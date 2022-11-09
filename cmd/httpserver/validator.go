package httpserver

import (
	"github.com/go-playground/validator/v10"
)

type APIValidator struct {
	validator *validator.Validate
}

func (cv *APIValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func (s *Server) Validator() {
	s.Server.Validator = &APIValidator{validator: validator.New()}
}
