package nrpc

import (
	"net"
	"sync"
)

type ServerOptions struct {
	NotFound *Handler
}

type Server struct {
	notFound *Handler

	services  map[string]map[string]*Handler
	listener  net.Listener
	waitConns *sync.WaitGroup
	numConns  int64
}

func NewServer(opts ServerOptions) *Server {
	if opts.NotFound == nil {
		opts.NotFound = NotFound
	}
	return &Server{
		notFound:  opts.NotFound,
		services:  map[string]map[string]*Handler{},
		waitConns: &sync.WaitGroup{},
	}
}

func (s *Server) Handle(service string, method string, h *Handler) {
	svc := s.services[service]
	if svc == nil {
		svc = map[string]*Handler{}
		s.services[service] = svc
	}
	svc[method] = h
}

func (s *Server) resolve(service, method string) *Handler {
	svc := s.services[service]
	if svc == nil {
		return s.notFound
	} else {
		h := svc[method]
		if h == nil {
			return s.notFound
		} else {
			return h
		}
	}
}

func (s *Server) ServeConn(conn net.Conn) {
	// TODO: implements
}

func (s *Server) Serve(l net.Listener) (err error) {
	for {
		var conn net.Conn
		if conn, err = l.Accept(); err != nil {
			return
		}
		s.waitConns.Add(1)
		go s.ServeConn(conn)
	}
}

func (s *Server) Start(addr string) (err error) {
	if s.listener != nil {
		return
	}
	var l net.Listener
	if l, err = net.Listen("tcp", addr); err != nil {
		return
	}
	s.listener = l
	go s.Serve(l)
	return
}

func (s *Server) Shutdown() {
	if s.listener == nil {
		return
	}
	_ = s.listener.Close()
	s.listener = nil
	s.waitConns.Wait()
}
