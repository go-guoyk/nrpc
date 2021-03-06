package nrpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

type TestIn struct {
	Hello string `json:"hello" query:"hell"`
}

type TestOut struct {
	Hello string `json:"hello"`
}

type TestService struct{}

func (s *TestService) Method1(ctx context.Context) (err error) {
	err = errors.New("test error")
	return
}

func (s *TestService) Method2(ctx context.Context, arg *TestIn) (err error) {
	err = Solid(fmt.Errorf("test error: %s", arg.Hello))
	return
}

func (s *TestService) Method3(ctx context.Context, arg *TestIn) (out TestOut, err error) {
	out.Hello = arg.Hello
	return
}

func (s *TestService) Method4(ctx context.Context) (out TestOut, err error) {
	out.Hello = "world"
	return
}

func TestExtractHandlers(t *testing.T) {
	hs := ExtractHandlers("TestService", &TestService{})
	assert.Equal(t, 4, len(hs))
	assert.Nil(t, hs["Method1"].inType)
	assert.Equal(t, reflect.Struct, hs["Method2"].inType.Kind())
}

func TestHandler_ServeHTTP(t *testing.T) {
	hs := ExtractHandlers("TestService", &TestService{})
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost:3000/TestService/Method1", nil)
	hs["Method1"].ServeHTTP(rw, req)
	assert.Equal(t, "test error", rw.Body.String())
	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rw.Header().Get("Content-Type"))
	assert.NotEmpty(t, rw.Header().Get(headerCorrelationID))

	buf := []byte(`{"hello":"world"}`)
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "http://localhost:3000/TestService/Method2", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Content-Length", strconv.Itoa(len(buf)))
	hs["Method2"].ServeHTTP(rw, req)
	assert.Equal(t, "test error: world", rw.Body.String())
	assert.Equal(t, http.StatusBadRequest, rw.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rw.Header().Get("Content-Type"))
	assert.NotEmpty(t, rw.Header().Get(headerCorrelationID))

	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "http://localhost:3000/TestService/Method3?hell=world", bytes.NewReader(buf))
	hs["Method3"].ServeHTTP(rw, req)
	assert.Equal(t, `{"hello":"world"}`, rw.Body.String())
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "application/json; charset=utf-8", rw.Header().Get("Content-Type"))
	assert.NotEmpty(t, rw.Header().Get(headerCorrelationID))

	buf = []byte(`something not json`)
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "http://localhost:3000/TestService/Method4", bytes.NewReader(buf))
	hs["Method4"].ServeHTTP(rw, req)
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("Content-Length", strconv.Itoa(len(buf)))
	assert.Equal(t, `{"hello":"world"}`, rw.Body.String())
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "application/json; charset=utf-8", rw.Header().Get("Content-Type"))
	assert.NotEmpty(t, rw.Header().Get(headerCorrelationID))
}
