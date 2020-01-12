package nrpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServer_HandleError(t *testing.T) {
	s := NewServer(ServerOptions{})
	s.HandleFunc("flake", "create", func(ctx context.Context, req *Request, res *Response) error {
		return &Error{
			Status:  StatusOK,
			Message: "test",
		}
	})
	err := s.Start(":18899")
	require.NoError(t, err)
	defer s.Shutdown()

	var m *Response
	c := NewClient(ClientOptions{})
	c.Register("flake", "127.0.0.1:18899")
	req := NewRequest("flake", "create")
	m, err = c.Invoke(context.Background(), req, nil)
	require.NotEmpty(t, m.Metadata.Get(MetadataKeyTrackId))
	require.NotEmpty(t, m.Metadata.Get(MetadataKeyHostname))
	require.Equal(t, StatusOK, m.Status)
	require.Equal(t, "test", m.Message)
	require.NoError(t, err)
}
