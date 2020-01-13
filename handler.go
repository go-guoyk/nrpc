package nrpc

import "context"

type Handler interface {
	ServeNRPC(ctx context.Context, nreq *Request, nres *Response) (err error)
}

type HandlerFunc func(ctx context.Context, nreq *Request, nres *Response) (err error)

func (h HandlerFunc) ServeNRPC(ctx context.Context, nreq *Request, nres *Response) error {
	return h(ctx, nreq, nres)
}

var (
	NotFound HandlerFunc = func(ctx context.Context, nreq *Request, nres *Response) (err error) {
		nres.Status = StatusErrNotImplemented
		nres.Message = "service or method not implemented"
		return
	}
)
