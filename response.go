package nrpc

import (
	"fmt"
	"io"
	"regexp"
)

var (
	statusPattern = regexp.MustCompile(`^[a-z0-9_-]+$`)
)

type Response struct {
	Status   string
	Metadata Metadata
	Message  string
	Payload  interface{}
}

func NewResponse() *Response {
	return &Response{Metadata: Metadata{}}
}

func (r *Response) Validate() (err error) {
	if !statusPattern.MatchString(r.Status) {
		err = fmt.Errorf("invalid status: %s", r.Status)
		return
	}
	return
}

func (r *Response) Decode(buf []byte) (err error) {
	if buf, err = decodeHeadline(buf, &r.Status); err != nil {
		return
	}
	if buf, err = decodeMetadata(buf, &r.Metadata); err != nil {
		return
	}
	if r.Status == StatusOK {
		if buf, err = decodePayload(buf, r.Payload); err != nil {
			return
		}
	} else {
		if buf, err = decodeMessage(buf, &r.Message); err != nil {
			return
		}
	}
	if err = r.Validate(); err != nil {
		return
	}
	return
}

// WriteTo serialize the response into io.Writer
func (r *Response) Encode(w io.Writer) (err error) {
	if err = r.Validate(); err != nil {
		return
	}
	if _, err = encodeHeadline(w, r.Status); err != nil {
		return
	}
	if _, err = encodeMetadata(w, r.Metadata); err != nil {
		return
	}
	if r.Status == StatusOK {
		if _, err = encodePayload(w, r.Payload); err != nil {
			return
		}
	} else {
		if _, err = encodeMessage(w, r.Message); err != nil {
			return
		}
	}
	return
}
