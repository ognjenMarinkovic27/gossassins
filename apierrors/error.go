package apierrors

type StatusError interface {
	Status() int
	error
}

type statusError struct {
	status int
	err    error
}

func NewStatusError(status int, err error) StatusError {
	return statusError{status, err}
}

func (e statusError) Error() string {
	return e.err.Error()
}

func (e statusError) Status() int {
	return e.status
}
