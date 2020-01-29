package main

import (
	"context"
	"errors"
	"go.guoyk.net/nrpc/v2"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

type AddService struct {
	Count int64
}

type AddIn struct {
	A int `json:"a" query:"a" default:"1"`
	B int `json:"b" query:"b" default:"1"`
}

type AddOut struct {
	V int `json:"v" query:"v"`
}

func (a *AddService) HealthCheck(ctx context.Context) error {
	if atomic.AddInt64(&a.Count, 1)%2 == 0 {
		return nil
	}
	return errors.New("test error")
}

func (a *AddService) Add(ctx context.Context, in *AddIn) (out AddOut, err error) {
	out.V = in.A + in.B
	return
}

func main() {
	s := nrpc.NewServer(nrpc.ServerOptions{Addr: "127.0.0.1:3000"})
	s.Register(&AddService{})
	s.Start(nil)
	defer s.Shutdown(context.Background())

	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)
	<-chSig
}
