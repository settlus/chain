package types

type NotFoundError struct {
	msg string
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{
		msg: msg,
	}
}

func (e *NotFoundError) Error() string { return e.msg }
