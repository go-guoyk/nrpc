package nrpc

import (
	"context"
	"fmt"
)

type Handler interface {
	ServeNRPC(ctx context.Context, nreq *Request, nres *Response) (err error)
}

type HandlerFunc func(ctx context.Context, nreq *Request, nres *Response) (err error)

func (h HandlerFunc) ServeNRPC(ctx context.Context, nreq *Request, nres *Response) error {
	return h(ctx, nreq, nres)
}

func InvokeHandler(ctx context.Context, h Handler, nreq *Request, nres *Response) (err error) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			nres.Status = StatusErrInternal
			nres.Message = fmt.Sprintf("%s", r)
			err = &Error{Status: nres.Status, Message: nres.Message}
		}
	}()
	if err = h.ServeNRPC(ctx, nreq, nres); err != nil {
		if ne, ok := err.(*Error); ok {
			nres.Status = ne.Status
			nres.Message = ne.Message
		} else {
			nres.Status = StatusErrInternal
			nres.Message = err.Error()
		}
	}
	return
}

var (
	NotFound HandlerFunc = func(ctx context.Context, nreq *Request, nres *Response) (err error) {
		nres.Status = StatusErrNotImplemented
		nres.Message = "service or method not implemented"
		return
	}
)
