package nrpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestCall_Do(t *testing.T) {
	s := NewServer(ServerOptions{Addr: "127.0.0.1:10099"})
	s.Register(&TestService{})
	s.Start(nil)
	defer s.Shutdown(context.Background())

	time.Sleep(time.Second)

	c := &Call{
		client:  http.DefaultClient,
		host:    "127.0.0.1:10099",
		service: "TestService",
		method:  "Method3",
		command: true,
	}
	in := &TestIn{Hello: "world2"}
	out := &TestOut{}
	c.in = in
	c.out = out

	err := c.Do(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "world2", out.Hello)

	c = &Call{
		client:  http.DefaultClient,
		host:    "127.0.0.1:10099",
		service: "TestService",
		method:  "Method3",
		command: false,
	}
	in = &TestIn{Hello: "world2"}
	out = &TestOut{}
	c.in = in
	c.out = out

	err = c.Do(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "world2", out.Hello)

	c = &Call{
		client:  http.DefaultClient,
		host:    "127.0.0.1:10099",
		service: "TestService",
		method:  "Method2",
		command: false,
	}
	in = &TestIn{Hello: "world2"}
	c.in = in

	err = c.Do(context.Background())
	assert.Error(t, err)
	assert.True(t, IsSolid(err))
	assert.Equal(t, "test error: world2", err.Error())
}
