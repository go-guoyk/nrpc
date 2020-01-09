package nrpc

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestServer_Handle(t *testing.T) {
	s := NewServer(ServerOptions{})
	s.HandleFunc("flake", "create", func(ctx context.Context, req *Request, res *Response) (err error) {
		res.Status = StatusOK
		res.Message = "OK"
		res.Payload = json.RawMessage{'{', '}'}
		return
	})
	err := s.Start(":18898")
	require.NoError(t, err)
	defer s.Shutdown()

	out := map[string]interface{}{}

	var m *Response
	c := NewClient(ClientOptions{})
	c.Register("flake", "127.0.0.1:18898")
	req := NewRequest("flake", "create")
	m, err = c.Invoke(context.Background(), req, &out)
	require.NotEmpty(t, m.Metadata.Get(MetadataKeyTrackId))
	require.NotEmpty(t, m.Metadata.Get(MetadataKeyHostname))
	require.Equal(t, StatusOK, m.Status)
	require.Equal(t, "OK", m.Message)
	require.NoError(t, err)
}

func TestServer_Shutdown(t *testing.T) {
	s := NewServer(ServerOptions{})
	s.HandleFunc("test", "test", func(ctx context.Context, req *Request, res *Response) (err error) {
		time.Sleep(time.Second)
		res.Message = "OK"
		res.Payload = map[string]string{"Hello": "World"}
		return
	})
	err := s.Start(":17777")
	require.NoError(t, err)
	req := NewRequest("test", "test")
	go func() {
		time.Sleep(time.Millisecond * 100)
		s.Shutdown()
	}()
	var res *Response
	var m map[string]string
	res, err = Invoke(context.Background(), "127.0.0.1:17777", req, &m)
	require.NoError(t, err)
	require.Equal(t, "OK", res.Message)
	require.Equal(t, "World", m["Hello"])
}
