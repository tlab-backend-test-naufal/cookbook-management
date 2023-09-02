package liberr

type ErrorDetails struct {
	// Message (required) is the user-defined error message.
	// E.g. "user email has invalid format".
	Message string

	// Code (required) is the user-defined error code string
	// E.g. "cookbook-MANAGEMENT_cookbook_NOT-FOUND".
	Code string
}

// NewErrorDetails creates a new ErrorDetails struct with the given parameters.
func NewErrorDetails(code, message string) *ErrorDetails {
	return &ErrorDetails{
		Message: message,
		Code:    code,
	}
}

// Error() is used to implement the Golang `error` interface.
func (e *ErrorDetails) Error() string {
	return e.Message
}
