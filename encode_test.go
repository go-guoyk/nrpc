package nrpc

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncodeHeadline(t *testing.T) {
	var n int
	var err error
	var buf *bytes.Buffer
	buf = &bytes.Buffer{}
	n, err = encodeHeadline(buf, "hello", "")
	require.NoError(t, err)
	require.Equal(t, 7, n)
	require.Equal(t, "hello,\n", buf.String())
}

func TestEncodeMetadata(t *testing.T) {
	var n int
	var m Metadata
	var err error
	var buf *bytes.Buffer
	buf = &bytes.Buffer{}
	n, err = encodeMetadata(buf, m)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, "\n", buf.String())
	buf = &bytes.Buffer{}
	m = Metadata{}
	n, err = encodeMetadata(buf, m)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, "\n", buf.String())
	buf = &bytes.Buffer{}
	m = Metadata{}
	m.Set("hello", "world")
	n, err = encodeMetadata(buf, m)
	require.NoError(t, err)
	require.Equal(t, 12, n)
	require.Equal(t, "hello=world\n", buf.String())
}

func TestEncodePayload(t *testing.T) {
	var n int
	var v interface{}
	var err error
	var buf *bytes.Buffer
	buf = &bytes.Buffer{}
	n, err = encodePayload(buf, v)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, "\n", buf.String())
	buf = &bytes.Buffer{}
	v = map[string]string{"hello": "world"}
	n, err = encodePayload(buf, v)
	require.NoError(t, err)
	require.Equal(t, 18, n)
	require.Equal(t, "{\"hello\":\"world\"}\n", buf.String())
}
