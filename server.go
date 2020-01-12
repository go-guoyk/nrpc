package nrpc

import (
	"context"
	"go.guoyk.net/trackid"
	"net"
	"sync"
)

type ServerOptions struct {
	DefaultHandler Handler
}

type Server struct {
	DefaultHandler Handler

	services map[string]map[string]Handler
	listener net.Listener
	conns    *sync.WaitGroup
}

func NewServer(opts ServerOptions) *Server {
	if opts.DefaultHandler == nil {
		opts.DefaultHandler = NotFound
	}
	return &Server{
		DefaultHandler: opts.DefaultHandler,
		services:       map[string]map[string]Handler{},
		conns:          &sync.WaitGroup{},
	}
}

func (s *Server) Handle(service string, method string, h Handler) {
	svc := s.services[service]
	if svc == nil {
		svc = map[string]Handler{}
		s.services[service] = svc
	}
	svc[method] = h
}

func (s *Server) HandleFunc(service string, method string, h HandlerFunc) {
	s.Handle(service, method, h)
}

func (s *Server) Handler(service, method string) Handler {
	svc := s.services[service]
	if svc == nil {
		return s.DefaultHandler
	} else {
		h := svc[method]
		if h == nil {
			return s.DefaultHandler
		} else {
			return h
		}
	}
}

func (s *Server) ServeConn(conn net.Conn) {
	defer s.conns.Done()
	defer conn.Close()
	var err error
	var req *Request
	if req, err = ReadRequest(conn); err != nil {
		return
	}

	ctx := context.Background()
	ctx = trackid.Set(ctx, req.Metadata.Get(MetadataKeyTrackId))

	h := s.Handler(req.Service, req.Method)

	res := NewResponse()

	res.Metadata.Set(MetadataKeyHostname, hostname)
	res.Metadata.Set(MetadataKeyTrackId, trackid.Get(ctx))

	if err = h.ServeNRPC(ctx, req, res); err != nil {
		if ne, ok := err.(*Error); ok {
			res.Status = ne.Status
			res.Message = ne.Message
		} else {
			res.Status = StatusErrInternal
			res.Message = err.Error()
		}
	}

	_, _ = res.WriteTo(conn)
}

func (s *Server) Serve(l net.Listener) (err error) {
	for {
		var conn net.Conn
		if conn, err = l.Accept(); err != nil {
			return
		}
		s.conns.Add(1)
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
	s.conns.Wait()
}
