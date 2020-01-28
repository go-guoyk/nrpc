package nrpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Call(t *testing.T) {
	s := NewServer(ServerOptions{Addr: "127.0.0.1:10087"})
	s.Register(&TestService{})
	go s.Start(nil)
	defer s.Shutdown(context.Background())

	time.Sleep(time.Second)

	c := NewClient(ClientOptions{MaxRetries: 1})
	c.Register("TestService", "127.0.0.1:10087")

	err := c.Query("TestService", "Method1").Do(context.Background())
	require.Error(t, err)
	require.Equal(t, "test error", err.Error())
	require.False(t, IsUserError(err))

	in := &TestIn{Hello: "world3"}
	err = c.Query("TestService", "Method2").In(in).Do(context.Background())
	require.Error(t, err)
	require.Equal(t, "test error: world3", err.Error())
	require.True(t, IsUserError(err))
}
