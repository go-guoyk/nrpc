package nrpc

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"strings"
	"time"
)

type ClientOptions struct {
	MaxRetries int
	Timeout    time.Duration
	Logger     *zerolog.Logger
}

type Client struct {
	maxRetries int
	client     *http.Client
	svcs       map[string]string
	logger     zerolog.Logger
}

func NewClient(opts ClientOptions) *Client {
	if opts.MaxRetries < 0 {
		opts.MaxRetries = 0
	} else if opts.MaxRetries == 0 {
		opts.MaxRetries = 3
	}
	if opts.Timeout == 0 {
		opts.Timeout = time.Second * 5
	}
	if opts.Logger == nil {
		opts.Logger = &log.Logger
	}
	return &Client{
		maxRetries: opts.MaxRetries,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{Timeout: opts.Timeout}).DialContext,
			},
		},
		svcs:   map[string]string{},
		logger: opts.Logger.With().Str("topic", "nrpc-client").Logger(),
	}
}

func (c *Client) Register(service, host string) {
	c.svcs[service] = host
}

func (c *Client) Query(target string, inout ...interface{}) *Call {
	return c.Call(target, false, inout...)
}

func (c *Client) Command(target string, inout ...interface{}) *Call {
	return c.Call(target, true, inout...)
}

func (c *Client) Call(target string, command bool, inout ...interface{}) *Call {
	splits := strings.Split(target, ".")
	if len(splits) != 2 {
		panic("invalid target: '" + target + "', must be in format of 'Service.Method'")
	}
	service, method := splits[0], splits[1]
	host := c.svcs[service]
	if len(host) == 0 {
		panic("no host registered for service: '" + service + "'")
	}
	call := &Call{
		client:  c.client,
		host:    host,
		service: service,
		method:  method,
		command: command,
		logger:  c.logger.With().Str("service", service).Str("method", method).Str("host", host).Logger(),

		maxRetries: c.maxRetries,
	}
	if len(inout) > 0 {
		call.in = inout[0]
	}
	if len(inout) > 1 {
		call.out = inout[1]
	}
	if len(inout) > 2 {
		panic("invalid arguments, 'inout' should have max length of 2")
	}
	return call
}
