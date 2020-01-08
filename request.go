package nrpc

import (
	"bufio"
	"errors"
	"io"
	"net/url"
	"strings"
)

var (
	ErrMissingService = errors.New("missing service")
	ErrMissingMethod  = errors.New("missing method")
)

type Request struct {
	Service  string
	Method   string
	Metadata url.Values

	// outgoing only
	Payload interface{}

	// incoming only
	br *bufio.Reader
}

// NewRequest create a new outgoing request
// outgoing only
func NewRequest(service, method string) *Request {
	return &Request{
		Service:  strings.ToLower(strings.TrimSpace(service)),
		Method:   strings.ToLower(strings.TrimSpace(method)),
		Metadata: url.Values{},
	}
}

// ReadRequest read a incoming request from io.Reader
// incoming only
func ReadRequest(r io.Reader) (req *Request, err error) {
	req = &Request{}
	br := bufio.NewReader(r)
	if err = DecodeHeadline(br, &req.Service, &req.Method); err != nil {
		return
	}
	if len(req.Service) == 0 {
		err = ErrMissingService
		return
	}
	if len(req.Method) == 0 {
		err = ErrMissingMethod
		return
	}
	if err = DecodeMetadata(br, &req.Metadata); err != nil {
		return
	}
	req.br = br
	return
}

// Unmarshal unmarshal the body
// incoming only
func (r *Request) Unmarshal(body interface{}) error {
	return DecodePayload(r.br, body)
}

// WriteTo serialize the request into io.Writer
// outgoing only
func (r *Request) WriteTo(w io.Writer) (tn int64, err error) {
	if len(r.Service) == 0 {
		err = ErrMissingService
		return
	}
	if len(r.Method) == 0 {
		err = ErrMissingMethod
		return
	}
	var n int
	if n, err = EncodeHeadline(w, r.Service, r.Method); err != nil {
		return
	}
	tn += int64(n)
	if n, err = EncodeMetadata(w, r.Metadata); err != nil {
		return
	}
	tn += int64(n)
	if n, err = EncodePayload(w, r.Payload); err != nil {
		return
	}
	tn += int64(n)
	return
}
