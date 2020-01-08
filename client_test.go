package nrpc

import (
	"github.com/stretchr/testify/require"
	"net"
	"net/url"
	"testing"
)

type testReq struct {
	Key1 string `json:"key_1"`
}

type testResp struct {
	Key2 string `json:"key_2"`
}

func handleTestConn(c net.Conn, t *testing.T) {
	defer c.Close()
	var err error
	var req *Message
	req, err = NewIncomingMessage(c)
	require.NoError(t, err)
	body := testReq{}
	err = req.Recv(&body)
	require.NoError(t, err)
	require.Equal(t, "1", body.Key1)
	om := NewOutgoingMessage(c)
	om.Title = StatusOK
	om.Subtitle = "ok"
	om.Metadata.Set("hello", req.Metadata.Get("hello"))
	pbody := testResp{Key2: req.Subtitle}
	err = om.Send(&pbody)
	require.NoError(t, err)
}

func TestClient_Invoke(t *testing.T) {
	l, err := net.Listen("tcp", ":10611")
	require.NoError(t, err)
	defer l.Close()

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				break
			}
			go handleTestConn(c, t)
		}
	}()

	c := NewClient()
	c.Register("test", "127.0.0.1:10611")
	mt := url.Values{}
	mt.Set("hello", "world")
	out := testResp{}
	var mp *Message
	mp, err = c.Invoke("test", "tmethod", mt, &testReq{Key1: "1"}, &out)
	require.NoError(t, err)
	require.Equal(t, "tmethod", out.Key2)
	require.Equal(t, "world", mp.Metadata.Get("hello"))
}
