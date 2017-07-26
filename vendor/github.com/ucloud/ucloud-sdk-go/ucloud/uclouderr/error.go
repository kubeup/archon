package uclouderr

import "fmt"

type UcloudError struct {
	// Return the Code of request
	RetCode int
	// Returns the error details message.
	Message string
}

type RequestFailed struct {
	UcloudError
	StatusCode int
}

func (r *RequestFailed) Error() string {
	return fmt.Sprintf("Request Error, status code: %d, return code: %d, message: %s",
		r.StatusCode, r.RetCode, r.Message)
}
