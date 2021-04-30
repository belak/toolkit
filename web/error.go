package web

type StatusError interface {
	error
	Status() int
}

type statusError struct {
	error
	status int
}

func NewError(statusCode int, err error) StatusError {
	if err == nil {
		return nil
	}

	return statusError{
		err,
		statusCode,
	}
}

func (e statusError) Status() int {
	return e.status
}

func (e statusError) Unwrap() error {
	return e.error
}
