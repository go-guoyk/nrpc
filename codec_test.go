package nrpc

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
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

func TestEncodeMessage(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := encodeMessage(buf, "hello\nworld")
	require.NoError(t, err)
	require.Equal(t, "hello%0Aworld\n", buf.String())
	require.Equal(t, 14, n)
}

func TestDecodeLine(t *testing.T) {
	var err error
	var line []byte
	buf := []byte(" line 1\n  line 2  \n line 3\n")
	line, buf, err = decodeLine(buf)
	require.NoError(t, err)
	require.Equal(t, "line 1", string(line))
	line, buf, err = decodeLine(buf)
	require.NoError(t, err)
	require.Equal(t, "line 2", string(line))
	line, buf, err = decodeLine(buf)
	require.NoError(t, err)
	require.Equal(t, "line 3", string(line))
	line, buf, err = decodeLine(buf)
	require.Equal(t, io.ErrUnexpectedEOF, err)
}

func TestDecodeHeadline(t *testing.T) {
	var err error
	var val1, val2, val3 string
	buf := []byte(" line, 1\n  line, 2,  hello%0Aworld \n line\n")
	buf, err = decodeHeadline(buf, &val1, &val2, &val3)
	require.NoError(t, err)
	require.Equal(t, "line", val1)
	require.Equal(t, "1", val2)
	require.Equal(t, "", val3)
	buf, err = decodeHeadline(buf, &val1, &val2, &val3)
	require.NoError(t, err)
	require.Equal(t, "line", val1)
	require.Equal(t, "2", val2)
	require.Equal(t, "hello\nworld", val3)
	buf, err = decodeHeadline(buf, &val1, &val2, &val3)
	require.NoError(t, err)
	require.Equal(t, "line", val1)
	require.Equal(t, "", val2)
	require.Equal(t, "", val3)
	buf, err = decodeHeadline(buf, &val1, &val2, &val3)
	require.Equal(t, io.ErrUnexpectedEOF, err)
}

func TestDecodeMetadata(t *testing.T) {
	var err error
	var m Metadata
	buf := []byte("\nkey1 = val1, keY2 = hello%0Aworld\n")
	buf, err = decodeMetadata(buf, &m)
	require.NoError(t, err)
	require.NotNil(t, m)
	buf, err = decodeMetadata(buf, &m)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, "val1", m.Get("key1"))
	require.Equal(t, "hello\nworld", m.Get("key2"))
}

func TestDecodePayload(t *testing.T) {
	var err error
	var m map[string]interface{}
	buf := []byte(`{"hello":"world"}` + "\n")
	buf, err = decodePayload(buf, &m)
	require.NoError(t, err)
	require.Equal(t, 0, len(buf))
	require.Equal(t, "world", m["hello"])
}

func TestDecodeMessage(t *testing.T) {
	var err error
	var m string
	buf := []byte("    hello%0Aworld \n")
	buf, err = decodeMessage(buf, &m)
	require.NoError(t, err)
	require.Equal(t, 0, len(buf))
	require.Equal(t, "hello\nworld", m)
}
