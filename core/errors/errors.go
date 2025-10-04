package errors

type TokenRecievingError struct {
	msg string
}

func NewTockenRecievingError(msg string) TokenRecievingError {
	return TokenRecievingError{
		msg: msg,
	}
}

func (e TokenRecievingError) Error() string {
	return e.msg
}
