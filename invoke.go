package nrpc

import (
	"context"
	"go.guoyk.net/trackid"
	"golang.org/x/sync/errgroup"
	"net"
)

func Invoke(ctx context.Context, addr string, req *Request, out interface{}) (resp *Response, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", addr); err != nil {
		return
	}
	defer conn.Close()

	req.Metadata.Set(MetadataKeyTrackId, trackid.Get(ctx))

	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		return do(ctx, func() (err error) {
			_, err = req.WriteTo(conn)
			return
		})
	})

	eg.Go(func() error {
		return do(ctx, func() (err error) {
			resp, err = ReadResponse(conn)
			return
		})
	})

	if err = eg.Wait(); err != nil {
		return
	}

	if err = resp.Unmarshal(out); err != nil {
		return
	}
	return
}
