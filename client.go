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
	err = backoff.Retry(func() (err error) {
		resp, err = Invoke(ctx, addr, req, out)
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
