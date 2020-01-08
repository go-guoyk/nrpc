package nrpc

import (
	"context"
	"encoding/json"
	"go.guoyk.net/trackid"
	"net"
	"sync"
)

type HandlerFunc func(ctx context.Context, req *Request, res *Response) (err error)

var DefaultServiceNotFound HandlerFunc = func(ctx context.Context, req *Request, res *Response) (err error) {
	res.Status = StatusErrNotImplemented
	res.Message = "service not implemented"
	return
}

var DefaultMethodNotFound HandlerFunc = func(ctx context.Context, req *Request, res *Response) (err error) {
	res.Status = StatusErrNotImplemented
	res.Message = "method not implemented"
	return
}

var EmptyResponsePayload = json.RawMessage("{}")

type Server struct {
	ServiceNotFound HandlerFunc
	MethodNotFound  HandlerFunc

	services  map[string]map[string]HandlerFunc
	servicesL *sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		services:        map[string]map[string]HandlerFunc{},
		servicesL:       &sync.RWMutex{},
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

func (s *Server) Method(service, method string) HandlerFunc {
	s.servicesL.RLock()
	defer s.servicesL.RUnlock()

	svc := s.services[service]
	if svc == nil {
		return s.ServiceNotFound
	} else {
		fn := svc[method]
		if fn == nil {
			return s.MethodNotFound
		} else {
			return fn
		}
	}
}

func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()
	var err error
	var req *Request
	if req, err = ReadRequest(conn); err != nil {
		return
	}

	ctx := trackid.Set(context.Background(), req.Metadata.Get("track_id"))

	fn := s.Method(req.Service, req.Method)

	resp := NewResponse()

	if err = fn(ctx, req, resp); err != nil {
		resp.Status = StatusErrInternal
		resp.Message = err.Error()
	}
	if resp.Payload == nil {
		resp.Payload = EmptyResponsePayload
	}

	_, _ = resp.WriteTo(conn)
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
