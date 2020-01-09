package nrpc

import (
	"context"
	"go.guoyk.net/trackid"
	"net"
)

func Invoke(ctx context.Context, addr string, req *Request, out interface{}) (resp *Response, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", addr); err != nil {
		return
	}
	defer conn.Close()

	req.Metadata.Set(MetadataKeyTrackId, trackid.Get(ctx))
	req.Metadata.Set(MetadataKeyHostname, hostname)

	go req.WriteTo(conn)

	if err = do(ctx, func() (err error) {
		resp, err = ReadResponse(conn)
		return
	}); err != nil {
		return
	}

	if err = resp.Unmarshal(out); err != nil {
		return
	}
	return
}
