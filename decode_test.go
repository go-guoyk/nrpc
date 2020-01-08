package nrpc

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/url"
	"testing"
)

func TestDecodeHeadline(t *testing.T) {
	var r *bufio.Reader
	var err error
	var val1, val2 string

	r = bufio.NewReader(bytes.NewReader([]byte("hello")))
	err = DecodeHeadline(r, &val1, &val2)
	require.Equal(t, io.ErrUnexpectedEOF, err)

	r = bufio.NewReader(bytes.NewReader([]byte("hello\n")))
	val2 = "xxx"
	err = DecodeHeadline(r, &val1, &val2)
	require.NoError(t, err)
	require.Equal(t, "hello", val1)
	require.Equal(t, "", val2)

	r = bufio.NewReader(bytes.NewReader([]byte("hello,world\n")))
	err = DecodeHeadline(r, &val1, &val2)
	require.NoError(t, err)
	require.Equal(t, "hello", val1)
	require.Equal(t, "world", val2)

	r = bufio.NewReader(bytes.NewReader([]byte("hello,world,world2\n")))
	err = DecodeHeadline(r, &val1, &val2)
	require.NoError(t, err)
	require.Equal(t, "hello", val1)
	require.Equal(t, "world,world2", val2)
}

func TestDecodeMetadata(t *testing.T) {
	var r *bufio.Reader
	var err error
	var v url.Values

	r = bufio.NewReader(bytes.NewReader([]byte("hello")))
	err = DecodeMetadata(r, &v)
	require.Equal(t, io.ErrUnexpectedEOF, err)

	r = bufio.NewReader(bytes.NewReader([]byte("key1=val1\n")))
	err = DecodeMetadata(r, &v)
	require.NoError(t, err)
	require.Equal(t, "val1", v.Get("key1"))

	v = nil
	r = bufio.NewReader(bytes.NewReader([]byte("\n")))
	err = DecodeMetadata(r, &v)
	require.NoError(t, err)
	require.NotNil(t, v)

	v = nil
	r = bufio.NewReader(bytes.NewReader([]byte("==%%=\n")))
	err = DecodeMetadata(r, &v)
	require.Error(t, err)
}

func TestDecodePayload(t *testing.T) {
	var r *bufio.Reader
	var err error
	var v map[string]string

	r = bufio.NewReader(bytes.NewReader([]byte(`{"hello":"world"}`)))
	err = DecodePayload(r, &v)
	require.NoError(t, err)
	require.Equal(t, "world", v["hello"])
}
