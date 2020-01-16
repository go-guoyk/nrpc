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
		buf.WriteString(s)
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
