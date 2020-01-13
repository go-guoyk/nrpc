package nrpc

import (
	"bufio"
	"errors"
	"io"
	"net/url"
	"strings"
)

var (
	ErrMissingStatus = errors.New("missing status")
)

type Response struct {
	Status   string
	Message  string
	Metadata url.Values

	// outgoing only
	Payload interface{}

	// incoming only
	br *bufio.Reader
}

// NewResponse create a new outgoing response
// outgoing only
func NewResponse() *Response {
	return &Response{
		Status:   StatusOK,
		Metadata: url.Values{},
	}
}

// ReadResponse read a incoming response from io.Reader
// incoming only
func ReadResponse(r io.Reader) (req *Response, err error) {
	req = &Response{}
	br := bufio.NewReader(r)
	if err = DecodeHeadline(br, &req.Status, &req.Message); err != nil {
		return
	}
	if len(req.Status) == 0 {
		err = ErrMissingStatus
		return
	}
	if req.Message, err = url.QueryUnescape(req.Message); err != nil {
		return
	}
	req.Status = strings.ToLower(req.Status)
	if err = DecodeMetadata(br, &req.Metadata); err != nil {
		return
	}
	req.br = br
	return
}

// Unmarshal unmarshal the body
// incoming only
func (r *Response) Unmarshal(body interface{}) error {
	return DecodePayload(r.br, body)
}

// WriteTo serialize the response into io.Writer
// outgoing only
func (r *Response) WriteTo(w io.Writer) (tn int64, err error) {
	if len(r.Status) == 0 {
		err = ErrMissingStatus
		return
	}
	var n int
	if n, err = EncodeHeadline(w, r.Status, url.QueryEscape(r.Message)); err != nil {
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
