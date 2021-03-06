package errors

import (
	"net/http"

	"github.com/ingot-cloud/ingot-go/pkg/framework/core/code"

	"github.com/pkg/errors"
)

var (
	// Wrap errors.Wrap
	Wrap = errors.Wrap
	// Wrapf errors.Wrapf
	Wrapf = errors.Wrapf
	// WithStack errors.WithStack
	WithStack = errors.WithStack
	// WithMessage errors.WithMessage
	WithMessage = errors.WithMessage
	// WithMessagef errors.WithMessagef
	WithMessagef = errors.WithMessagef
)

// New error
func New(statusCode int, code string, message string) error {
	return &E{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// Unpack error
func Unpack(err error) *E {
	if e, ok := err.(*E); ok {
		return e
	}

	message := "Unknow"
	if err != nil {
		message = err.Error()
	}

	return &E{
		StatusCode: http.StatusInternalServerError,
		Code:       code.InternalServerError,
		Message:    message,
	}
}

// E is error wrapper
type E struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *E) Error() string {
	return e.Message
}
