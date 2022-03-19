package vz_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-validation/vz"
)

func TestValidateStruct(t *testing.T) {
	vz.MustRegisterValidator("custom-validator", vz.RegexpValidatorFactory(regexp.MustCompile("^valid$")))

	type validatableStruct struct {
		First  string `json:"first" validate:"required"`
		Second string `validate:"custom-validator"`
	}

	err := vz.ValidateStruct(&validatableStruct{})
	require.Error(t, err)
	require.Equal(t, "failed validation: Key: 'validatableStruct.first' Error:Field validation for 'first' failed on the 'required' tag\nKey: 'validatableStruct.Second' Error:Field validation for 'Second' failed on the 'custom-validator' tag", err.Error())
	require.Equal(t, "failed-validation", errorz.GetID(err).String())
	require.Equal(t, http.StatusBadRequest, errorz.GetStatus(err).Int())
	require.Equal(t, errorz.Metadata{"fields": map[string]interface{}{"first": "required", "Second": "custom-validator"}}, errorz.GetMetadata(err))
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[0], "vz_test.TestValidateStruct"))

	require.PanicsWithError(t, "failed validation: Key: 'validatableStruct.first' Error:Field validation for 'first' failed on the 'required' tag\nKey: 'validatableStruct.Second' Error:Field validation for 'Second' failed on the 'custom-validator' tag", func() {
		vz.MustValidateStruct(&validatableStruct{})
	})

	fixturez.RequireNoError(t, vz.ValidateStruct(&validatableStruct{
		First:  "required",
		Second: "valid",
	}))

	require.NotPanics(t, func() {
		vz.MustValidateStruct(&validatableStruct{
			First:  "required",
			Second: "valid",
		})
	})
}

type MockValidator struct {
	valid bool
	err   error

	simpleValidator bool
	validator       bool
}

func (v *MockValidator) Valid() bool {
	v.simpleValidator = true
	return v.valid
}

func (v *MockValidator) Validate() error {
	v.validator = true
	return v.err
}

func TestValidate_Ok(t *testing.T) {
	v := &MockValidator{valid: true, err: nil}
	fixturez.RequireNoError(t, vz.Validate(v))
	require.True(t, v.simpleValidator)
	require.True(t, v.validator)
}

func TestMustValidate_Ok(t *testing.T) {
	v := &MockValidator{valid: true, err: nil}
	require.NotPanics(t, func() {
		vz.MustValidate(v)
	})
	require.True(t, v.simpleValidator)
	require.True(t, v.validator)
}

func TestValidate_SimpleValidator_Error(t *testing.T) {
	v := &MockValidator{valid: false, err: nil}
	require.Error(t, vz.Validate(v))
	require.True(t, v.simpleValidator)
	require.False(t, v.validator)
}

func TestValidate_Validator_Error(t *testing.T) {
	v := &MockValidator{valid: true, err: fmt.Errorf("test error")}
	require.Error(t, vz.Validate(v))
	require.True(t, v.simpleValidator)
	require.True(t, v.validator)
}

func TestMustValidate_Error(t *testing.T) {
	v := &MockValidator{valid: false, err: nil}
	require.Panics(t, func() {
		vz.MustValidate(v)
	})
	require.True(t, v.simpleValidator)
	require.False(t, v.validator)
}

func TestValidate_Invalid(t *testing.T) {
	require.PanicsWithError(t, "value must implement SimpleValidator or Validator", func() {
		_ = vz.Validate(struct{}{})
	})
}

type TestValidator struct {
	// intentionally empty
}

func (*TestValidator) Validate() error {
	return nil
}

type TestSimpleValidator struct {
	// intentionally empty
}

func (*TestSimpleValidator) Valid() bool {
	return true
}

func TestIsValidatable(t *testing.T) {
	require.True(t, vz.IsValidatable(&TestValidator{}))
	require.True(t, vz.IsValidatable(&TestSimpleValidator{}))
	require.False(t, vz.IsValidatable(struct{}{}))
	require.False(t, vz.IsValidatable(""))
}
