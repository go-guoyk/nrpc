package nrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
)

func encodeHeadline(w io.Writer, subs ...string) (int, error) {
	buf := &bytes.Buffer{}
	for _, s := range subs {
		if buf.Len() > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(url.QueryEscape(s))
	}
	buf.WriteRune('\n')
	return w.Write(buf.Bytes())
}

func encodeMetadata(w io.Writer, metadata Metadata) (int, error) {
	buf := append(metadata.Encode(), '\n')
	return w.Write(buf)
}

func encodePayload(w io.Writer, payload interface{}) (int, error) {
	var buf []byte
	if payload != nil {
		var err error
		if buf, err = json.Marshal(payload); err != nil {
			return 0, err
		}
	}
	return w.Write(append(buf, '\n'))
}

func encodeMessage(w io.Writer, message string) (int, error) {
	return w.Write(append([]byte(url.QueryEscape(message)), '\n'))
}

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
			if *s, err = url.QueryUnescape(string(bytes.TrimSpace(subs[i]))); err != nil {
				return
			}
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
