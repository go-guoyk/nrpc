package nrpc

import (
	"context"
	"fmt"
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

func InvokeHandler(ctx context.Context, h *Handler, nreq *Request, nres *Response) (err error) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			nres.Status = StatusErrInternal
			nres.Message = fmt.Sprintf("%s", r)
			err = &Error{Status: nres.Status, Message: nres.Message}
		}
	}()
	// payload
	var p interface{}
	if h.Alloc != nil {
		p = h.Alloc()
	}
	if err = nreq.Unmarshal(p); err != nil {
		return
	}
	nreq.Payload = p
	// execute
	if err = h.Serve(ctx, nreq, nres); err != nil {
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
