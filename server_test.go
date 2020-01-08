package nrpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestServer_Handle(t *testing.T) {
	var l net.Listener
	l, err := net.Listen("tcp", ":18898")
	require.NoError(t, err)
	defer l.Close()

	s := NewServer()
	s.Register("flake", "create", func(ctx context.Context, req *Message, res *Message) (err error) {
		res.Subject = StatusOK
		res.SecondarySubject = "OK"
		return res.Send(map[string]interface{}{
			"id": 0,
		})
	})
	go s.Serve(l)

	out := map[string]interface{}{}

	var m *Message
	c := NewClient()
	c.Register("flake", "127.0.0.1:18898")
	m, err = c.Invoke("flake", "create", nil, nil, &out)
	require.Equal(t, StatusOK, m.Subject)
	require.Equal(t, "OK", m.SecondarySubject)
	require.NoError(t, err)
}
