package nrpc

import (
	"context"
	"net"
)

var DefaultTransport RoundTripper = &Transport{}

type Transport struct{}

func (st *Transport) RoundTrip(ctx context.Context, addr string, nreq *Request, out interface{}) (nres *Response, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", addr); err != nil {
		return
	}
	defer conn.Close()

	go nreq.WriteTo(conn)

	if err = do(ctx, func() (err error) {
		nres, err = ReadResponse(conn)
		return
	}); err != nil {
		return
	}

	// Unmarshal only on success
	if nres.Status == StatusOK {
		if err = nres.Unmarshal(out); err != nil {
			return
		}
	}
	return
}
