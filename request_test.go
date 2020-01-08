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

func TestRequest_WriteTo(t *testing.T) {
	req := NewRequest("hello", "world")
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
	require.Equal(t, sample1, buf.String())
}
