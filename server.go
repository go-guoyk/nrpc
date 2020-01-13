package nrpc

import (
	"context"
	"go.guoyk.net/trackid"
	"net"
	"sync"
	"sync/atomic"
)

type ServerOptions struct {
	DefaultHandler Handler
}

type Server struct {
	defaultHandler Handler

	services  map[string]map[string]Handler
	listener  net.Listener
	waitConns *sync.WaitGroup
	numConns  int64
}

func NewServer(opts ServerOptions) *Server {
	if opts.DefaultHandler == nil {
		opts.DefaultHandler = NotFound
	}
	return &Server{
		defaultHandler: opts.DefaultHandler,
		services:       map[string]map[string]Handler{},
		waitConns:      &sync.WaitGroup{},
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
		return s.defaultHandler
	} else {
		h := svc[method]
		if h == nil {
			return s.defaultHandler
		} else {
			return h
		}
	}
}

func (s *Server) ServeConn(conn net.Conn) {
	defer s.waitConns.Done()
	defer conn.Close()

	// update num conns
	atomic.AddInt64(&s.numConns, 1)
	defer atomic.AddInt64(&s.numConns, -1)

	// read request
	var err error
	var req *Request
	if req, err = ReadRequest(conn); err != nil {
		return
	}

	// prepare context
	ctx := context.Background()
	ctx = trackid.Set(ctx, req.Metadata.Get(MetadataKeyTrackId))

	// prepare response
	res := NewResponse()
	res.Metadata.Set(MetadataKeyHostname, hostname)
	res.Metadata.Set(MetadataKeyTrackId, trackid.Get(ctx))

	// find handler
	h := s.Handler(req.Service, req.Method)

	// execute handler
	if err = h.ServeNRPC(ctx, req, res); err != nil {
		if ne, ok := err.(*Error); ok {
			res.Status = ne.Status
			res.Message = ne.Message
		} else {
			res.Status = StatusErrInternal
			res.Message = err.Error()
		}
	}

	// write response
	_, _ = res.WriteTo(conn)
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
