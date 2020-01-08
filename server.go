package nrpc

import (
	"context"
	"net"
	"sync"
)

type HandlerFunc func(ctx context.Context, req *Message, res *Message) (err error)

var DefaultServiceNotFound HandlerFunc = func(ctx context.Context, req *Message, res *Message) (err error) {
	res.Title = StatusErrNotImplemented
	res.Subtitle = "service not implemented"
	return res.Send(nil)
}

var DefaultMethodNotFound HandlerFunc = func(ctx context.Context, req *Message, res *Message) (err error) {
	res.Title = StatusErrNotImplemented
	res.Subtitle = "method not implemented"
	return res.Send(nil)
}

type Server struct {
	ServiceNotFound HandlerFunc
	MethodNotFound  HandlerFunc

	services  map[string]map[string]HandlerFunc
	servicesL sync.Locker
}

func NewServer() *Server {
	return &Server{
		services:        map[string]map[string]HandlerFunc{},
		servicesL:       &sync.Mutex{},
		ServiceNotFound: DefaultServiceNotFound,
		MethodNotFound:  DefaultMethodNotFound,
	}
}

func (s *Server) Register(service string, method string, hf HandlerFunc) {
	s.servicesL.Lock()
	defer s.servicesL.Unlock()
	svc := s.services[service]
	if svc == nil {
		svc = map[string]HandlerFunc{}
		s.services[service] = svc
	}
	svc[method] = hf
}

func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()
	var err error
	var mr *Message
	if mr, err = NewIncomingMessage(conn); err != nil {
		return
	}
	ms := NewOutgoingMessage(conn)
	svc := s.services[mr.Title]
	if svc == nil {
		_ = s.ServiceNotFound(context.Background(), mr, ms)
		return
	}
	mtd := svc[mr.Subtitle]
	if mtd == nil {
		_ = s.MethodNotFound(context.Background(), mr, ms)
		return
	}
	_ = mtd(context.Background(), mr, ms)
}

func (s *Server) Serve(l net.Listener) (err error) {
	for {
		var conn net.Conn
		if conn, err = l.Accept(); err != nil {
			return
		}
		go s.Handle(conn)
	}
}
