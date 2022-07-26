package exceptions

type ResourceNotFoundError struct {
	Message string
}

func (e ResourceNotFoundError) Error() string {
	return e.Message
}

type ResourceAlreadyExistsError struct {
	Message string
}

func (e ResourceAlreadyExistsError) Error() string {
	return e.Message
}

type InvalidRequestError struct {
	Message string
}

func (e InvalidRequestError) Error() string {
	return e.Message
}

type NotAuthorizedError struct {
	Message string
}

func (e NotAuthorizedError) Error() string {
	return e.Message
}
