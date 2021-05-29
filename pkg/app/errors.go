package app

// UnauthorizedError principle does not have access
type UnauthorizedError struct {
	s string
}

// Error interface
func (e *UnauthorizedError) Error() string {
	return e.s
}

// Unauthorized always returns true for this error type. Can be used by the caller for error checking by behavior
func (e *UnauthorizedError) Unauthorized() bool {
	return true
}

// NotFoundError not found in store
type NotFoundError struct {
	s string
}

// Error interface
func (e *NotFoundError) Error() string {
	return e.s
}

// NotFound always returns true for this error type. Can be used by the caller for error checking by behavior
func (e *NotFoundError) NotFound() bool {
	return true
}
