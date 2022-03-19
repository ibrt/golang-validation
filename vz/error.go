package vz

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ibrt/golang-errors/errorz"
)

const (
	// ErrIDFailedValidation is an error ID.
	ErrIDFailedValidation = errorz.ID("failed-validation")
)

// NewErrFailedValidation creates a new failed validation error.
func NewErrFailedValidation(format string, fields map[string]interface{}, options ...errorz.Option) error {
	return errorz.Errorf(format, append(options,
		ErrIDFailedValidation,
		errorz.Status(http.StatusBadRequest),
		errorz.Prefix("failed validation"),
		newErrFailedValidationOption(fields),
		errorz.Skip())...)
}

// WrapErrFailedValidation wraps an error as failed validation error.
func WrapErrFailedValidation(err error, options ...errorz.Option) error {
	if errorz.GetID(err) == ErrIDFailedValidation {
		return errorz.Wrap(err, options...)
	}

	return errorz.Wrap(err, append(options,
		ErrIDFailedValidation,
		errorz.Status(http.StatusBadRequest),
		errorz.Prefix("failed validation"),
		newErrFailedValidationOption(nil),
		errorz.Skip())...)
}

func newErrFailedValidationOption(extraFields map[string]interface{}) errorz.Option {
	return errorz.OptionFunc(func(err error) {
		finalFields := make(map[string]interface{})

		if errs, ok := errorz.Unwrap(err).(validator.ValidationErrors); ok {
			for k, v := range fieldsFromValidationErrors(errs) {
				finalFields[k] = v
			}
		}

		for k, v := range extraFields {
			finalFields[k] = v
		}

		errorz.M("fields", finalFields)(err)
	})
}

func fieldsFromValidationErrors(errs validator.ValidationErrors) map[string]interface{} {
	fields := map[string]interface{}{}

	for _, err := range errs {
		n := err.Namespace()
		if i := strings.Index(err.Namespace(), "."); i >= 0 {
			n = err.Namespace()[i+1:]
		}
		fields[n] = err.Tag()
	}

	return fields
}
