package nrpc

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrServiceNotRegistered = errors.New("service not registered")
)

type Client struct {
	services  map[string]string
	servicesL sync.Locker
}

func NewClient() *Client {
	return &Client{
		services:  map[string]string{},
		servicesL: &sync.Mutex{},
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
	return Invoke(ctx, addr, req, out)
}
