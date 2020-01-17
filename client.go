package nrpc

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"go.guoyk.net/trackid"
	"sync"
)

var (
	ErrServiceNotRegistered = errors.New("service not registered")
)

type ClientOptions struct {
	MaxRetries   uint64
	RoundTripper RoundTripper
}

type Client struct {
	roundTripper RoundTripper
	maxRetries   uint64
	services     map[string]string
	servicesL    sync.Locker
}

func NewClient(opts ClientOptions) *Client {
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 3
	}
	if opts.RoundTripper == nil {
		opts.RoundTripper = DefaultTransport
	}
	return &Client{
		roundTripper: opts.RoundTripper,
		maxRetries:   opts.MaxRetries,
		services:     map[string]string{},
		servicesL:    &sync.Mutex{},
	}
}

func (c *Client) Register(service string, addr string) {
	c.servicesL.Lock()
	defer c.servicesL.Unlock()
	c.services[service] = addr
}

func (c *Client) Invoke(ctx context.Context, nreq *Request, out interface{}) (nres *Response, err error) {
	addr := c.services[nreq.Service]
	if len(addr) == 0 {
		err = ErrServiceNotRegistered
		return
	}

	nreq.Metadata.Set(MetadataKeyTrackId, trackid.Get(ctx))
	nreq.Metadata.Set(MetadataKeyHostname, hostname)

	var tried int
	err = backoff.Retry(func() (err error) {
		// non-success is error too
		tried++
		nres = &Response{}
		nres.Payload = out
		if err = c.roundTripper.RoundTrip(ctx, addr, nreq, nres); err == nil {
			if nres.Status != StatusOK {
				err = &Error{Status: nres.Status, Message: nres.Message, Tried: tried}
			}
		}
		return
	},
		backoff.WithContext(
			backoff.WithMaxRetries(backoff.NewExponentialBackOff(),
				c.maxRetries,
			),
			ctx,
		),
	)
	return
}
