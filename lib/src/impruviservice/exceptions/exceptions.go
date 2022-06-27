package exceptions

type ResourceNotFoundError struct {
	Message string
}

func (e ResourceNotFoundError) Error() string {
	return e.Message
}

type InvalidRequestError struct {
	message string
}

func (e *InvalidRequestError) Error() string {
	return e.message
}
