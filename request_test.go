package nrpc

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

const sample1 = `hello,world
key1=val1&key2=val2
{"hello":"world"}
`

const sample2 = ` hellO  , wOrld  
  key1=val1&key2=val2  
{"hello":"world"}
`

func TestRequest_WriteTo(t *testing.T) {
	nreq := NewRequest("hello", "world")
	nreq.Metadata.Set("key1", "val1")
	nreq.Metadata.Set("key2", "val2")
	nreq.Payload = map[string]string{
		"hello": "world",
	}
	buf := &bytes.Buffer{}
	n, err := nreq.WriteTo(buf)
	require.NoError(t, err)
	require.Equal(t, int64(50), n)
	require.Equal(t, 50, buf.Len())
	require.Equal(t, sample1, buf.String())
}

func TestReadRequest(t *testing.T) {
	nreq, err := ReadRequest(bytes.NewReader([]byte(sample1)))
	require.NoError(t, err)
	require.Equal(t, "hello", nreq.Service)
	require.Equal(t, "world", nreq.Method)
	require.Equal(t, "val1", nreq.Metadata.Get("key1"))
	require.Equal(t, "val2", nreq.Metadata.Get("key2"))
	var p map[string]string
	err = nreq.Unmarshal(&p)
	require.NoError(t, err)
	require.Equal(t, "world", p["hello"])

	nreq, err = ReadRequest(bytes.NewReader([]byte(sample2)))
	require.NoError(t, err)
	require.Equal(t, "hello", nreq.Service)
	require.Equal(t, "world", nreq.Method)
	require.Equal(t, "val1", nreq.Metadata.Get("key1"))
	require.Equal(t, "val2", nreq.Metadata.Get("key2"))
	p = nil
	err = nreq.Unmarshal(&p)
	require.NoError(t, err)
	require.Equal(t, "world", p["hello"])
}
