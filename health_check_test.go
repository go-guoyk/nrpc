package nrpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestHealthCheck struct {
	err error
}

func (hc *TestHealthCheck) HealthCheck(ctx context.Context) error {
	return hc.err
}

func TestHealthChecks(t *testing.T) {
	hc1 := &TestHealthCheck{}
	hc2 := &TestHealthCheck{}
	hcs := &HealthChecks{}
	hcs.Add(hc1)
	hcs.Add(hc2)
	require.NoError(t, hcs.HealthCheck(context.Background()))

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:3000", nil)
	hcs.ServeHTTP(rw, req)
	require.Equal(t, http.StatusOK, rw.Code)

	hc1.err = errors.New("test error")
	require.Error(t, hcs.HealthCheck(context.Background()))
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:3000", nil)
	hcs.ServeHTTP(rw, req)
	require.Equal(t, http.StatusInternalServerError, rw.Code)
	require.Equal(t, "test error", rw.Body.String())
}
