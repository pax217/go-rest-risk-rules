package exceptions

type InvalidRequestException interface {
	Error() string
	IsInvalidRequestException() bool
	Causes() Causes
}

type invalidRequestException struct {
	ErrMessage string
	ErrCause   Causes
}

func (exception *invalidRequestException) Error() string {
	return exception.ErrMessage
}

func (exception *invalidRequestException) IsInvalidRequestException() bool {
	return true
}
func (exception *invalidRequestException) Causes() Causes {
	return exception.ErrCause
}
func NewInvalidRequest(msg string) InvalidRequestException {
	return &invalidRequestException{ErrMessage: msg}
}
func NewInvalidRequestWithCauses(msg string, causes Causes) InvalidRequestException {
	return &invalidRequestException{ErrMessage: msg, ErrCause: causes}
}
