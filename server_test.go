package nrpc

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

type IDNewIn struct {
	Count uint64 `json:"count" query:"count" default:"1"`
}

type IDNewOut struct {
	ID uint64 `json:"id"`
}

type IDService struct {
	id uint64
}

func (ids *IDService) New(ctx context.Context, in *IDNewIn) (out IDNewOut, err error) {
	out.ID = atomic.AddUint64(&ids.id, in.Count)
	return
}

func TestServer(t *testing.T) {
	s := NewServer(ServerOptions{Addr: ":10888"})
	s.Register(&IDService{})
	s.Start(nil)
	defer s.Shutdown(context.Background())

	time.Sleep(time.Second)

	resp, err := http.Get("http://127.0.0.1:10888/IDService/New?count=2")
	require.NoError(t, err)
	defer resp.Body.Close()
	buf := []byte{}
	buf, err = ioutil.ReadAll(resp.Body)
	out := IDNewOut{}
	err = json.Unmarshal(buf, &out)
	require.NoError(t, err)
	require.Equal(t, uint64(2), out.ID)

	resp, err = http.Get("http://127.0.0.1:10888/IDService/New?count=3")
	require.NoError(t, err)
	defer resp.Body.Close()
	buf = []byte{}
	buf, err = ioutil.ReadAll(resp.Body)
	out = IDNewOut{}
	err = json.Unmarshal(buf, &out)
	require.NoError(t, err)
	require.Equal(t, uint64(5), out.ID)
}
