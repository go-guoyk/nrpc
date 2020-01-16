package nrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
)

func decodeLine(buf []byte) (line []byte, ret []byte, err error) {
	ret = buf
	if i := bytes.IndexByte(buf, '\n'); i < 0 {
		err = io.ErrUnexpectedEOF
	} else {
		line, ret = bytes.TrimSpace(buf[:i]), buf[i+1:]
	}
	return
}

func decodeHeadline(buf []byte, outs ...*string) (ret []byte, err error) {
	var line []byte
	if line, ret, err = decodeLine(buf); err != nil {
		return
	}
	subs := bytes.SplitN(line, []byte{','}, len(outs))
	for i, s := range outs {
		if i >= len(subs) {
			*s = ""
		} else {
			*s = string(bytes.TrimSpace(subs[i]))
		}
	}
	return
}

func decodeMetadata(buf []byte, out *Metadata) (ret []byte, err error) {
	var line []byte
	if line, ret, err = decodeLine(buf); err != nil {
		return
	}
	var m Metadata
	if m, err = ParseMetadata(line); err != nil {
		return
	}
	*out = m
	return
}

func decodePayload(buf []byte, out interface{}) (ret []byte, err error) {
	var line []byte
	if line, ret, err = decodeLine(buf); err != nil {
		return
	}
	if out == nil {
		return
	}
	if err = json.Unmarshal(line, out); err != nil {
		return
	}
	return
}

func decodeMessage(buf []byte, out *string) (ret []byte, err error) {
	var line []byte
	if line, ret, err = decodeLine(buf); err != nil {
		return
	}
	if *out, err = url.QueryUnescape(string(line)); err != nil {
		return
	}
	return
}
