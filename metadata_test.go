package nrpc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetadata_Encode(t *testing.T) {
	m := Metadata{}
	m.Set("hello", "world")
	require.Equal(t, "hello=world", m.Encode())
	m = Metadata{}
	m.Set("hello", "wor\nld")
	require.Equal(t, "hello=wor%0Ald", m.Encode())
}

func TestParseMetadata(t *testing.T) {
	v := "hello=%0Ald"
	m, err := ParseMetadata(v)
	require.NoError(t, err)
	require.Equal(t, "\nld", m.Get("hello"))
}
