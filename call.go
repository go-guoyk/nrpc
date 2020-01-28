package nrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"go.guoyk.net/trackid"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Call struct {
	client  *http.Client
	host    string
	service string
	method  string
	command bool
	in      interface{}
	out     interface{}

	maxRetries int
}

func (c *Call) In(ptr interface{}) *Call {
	c.in = ptr
	return c
}

func (c *Call) Out(ptr interface{}) *Call {
	c.out = ptr
	return c
}

func (c *Call) Do(ctx context.Context) (err error) {
	if c.host == "" {
		err = UserError(fmt.Errorf("unknown host for service '%s'", c.service))
		return
	}

	// build url
	uri := &url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "/" + c.service + "/" + c.method,
	}

	// build request
	var req *http.Request

	if c.command {
		if c.in == nil {
			err = UserError(errors.New("request with type 'command' must have body"))
			return
		}
		var buf []byte
		if buf, err = json.Marshal(c.in); err != nil {
			err = UserError(err)
			return
		}
		body := bytes.NewReader(buf)
		if req, err = http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), body); err != nil {
			err = UserError(err)
			return
		}
		req.Header.Set(headerContentType, mimeApplicationJSONCharsetUTF8)
	} else {
		if c.in != nil {
			var q url.Values
			enc := form.NewEncoder()
			enc.SetTagName("query")
			if q, err = enc.Encode(c.in); err != nil {
				err = UserError(err)
				return
			}
			uri.RawQuery = q.Encode()
		}
		if req, err = http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil); err != nil {
			err = UserError(err)
			return
		}
	}

	// correlation id
	req.Header.Set(headerCorrelationID, trackid.Get(ctx))

	// execute request
	var resp *http.Response
	if resp, err = c.client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	// on error
	if resp.StatusCode >= 400 {
		sb := &strings.Builder{}
		if _, err = io.Copy(sb, resp.Body); err != nil {
			return
		}
		if resp.StatusCode >= 400 && resp.StatusCode < 400 {
			err = UserError(errors.New(sb.String()))
		} else {
			err = errors.New(sb.String())
		}
		return
	}

	// if out is required
	if c.out != nil {
		if !strings.HasPrefix(resp.Header.Get(headerContentType), mimeApplicationJSON) {
			err = errors.New("not a application/json response")
			return
		}
		dec := json.NewDecoder(resp.Body)
		if err = dec.Decode(c.out); err != nil {
			return
		}
	}

	return
}
