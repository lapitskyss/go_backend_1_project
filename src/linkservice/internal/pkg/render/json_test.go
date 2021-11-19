package render

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	response := struct {
		Status bool `json:"status"`
	}{true}

	JSON(recorder, response, http.StatusOK)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, "{\"status\":true}", recorder.Body.String())
}

func TestNotFoundError(t *testing.T) {
	recorder := httptest.NewRecorder()

	NotFoundError(recorder)
	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.JSONEq(t, "{\"status\":false, \"error\":\"Not found.\"}", recorder.Body.String())
}

func TestBadRequestError(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := errors.New("unexpected error")

	BadRequestError(recorder, err)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.JSONEq(t, "{\"status\":false, \"error\":\"unexpected error\"}", recorder.Body.String())
}

func TestInternalServerError(t *testing.T) {
	recorder := httptest.NewRecorder()

	InternalServerError(recorder)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}
