package nrpc

import (
	"errors"
	"net"
	"net/url"
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

func (c *Client) Invoke(service, method string, metadata url.Values, body interface{}, out interface{}) (resp *Message, err error) {
	addr := c.services[service]
	if len(addr) == 0 {
		err = ErrServiceNotRegistered
		return
	}
	return Invoke(addr, service, method, metadata, body, out)
}

func Invoke(addr, service, method string, metadata url.Values, body interface{}, out interface{}) (resp *Message, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", addr); err != nil {
		return
	}
	defer conn.Close()
	m := NewOutgoingMessage(conn)
	m.Subject = service
	m.SecondarySubject = method
	if metadata != nil {
		m.Metadata = metadata
	}
	if err = m.Send(body); err != nil {
		return
	}
	if resp, err = NewIncomingMessage(conn); err != nil {
		return
	}
	err = resp.Recv(out)
	return
}
