package nrpc

import (
	"fmt"
	"io"
	"regexp"
)

var (
	servicePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
	methodPattern  = regexp.MustCompile(`^[a-z0-9_-]+$`)
)

type Request struct {
	Service  string
	Method   string
	Metadata Metadata
	Payload  interface{}
}

func NewRequest() *Request {
	return &Request{Metadata: Metadata{}}
}

func (r *Request) Validate() (err error) {
	if !servicePattern.MatchString(r.Service) {
		err = fmt.Errorf("invalid service name: %s", r.Service)
		return
	}
	if !methodPattern.MatchString(r.Method) {
		err = fmt.Errorf("invalid method name: %s", r.Method)
		return
	}
	return
}

func (r *Request) Decode(buf []byte) (err error) {
	if buf, err = decodeHeadline(buf, &r.Service, &r.Method); err != nil {
		return
	}
	if buf, err = decodeMetadata(buf, &r.Metadata); err != nil {
		return
	}
	if buf, err = decodePayload(buf, r.Payload); err != nil {
		return
	}
	if err = r.Validate(); err != nil {
		return
	}
	return
}

func (r *Request) Encode(w io.Writer) (err error) {
	if err = r.Validate(); err != nil {
		return
	}
	if _, err = encodeHeadline(w, r.Service, r.Method); err != nil {
		return
	}
	if _, err = encodeMetadata(w, r.Metadata); err != nil {
		return
	}
	if _, err = encodePayload(w, r.Payload); err != nil {
		return
	}
	return
}
