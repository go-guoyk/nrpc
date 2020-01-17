package nrpc

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	requestSample1 = `flake, create
hostname = MicroPC, track_Id = 111
{"hello":"world"}
`
	requestSample2 = `flake,create
hello=++hello%0Aworld
{"hello":"world"}
`
)

type testStruct struct {
	Hello string `json:"hello"`
}

func TestRequest_Decode(t *testing.T) {
	r := &Request{}
	s := testStruct{}
	r.Payload = &s
	err := r.Decode([]byte(requestSample1))
	require.NoError(t, err)
	require.Equal(t, "111", r.Metadata.Get("track_id"))
	require.Equal(t, "MicroPC", r.Metadata.Get("hostname"))
	require.Equal(t, "create", r.Method)
	require.Equal(t, "flake", r.Service)
	require.Equal(t, "world", s.Hello)
}

func TestRequest_Encode(t *testing.T) {
	var err error
	buf := &bytes.Buffer{}
	r := &Request{Metadata: Metadata{}}
	err = r.Encode(buf)
	require.Error(t, err)
	r.Service = "flake"
	r.Method = "create"
	r.Metadata.Set("hello", "  hello\nworld")
	r.Payload = &testStruct{Hello: "world"}
	err = r.Encode(buf)
	require.NoError(t, err)
	require.Equal(t, requestSample2, buf.String())
}
