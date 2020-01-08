package nrpc

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestEncodeHeadline(t *testing.T) {
	var n int
	var err error
	var buf *bytes.Buffer
	buf = &bytes.Buffer{}
	n, err = EncodeHeadline(buf, "hello", "")
	require.NoError(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, "hello\n", buf.String())
	buf = &bytes.Buffer{}
	n, err = EncodeHeadline(buf, "  hello  ", "")
	require.NoError(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, "hello\n", buf.String())
	buf = &bytes.Buffer{}
	n, err = EncodeHeadline(buf, "  hello   ", "  world   ")
	require.NoError(t, err)
	require.Equal(t, 12, n)
	require.Equal(t, "hello,world\n", buf.String())
}

func TestEncodeMetadata(t *testing.T) {
	var n int
	var v url.Values
	var err error
	var buf *bytes.Buffer
	buf = &bytes.Buffer{}
	n, err = EncodeMetadata(buf, v)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, "\n", buf.String())
	buf = &bytes.Buffer{}
	v = url.Values{}
	n, err = EncodeMetadata(buf, v)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, "\n", buf.String())
	buf = &bytes.Buffer{}
	v = url.Values{}
	v.Set("hello", "world")
	n, err = EncodeMetadata(buf, v)
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
	n, err = EncodePayload(buf, v)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, "\n", buf.String())
	buf = &bytes.Buffer{}
	v = map[string]string{"hello": "world"}
	n, err = EncodePayload(buf, v)
	require.NoError(t, err)
	require.Equal(t, 18, n)
	require.Equal(t, "{\"hello\":\"world\"}\n", buf.String())
}
