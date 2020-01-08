package nrpc

import (
	"encoding/json"
	"io"
	"net/url"
)

type countableWriter struct {
	w io.Writer
	n *int
}

func (w *countableWriter) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)
	*w.n += n
	return
}

func EncodeHeadline(w io.Writer, subject1, subject2 string) (int, error) {
	buf := make([]byte, 0, len(subject1)+len(subject2)+2)
	buf = append(buf, []byte(subject1)...)
	buf = append(buf, ',')
	buf = append(buf, []byte(subject2)...)
	buf = append(buf, '\n')
	return w.Write(buf)
}

func EncodeMetadata(w io.Writer, metadata url.Values) (int, error) {
	buf := []byte(metadata.Encode())
	buf = append(buf, '\n')
	return w.Write(buf)
}

func EncodePayload(w io.Writer, payload interface{}) (n int, err error) {
	if payload == nil {
		return w.Write([]byte{'\n'})
	}
	wc := &countableWriter{w: w, n: &n}
	enc := json.NewEncoder(wc)
	err = enc.Encode(payload)
	return
}
