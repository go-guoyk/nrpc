package nrpc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetadata_Encode(t *testing.T) {
	m := Metadata{}
	m.Set("heLlo", "world")
	require.Equal(t, "hello=world", m.Encode())
	m = Metadata{}
	m.Set("hellO", "wor\nld")
	m.Set("hellO2", "world2")
	require.Equal(t, "hello=wor%0Ald;hello2=world2", m.Encode())
}

func TestParseMetadata(t *testing.T) {
	v := "hello=%0Ald"
	m, err := ParseMetadata(v)
	require.NoError(t, err)
	require.Equal(t, "\nld", m.Get("hello"))

	v = "hello = %0Ald; hEllo2 = worldd"
	m, err = ParseMetadata(v)
	require.NoError(t, err)
	require.Equal(t, "\nld", m.Get("hello"))
	require.Equal(t, "worldd", m.Get("hello2"))

	v = "hello = %0Ald; hEllo2 = worldd;;;   ; ;"
	m, err = ParseMetadata(v)
	require.NoError(t, err)
	require.Equal(t, "\nld", m.Get("hello"))
	require.Equal(t, "worldd", m.Get("hello2"))

	v = "hello = %0Ald%%; hEllo2 = worldd;;;   ; ;"
	m, err = ParseMetadata(v)
	require.Error(t, err)
}
