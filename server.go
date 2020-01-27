package nrpc

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

type ServerOptions struct {
	Addr string
}

type Server struct {
	s   *http.Server
	mux *http.ServeMux
}

// Register register a rpc object with default name
func (s *Server) Register(r interface{}) {
	t := reflect.TypeOf(r)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	s.RegisterName(t.Name(), r)
}

// RegisterName register a rpc object with given name
func (s *Server) RegisterName(name string, r interface{}) {
	hs := ExtractHandlers(name, r)
	for m, h := range hs {
		s.mux.Handle(fmt.Sprintf("/%s/%s", name, m), h)
	}
}

func NewServer(opts ServerOptions) *Server {
	if opts.Addr == "" {
		opts.Addr = ":3000"
	}
	mux := http.NewServeMux()
	return &Server{
		s: &http.Server{
			Addr:    opts.Addr,
			Handler: mux,
		},
		mux: mux,
	}
}

func (s *Server) Start(ech chan error) {
	go func() {
		ech <- s.s.ListenAndServe()
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
