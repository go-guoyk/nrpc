package nrpc

import (
	"errors"
	"io"
	"regexp"
)

var (
	servicePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
	methodPattern  = regexp.MustCompile(`^[a-z0-9_-]+$`)
)

var (
	errInvalidServiceName = errors.New("invalid service name")
	errInvalidMethodName  = errors.New("invalid method name")
)

type Request struct {
	Service  string
	Method   string
	Metadata Metadata
	Payload  interface{}
}

// NewRequest create a new request
func NewRequest() *Request {
	return &Request{
		Metadata: Metadata{},
	}
}

func (q *Request) Validate() (err error) {
	if !servicePattern.MatchString(q.Service) {
		err = errInvalidServiceName
		return
	}
	if !methodPattern.MatchString(q.Method) {
		err = errInvalidMethodName
		return
	}
	return
}

func (q *Request) Decode(buf []byte) (err error) {
	if buf, err = decodeHeadline(buf, &q.Service, &q.Method); err != nil {
		return
	}
	if buf, err = decodeMetadata(buf, &q.Metadata); err != nil {
		return
	}
	if buf, err = decodePayload(buf, q.Payload); err != nil {
		return
	}
	if err = q.Validate(); err != nil {
		return
	}
	return
}

func (q *Request) Encode(w io.Writer) (err error) {
	if err = q.Validate(); err != nil {
		return
	}
	if _, err = encodeHeadline(w, q.Service, q.Method); err != nil {
		return
	}
	if _, err = encodeMetadata(w, q.Metadata); err != nil {
		return
	}
	if _, err = encodePayload(w, q.Payload); err != nil {
		return
	}
	return
}
