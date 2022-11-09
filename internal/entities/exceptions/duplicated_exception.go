package exceptions

type DuplicatedException interface {
	Error() string
	IsDuplicatedError() bool
	Causes() Causes
}

type duplicatedException struct {
	ErrMessage string
	ErrCause   Causes
}

type Causes struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

func (exception *duplicatedException) Error() string {
	return exception.ErrMessage
}

func (exception *duplicatedException) Causes() Causes {
	return exception.ErrCause
}

func (exception *duplicatedException) IsDuplicatedError() bool {
	return true
}

func NewDuplicatedException(message string) DuplicatedException {
	return &duplicatedException{ErrMessage: message}
}

func NewDuplicatedExceptionWithCause(message string, causes Causes) DuplicatedException {
	return &duplicatedException{ErrMessage: message, ErrCause: causes}
}
