package errors

import "fmt"

type ApiError struct {
	Code    int
	Message string
	RawErr  error
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("API Error: [%d] %s", e.Code, e.Message)
}

func IsApiError(err error) bool {
	_, ok := err.(*ApiError)
	return ok
}
