package nrpc

import (
	"bufio"
	"encoding/json"
	"io"
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

func DecodeMetadata(br *bufio.Reader, mOut *Metadata) (err error) {
	var line string
	if line, err = br.ReadString('\n'); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		*mOut = Metadata{}
		return
	}
	var m Metadata
	if m, err = ParseMetadata(line); err != nil {
		return
	}
	*mOut = m
	return
}

func DecodePayload(br *bufio.Reader, payload interface{}) (err error) {
	if payload == nil {
		_, err = br.ReadBytes('\n')
		return
	}
	dec := json.NewDecoder(br)
	if err = dec.Decode(payload); err != nil {
		return
	}
	r := io.MultiReader(dec.Buffered(), br)
	nbr := bufio.NewReaderSize(r, 10)
	_, err = nbr.ReadBytes('\n')
	return
}
