package nrpc

import (
	"context"
	"go.guoyk.net/trackid"
	"net"
)

func Invoke(ctx context.Context, addr string, nreq *Request, out interface{}) (nres *Response, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", addr); err != nil {
		return
	}
	defer conn.Close()

	nreq.Metadata.Set(MetadataKeyTrackId, trackid.Get(ctx))
	nreq.Metadata.Set(MetadataKeyHostname, hostname)

	go nreq.WriteTo(conn)

	if err = do(ctx, func() (err error) {
		nres, err = ReadResponse(conn)
		return
	}); err != nil {
		return
	}

	// unmarshal only on success
	if nres.Status == StatusOK {
		if err = nres.Unmarshal(out); err != nil {
			return
		}
	}
	return
}
