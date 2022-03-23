package vz

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ibrt/golang-errors/errorz"
)

var (
	validate = validator.New()
)

func init() {
	validate.RegisterTagNameFunc(validatorTagName)
}

func validatorTagName(f reflect.StructField) string {
	if name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]; name != "" && name != "-" {
		return name
	}

	return f.Name
}

// ValidateStruct validates a struct.
func ValidateStruct(v interface{}) error {
	if err := validate.Struct(v); err != nil {
		return WrapErrFailedValidation(err, errorz.SkipPackage())
	}
	return nil
}

// MustValidateStruct is like ValidateStruct, panics on error.
func MustValidateStruct(v interface{}) {
	errorz.MaybeMustWrap(ValidateStruct(v), errorz.SkipPackage())
}

// RegexpValidatorFactory creates a validator that matches a regexp.
func RegexpValidatorFactory(regexp *regexp.Regexp) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return regexp.MatchString(fl.Field().String())
	}
}

// MustRegisterValidator registers a validator.
func MustRegisterValidator(tag string, validator validator.Func) {
	errorz.MaybeMustWrap(validate.RegisterValidation(tag, validator), errorz.SkipPackage())
}

// SimpleValidator describes a type that can validate itself, returning only a "valid" bool.
type SimpleValidator interface {
	Valid() bool
}

// Validator describes a type that can validate itself.
type Validator interface {
	Validate() error
}

// IsValidatable returns true if the given struct implements at least one of Validator or SimpleValidator.
func IsValidatable(v interface{}) bool {
	_, isValidator := v.(Validator)
	_, isSimpleValidator := v.(SimpleValidator)
	return isValidator || isSimpleValidator
}

// Validate calls Valid and/or Validate if the given value implements SimpleValidator or Validator.
// Panics if the given value implements none of the two.
func Validate(v interface{}) error {
	validated := false

	if v, ok := v.(SimpleValidator); ok {
		validated = true
		if !v.Valid() {
			return NewErrFailedValidation("invalid", nil, errorz.SkipPackage())
		}
	}

	if v, ok := v.(Validator); ok {
		validated = true
		if err := v.Validate(); err != nil {
			return WrapErrFailedValidation(err, errorz.SkipPackage())
		}
	}

	errorz.Assertf(validated, "value must implement SimpleValidator or Validator", errorz.SkipPackage())
	return nil
}

// MustValidate is like Validate but panics on error.
func MustValidate(v interface{}) {
	errorz.MaybeMustWrap(Validate(v), errorz.SkipPackage())
}
