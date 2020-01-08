package nrpc

import (
	"bufio"
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

func DecodeHeadline(br *bufio.Reader, val1, val2 *string) (err error) {
	var line string
	if line, err = br.ReadString('\n'); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return
	}
	splits := strings.SplitN(line, ",", 2)
	*val1 = strings.TrimSpace(splits[0])
	if len(splits) > 1 {
		*val2 = strings.TrimSpace(splits[1])
	} else {
		*val2 = ""
	}
	return
}

func DecodeMetadata(br *bufio.Reader, m *url.Values) (err error) {
	var line string
	if line, err = br.ReadString('\n'); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		*m = url.Values{}
		return
	}
	var values url.Values
	if values, err = url.ParseQuery(line); err != nil {
		return
	}
	*m = values
	return
}

func DecodePayload(br *bufio.Reader, payload interface{}) error {
	if payload == nil {
		return nil
	}
	dec := json.NewDecoder(br)
	return dec.Decode(payload)
}
