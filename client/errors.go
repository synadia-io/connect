package client

import (
	"errors"
	"fmt"
)

const GenericErrMsg = "an error occurred"

type ConnectServiceError struct {
	Code          int    // nats service error code
	Description   string // friendly & safe error description
	internalError error
	id            string
}

// Error implements the builtin/error interface and returns the full internal error
func (ce ConnectServiceError) Error() string {
	return ce.internalError.Error()
}

// Body implements ClientError interface (control-plane) and returns a friendly & sanitized error
func (ce ConnectServiceError) Body() string {
	return fmt.Sprintf("connect[%d]: %s", ce.Code, ce.Description)
}

// ID implements IdableClientError interface (control-plane)
func (ce ConnectServiceError) ID() string {
	return ce.id
}

// ErrorMessage returns the user-friendly message from an error.
// If the error is a ConnectServiceError, returns Body().
// Otherwise returns err.Error().
func ErrorMessage(err error) string {
	var svcErr *ConnectServiceError
	if errors.As(err, &svcErr) {
		return svcErr.Body()
	}
	return err.Error()
}
