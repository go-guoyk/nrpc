package nrpc

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"sync"
)

var (
	ErrServiceNotRegistered = errors.New("service not registered")
)

// RoundTripper abstract the execution of nrpc
type RoundTripper interface {
	// RoundTrip send the Request and receive the Response
	RoundTrip(ctx context.Context, addr string, nreq *Request, out interface{}) (nres *Response, err error)
}

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
		opts.RoundTripper = SimpleTransport
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
	var tried int
	err = backoff.Retry(func() (err error) {
		// non-success is error too
		tried++
		if nres, err = c.roundTripper.RoundTrip(ctx, addr, nreq, out); err == nil {
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
