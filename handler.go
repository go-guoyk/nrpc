package nrpc

import "context"

type Handler interface {
	ServeNRPC(ctx context.Context, req *Request, res *Response) (err error)
}

type HandlerFunc func(ctx context.Context, req *Request, res *Response) (err error)

func (h HandlerFunc) ServeNRPC(ctx context.Context, req *Request, res *Response) error {
	return h(ctx, req, res)
}

var (
	NotFound HandlerFunc = func(ctx context.Context, req *Request, res *Response) (err error) {
		res.Status = StatusErrNotImplemented
		res.Message = "service or method not implemented"
		return
	}
)
