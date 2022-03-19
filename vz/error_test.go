package vz_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ibrt/golang-errors/errorz"
	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-validation/vz"
)

func TestErrFailedValidation(t *testing.T) {
	err := vz.NewErrFailedValidation("test error", map[string]interface{}{"field": "missing"}, errorz.Prefix("test prefix"))
	require.Error(t, err)
	require.Equal(t, "failed validation: test prefix: test error", err.Error())
	require.Equal(t, "failed-validation", errorz.GetID(err).String())
	require.Equal(t, http.StatusBadRequest, errorz.GetStatus(err).Int())
	require.Equal(t, errorz.Metadata{"fields": map[string]interface{}{"field": "missing"}}, errorz.GetMetadata(err))
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[0], "vz_test.TestErrFailedValidation"))

	err = vz.WrapErrFailedValidation(errorz.Errorf("test error"), errorz.Prefix("test prefix"))
	require.Error(t, err)
	require.Equal(t, "failed validation: test prefix: test error", err.Error())
	require.Equal(t, "failed-validation", errorz.GetID(err).String())
	require.Equal(t, http.StatusBadRequest, errorz.GetStatus(err).Int())
	require.Equal(t, errorz.Metadata{"fields": map[string]interface{}{}}, errorz.GetMetadata(err))
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[0], "vz_test.TestErrFailedValidation"))
	require.Equal(t, vz.WrapErrFailedValidation(err), err)
}
