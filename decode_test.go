package nrpc

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
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
	var m Metadata

	r = bufio.NewReader(bytes.NewReader([]byte("hello")))
	err = DecodeMetadata(r, &m)
	require.Equal(t, io.ErrUnexpectedEOF, err)

	r = bufio.NewReader(bytes.NewReader([]byte("key1=val1\n")))
	err = DecodeMetadata(r, &m)
	require.NoError(t, err)
	require.Equal(t, "val1", m.Get("key1"))

	m = nil
	r = bufio.NewReader(bytes.NewReader([]byte("\n")))
	err = DecodeMetadata(r, &m)
	require.NoError(t, err)
	require.NotNil(t, m)

	m = nil
	r = bufio.NewReader(bytes.NewReader([]byte("key==%%=\n")))
	err = DecodeMetadata(r, &m)
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

	err = DecodePayload(r, nil)
	require.NoError(t, err)
}
