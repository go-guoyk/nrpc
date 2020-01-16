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
	s.Handle("flake", "create", &Handler{
		Serve: func(ctx context.Context, req *Request, res *Response) (err error) {
			res.Status = StatusOK
			res.Message = "OK"
			res.Payload = json.RawMessage{'{', '}'}
			return
		},
	})
	err := s.Start(":18898")
	require.NoError(t, err)
	defer s.Shutdown()

	out := map[string]interface{}{}

	var nres *Response
	c := NewClient(ClientOptions{})
	c.Register("flake", "127.0.0.1:18898")
	nreq := NewRequest("flake", "create")
	nres, err = c.Invoke(context.Background(), nreq, &out)
	require.NotEmpty(t, nres.Metadata.Get(MetadataKeyTrackId))
	require.NotEmpty(t, nres.Metadata.Get(MetadataKeyHostname))
	require.Equal(t, StatusOK, nres.Status)
	require.Equal(t, "OK", nres.Message)
	require.NoError(t, err)
}

func TestServer_HandlePanic(t *testing.T) {
	s := NewServer(ServerOptions{})
	s.Handle("flake", "create", &Handler{
		Serve: func(ctx context.Context, req *Request, res *Response) (err error) {
			panic("test")
			return
		},
	})
	err := s.Start(":18878")
	require.NoError(t, err)
	defer s.Shutdown()

	out := map[string]interface{}{}

	var nres *Response
	c := NewClient(ClientOptions{})
	c.Register("flake", "127.0.0.1:18878")
	nreq := NewRequest("flake", "create")
	nres, err = c.Invoke(context.Background(), nreq, &out)
	require.NotEmpty(t, nres.Metadata.Get(MetadataKeyTrackId))
	require.NotEmpty(t, nres.Metadata.Get(MetadataKeyHostname))
	require.Equal(t, StatusErrInternal, nres.Status)
	require.Equal(t, "test", nres.Message)
}

func TestServer_Shutdown(t *testing.T) {
	s := NewServer(ServerOptions{})
	s.Handle("test", "test", &Handler{
		Serve: func(ctx context.Context, req *Request, res *Response) (err error) {
			time.Sleep(time.Second)
			res.Message = "OK"
			res.Payload = map[string]string{"Hello": "World"}
			return
		},
	})
	err := s.Start(":17777")
	require.NoError(t, err)
	nreq := NewRequest("test", "test")
	go func() {
		time.Sleep(time.Millisecond * 100)
		s.Shutdown()
	}()
	var nres *Response
	var m map[string]string
	nres, err = DefaultTransport.RoundTrip(context.Background(), "127.0.0.1:17777", nreq, &m)
	require.NoError(t, err)
	require.Equal(t, "OK", nres.Message)
	require.Equal(t, "World", m["Hello"])
}
