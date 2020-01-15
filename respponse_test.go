package nrpc

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

const sampleR1 = `hello,world
key1=val1;key2=val2
{"hello":"world"}
`

const sampleR2 = ` hellO  , wOrld  
  key1=val1;key2=val2  
{"hello":"world"}
`

func TestResponse_WriteTo(t *testing.T) {
	nres := NewResponse()
	nres.Status = "hello"
	nres.Message = "world"
	nres.Metadata.Set("key1", "val1")
	nres.Metadata.Set("key2", "val2")
	nres.Payload = map[string]string{
		"hello": "world",
	}
	buf := &bytes.Buffer{}
	n, err := nres.WriteTo(buf)
	require.NoError(t, err)
	require.Equal(t, int64(50), n)
	require.Equal(t, 50, buf.Len())
	require.Equal(t, sampleR1, buf.String())
}

func TestReadResponse(t *testing.T) {
	nres, err := ReadResponse(bytes.NewReader([]byte(sampleR1)))
	require.NoError(t, err)
	require.Equal(t, "hello", nres.Status)
	require.Equal(t, "world", nres.Message)
	require.Equal(t, "val1", nres.Metadata.Get("key1"))
	require.Equal(t, "val2", nres.Metadata.Get("key2"))
	var p map[string]string
	err = nres.Unmarshal(&p)
	require.NoError(t, err)
	require.Equal(t, "world", p["hello"])

	nres, err = ReadResponse(bytes.NewReader([]byte(sampleR2)))
	require.NoError(t, err)
	require.Equal(t, "hello", nres.Status)
	require.Equal(t, "wOrld", nres.Message)
	require.Equal(t, "val1", nres.Metadata.Get("key1"))
	require.Equal(t, "val2", nres.Metadata.Get("key2"))
	p = nil
	err = nres.Unmarshal(&p)
	require.NoError(t, err)
	require.Equal(t, "world", p["hello"])
}
