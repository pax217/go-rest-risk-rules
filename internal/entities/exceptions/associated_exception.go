package exceptions

type AssociatedException interface {
	Error() string
	IsAssociatedError() bool
	Causes() Causes
}

type associatedException struct {
	ErrMessage string
	ErrCause   Causes
}

func (exception *associatedException) Error() string {
	return exception.ErrMessage
}

func (exception *associatedException) Causes() Causes {
	return exception.ErrCause
}

func (exception *associatedException) IsAssociatedError() bool {
	return true
}

func NewAssociatedExceptionWithCause(message string, causes Causes) AssociatedException {
	return &associatedException{ErrMessage: message, ErrCause: causes}
}
