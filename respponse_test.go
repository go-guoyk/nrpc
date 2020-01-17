package nrpc

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	responseSample1 = `ok
hostname = MicroPC, track_Id = 111
{"hello":"world"}
`

	responseSample11 = `err_internal
hostname = MicroPC, track_Id = 111
hello%0Aworld
`
	responseSample2 = `ok
hello=++hello%0Aworld
{"hello":"world"}
`

	responseSample21 = `err_internal
hello=++hello%0Aworld
hello%0Aworld
`
)

func TestResponse_Decode(t *testing.T) {
	r := &Response{}
	s := testStruct{}
	r.Payload = &s
	err := r.Decode([]byte(responseSample1))
	require.NoError(t, err)
	require.Equal(t, "111", r.Metadata.Get("track_id"))
	require.Equal(t, "MicroPC", r.Metadata.Get("hostname"))
	require.Equal(t, "ok", r.Status)
	require.Equal(t, "", r.Message)
	require.Equal(t, "world", s.Hello)

	r = &Response{}
	s = testStruct{}
	r.Payload = &s
	err = r.Decode([]byte(responseSample11))
	require.NoError(t, err)
	require.Equal(t, "111", r.Metadata.Get("track_id"))
	require.Equal(t, "MicroPC", r.Metadata.Get("hostname"))
	require.Equal(t, "err_internal", r.Status)
	require.Equal(t, "hello\nworld", r.Message)
}

func TestResponse_Encode(t *testing.T) {
	var err error
	buf := &bytes.Buffer{}
	r := &Response{Metadata: Metadata{}}
	err = r.Encode(buf)
	require.Error(t, err)
	r.Status = StatusOK
	r.Metadata.Set("hello", "  hello\nworld")
	r.Payload = &testStruct{Hello: "world"}
	err = r.Encode(buf)
	require.NoError(t, err)
	require.Equal(t, responseSample2, buf.String())

	buf = &bytes.Buffer{}
	r = &Response{Metadata: Metadata{}}
	err = r.Encode(buf)
	require.Error(t, err)
	r.Status = StatusErrInternal
	r.Metadata.Set("hello", "  hello\nworld")
	r.Message = "hello\nworld"
	r.Payload = &testStruct{Hello: "world"}
	err = r.Encode(buf)
	require.NoError(t, err)
	require.Equal(t, responseSample21, buf.String())
}
