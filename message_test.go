package nrpc

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

type testStruct struct {
	ValA string `json:"a"`
}

const sampleM1 = `service1,method1
track_id: aaa	 
instance_id: localhost

{"a":"b"}`

const sampleM2 = `service1,method1

`

const sampleM3 = `err_internal,internal error
track_id: 11111

{"key1":"val1","key2":2}
`

func TestNewIncomingMessage(t *testing.T) {
	r := bytes.NewReader([]byte(sampleM1))

	m, err := NewIncomingMessage(r)
	require.NoError(t, err)
	require.Equal(t, "aaa", m.Metadata.Get("track_id"))
	require.Equal(t, "localhost", m.Metadata.Get("instance_id"))

	body := testStruct{}
	require.NoError(t, m.Recv(&body))
	require.Equal(t, "b", body.ValA)

	r = bytes.NewReader([]byte(sampleM2))

	m, err = NewIncomingMessage(r)
	require.NoError(t, err)
	require.Equal(t, 0, len(m.Metadata))

	body = testStruct{}
	require.NoError(t, m.Recv(nil))
}

func TestNewOutgoingMessage(t *testing.T) {
	buf := &bytes.Buffer{}
	m := NewOutgoingMessage(buf)
	m.Subject = StatusErrInternal
	m.SecondarySubject = "internal error"
	m.Metadata.Add("track_id", "11111")
	err := m.Send(map[string]interface{}{
		"key1": "val1",
		"key2": 2,
	})
	require.NoError(t, err)
	require.Equal(t, sampleM3, buf.String())
}
