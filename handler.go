package nrpc

import (
	"context"
)

var (
	NotFound = &Handler{
		Alloc: nil,
		Serve: func(ctx context.Context, nreq *Request, nres *Response) (err error) {
			nres.Status = StatusErrNotImplemented
			nres.Message = "service or method not implemented"
			return
		},
	}
)

type Handler struct {
	Alloc func() interface{}
	Serve func(ctx context.Context, nreq *Request, nres *Response) (err error)
}
