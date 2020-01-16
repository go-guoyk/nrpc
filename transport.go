package nrpc

import (
	"bufio"
	"bytes"
	"context"
	"net"
	"sync"
)

type conn struct {
	c  net.Conn
	br *bufio.Reader
}

// RoundTripper abstract the execution of nrpc
type RoundTripper interface {
	// RoundTrip send the Request and receive the Response
	RoundTrip(ctx context.Context, addr string, nreq *Request, nres *Response) (err error)
}

var DefaultTransport RoundTripper = NewTransport()

type Transport struct {
	bufPool *sync.Pool
	dialer  *net.Dialer
}

func NewTransport() *Transport {
	return &Transport{
		bufPool: &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
		dialer: &net.Dialer{},
	}
}

func (t *Transport) retrieveConn(ctx context.Context, addr string) (co *conn, ret func(close bool), err error) {
	// TODO: implements connection reuse
	var c net.Conn
	if c, err = t.dialer.DialContext(ctx, "tcp", addr); err != nil {
		return
	}
	ret = func(close bool) {
		_ = c.Close()
	}
	co = &conn{c: c, br: bufio.NewReader(c)}
	return
}

func (t *Transport) retrieveBuf() (buf *bytes.Buffer, ret func()) {
	buf = t.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	ret = func() {
		t.bufPool.Put(buf)
	}
	return
}

func (t *Transport) RoundTrip(ctx context.Context, addr string, nreq *Request, nres *Response) (err error) {
	var conn *conn
	var retConn func(close bool)
	if conn, retConn, err = t.retrieveConn(ctx, addr); err != nil {
		return
	}

	// encode
	wBuf, retWBuf := t.retrieveBuf()
	if err = nreq.Encode(wBuf); err != nil {
		retConn(false)
		retWBuf()
		return
	}

	// write
	if _, err = conn.c.Write(wBuf.Bytes()); err != nil {
		retConn(true)
		retWBuf()
		return
	}
	retWBuf()

	// read three lines
	var rBuf []byte
	for i := 0; i < 3; i++ {
		var buf []byte
		if buf, err = conn.br.ReadBytes('\n'); err != nil {
			retConn(true)
			return
		}
		rBuf = append(rBuf, buf...)
	}
	retConn(false)

	// decode
	if err = nres.Decode(rBuf); err != nil {
		return
	}

	return
}
