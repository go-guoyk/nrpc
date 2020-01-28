package nrpc

import (
	"net"
	"net/http"
	"time"
)

type ClientOptions struct {
	MaxRetries int
	Timeout    time.Duration
}

type Client struct {
	maxRetries int
	client     *http.Client
	svcs       map[string]string
}

func NewClient(opts ClientOptions) *Client {
	if opts.MaxRetries < 0 {
		opts.MaxRetries = 3
	}
	if opts.Timeout == 0 {
		opts.Timeout = time.Second * 5
	}
	return &Client{
		maxRetries: opts.MaxRetries,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{Timeout: opts.Timeout}).DialContext,
			},
		},
		svcs: map[string]string{},
	}
}

func (c *Client) Register(service, host string) {
	c.svcs[service] = host
}

func (c *Client) Query(service, method string) *Call {
	return c.Call(service, method, false)
}

func (c *Client) Command(service, method string) *Call {
	return c.Call(service, method, true)
}

func (c *Client) Call(service, method string, command bool) *Call {
	return &Call{
		client:  c.client,
		host:    c.svcs[service],
		service: service,
		method:  method,
		command: command,

		maxRetries: c.maxRetries,
	}
}
