package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Details interface{}
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewBatExecutionError(details interface{}, err error) *AppError {
	return &AppError{
		Code:    "BAT_EXECUTION_ERROR",
		Message: "Failed to execute batch file",
		Details: details,
		Err:     err,
	}
}

func NewValidationError(details interface{}) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: "Invalid input parameters",
		Details: details,
	}
}

func NewTimeoutError(details interface{}) *AppError {
	return &AppError{
		Code:    "TIMEOUT_ERROR",
		Message: "Operation timed out",
		Details: details,
	}
}
