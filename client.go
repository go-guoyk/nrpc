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

type ClientOptions struct {
	MaxRetries uint64
}

type Client struct {
	maxRetries uint64
	services   map[string]string
	servicesL  sync.Locker
}

func NewClient(opts ClientOptions) *Client {
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 3
	}
	return &Client{
		maxRetries: opts.MaxRetries,
		services:   map[string]string{},
		servicesL:  &sync.Mutex{},
	}
}

func (c *Client) Register(service string, addr string) {
	c.servicesL.Lock()
	defer c.servicesL.Unlock()
	c.services[service] = addr
}

func (c *Client) Invoke(ctx context.Context, req *Request, out interface{}) (resp *Response, err error) {
	addr := c.services[req.Service]
	if len(addr) == 0 {
		err = ErrServiceNotRegistered
		return
	}
	var tried int
	err = backoff.Retry(func() (err error) {
		// non-success is error too
		tried++
		if resp, err = Invoke(ctx, addr, req, out); err == nil {
			if resp.Status != StatusOK {
				err = &Error{Status: resp.Status, Message: resp.Message, Tried: tried}
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
