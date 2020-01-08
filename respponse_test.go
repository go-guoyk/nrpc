package nrpc

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

const sampleR1 = `hello,world
key1=val1&key2=val2
{"hello":"world"}
`

const sampleR2 = ` hellO  , wOrld  
  key1=val1&key2=val2  
{"hello":"world"}
`

func TestResponse_WriteTo(t *testing.T) {
	req := NewResponse()
	req.Status = "hello"
	req.Message = "world"
	req.Metadata.Set("key1", "val1")
	req.Metadata.Set("key2", "val2")
	req.Payload = map[string]string{
		"hello": "world",
	}
	buf := &bytes.Buffer{}
	n, err := req.WriteTo(buf)
	require.NoError(t, err)
	require.Equal(t, int64(50), n)
	require.Equal(t, 50, buf.Len())
	require.Equal(t, sampleR1, buf.String())
}

func TestReadResponse(t *testing.T) {
	req, err := ReadResponse(bytes.NewReader([]byte(sampleR1)))
	require.NoError(t, err)
	require.Equal(t, "hello", req.Status)
	require.Equal(t, "world", req.Message)
	require.Equal(t, "val1", req.Metadata.Get("key1"))
	require.Equal(t, "val2", req.Metadata.Get("key2"))
	var p map[string]string
	err = req.Unmarshal(&p)
	require.NoError(t, err)
	require.Equal(t, "world", p["hello"])

	req, err = ReadResponse(bytes.NewReader([]byte(sampleR2)))
	require.NoError(t, err)
	require.Equal(t, "hello", req.Status)
	require.Equal(t, "wOrld", req.Message)
	require.Equal(t, "val1", req.Metadata.Get("key1"))
	require.Equal(t, "val2", req.Metadata.Get("key2"))
	p = nil
	err = req.Unmarshal(&p)
	require.NoError(t, err)
	require.Equal(t, "world", p["hello"])
}
