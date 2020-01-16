package nrpc

import (
	"io"
)

type Response struct {
	Status   string
	Metadata Metadata
	Message  string
	Payload  interface{}
}

// NewResponse create a new response
func NewResponse() *Response {
	return &Response{
		Metadata: Metadata{},
	}
}

func (p *Response) Decode(buf []byte) (err error) {
	if buf, err = decodeHeadline(buf, &p.Status); err != nil {
		return
	}
	if buf, err = decodeMetadata(buf, &p.Metadata); err != nil {
		return
	}
	if p.Status == StatusOK {
		if buf, err = decodePayload(buf, p.Payload); err != nil {
			return
		}
	} else {
		if buf, err = decodeMessage(buf, &p.Message); err != nil {
			return
		}
	}
	return
}

// WriteTo serialize the response into io.Writer
func (p *Response) Encode(w io.Writer) (err error) {
	if _, err = encodeHeadline(w, p.Status); err != nil {
		return
	}
	if _, err = encodeMetadata(w, p.Metadata); err != nil {
		return
	}
	if p.Status == StatusOK {
		if _, err = encodePayload(w, p.Payload); err != nil {
			return
		}
	} else {
		if _, err = encodeMessage(w, p.Message); err != nil {
			return
		}
	}
	return
}
