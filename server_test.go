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
	s.Register("flake", "create", func(ctx context.Context, req *Request, res *Response) (err error) {
		res.Status = StatusOK
		res.Message = "OK"
		return
	})
	go s.Serve(l)

	out := map[string]interface{}{}

	var m *Response
	c := NewClient()
	c.Register("flake", "127.0.0.1:18898")
	req := NewRequest("flake", "create")
	m, err = c.Invoke(context.Background(), req, &out)
	require.Equal(t, StatusOK, m.Status)
	require.Equal(t, "OK", m.Message)
	require.NoError(t, err)
}
