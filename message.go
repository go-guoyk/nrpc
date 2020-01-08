package nrpc

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"strings"
)

var (
	ErrInvalidHeadline     = errors.New("invalid headline")
	ErrInvalidMetadata     = errors.New("invalid metadata")
	ErrNotAIncomingMessage = errors.New("not a incoming message")
	ErrNotAOutgoingMessage = errors.New("not a outgoing message")
	ErrPerformed           = errors.New("message already performed")
)

type Message struct {
	Title    string
	Subtitle string
	Metadata url.Values

	br *bufio.Reader
	bw *bufio.Writer

	performed bool
}

func (m *Message) Performed() bool {
	return m.performed
}

func (m *Message) Recv(body interface{}) (err error) {
	if m.performed {
		err = ErrPerformed
		return
	}
	if m.br == nil {
		err = ErrNotAIncomingMessage
		return
	}
	if body == nil {
		return
	}
	dec := json.NewDecoder(m.br)
	err = dec.Decode(body)
	return
}

func (m *Message) Send(body interface{}) (err error) {
	if m.performed {
		err = ErrPerformed
		return
	}
	if m.bw == nil {
		err = ErrNotAOutgoingMessage
		return
	}
	if err = m.sendHead(); err != nil {
		return
	}
	if body == nil {
		if err = m.bw.WriteByte('\n'); err != nil {
			return
		}
		if err = m.bw.Flush(); err != nil {
			return
		}
		return
	}
	enc := json.NewEncoder(m.bw)
	if err = enc.Encode(body); err != nil {
		return
	}
	err = m.bw.Flush()
	return
}

func (m *Message) sendHead() (err error) {
	if _, err = m.bw.WriteString(m.Title); err != nil {
		return
	}
	if err = m.bw.WriteByte(','); err != nil {
		return
	}
	if _, err = m.bw.WriteString(m.Subtitle); err != nil {
		return
	}
	if err = m.bw.WriteByte('\n'); err != nil {
		return
	}
	for k, vs := range m.Metadata {
		for _, v := range vs {
			if _, err = m.bw.WriteString(strings.ToLower(k)); err != nil {
				return
			}
			if _, err = m.bw.Write([]byte{':', ' '}); err != nil {
				return
			}
			if _, err = m.bw.WriteString(v); err != nil {
				return
			}
			if err = m.bw.WriteByte('\n'); err != nil {
				return
			}
		}
	}
	if err = m.bw.WriteByte('\n'); err != nil {
		return
	}
	return
}

func (m *Message) readline() (line string, err error) {
	if line, err = m.br.ReadString('\n'); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return
	}
	return
}

func NewIncomingMessage(r io.Reader) (m *Message, err error) {
	m = &Message{
		Metadata: url.Values{},
		br:       bufio.NewReader(r),
	}

	var line string
	if line, err = m.readline(); err != nil {
		return
	}

	sp := strings.SplitN(line, ",", 2)
	if len(sp) != 2 {
		err = ErrInvalidHeadline
		return
	}
	m.Title, m.Subtitle = strings.TrimSpace(sp[0]), strings.TrimSpace(sp[1])
	if len(m.Title) == 0 || len(m.Subtitle) == 0 {
		err = ErrInvalidHeadline
		return
	}

	// metadata
	m.Metadata = url.Values{}

	for {
		if line, err = m.readline(); err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			break
		}
		sp = strings.SplitN(line, ":", 2)
		if len(sp) != 2 {
			err = ErrInvalidMetadata
			return
		}
		m.Metadata.Add(strings.ToLower(strings.TrimSpace(sp[0])), strings.TrimSpace(sp[1]))
	}

	return
}

func NewOutgoingMessage(w io.Writer) *Message {
	return &Message{Metadata: url.Values{}, bw: bufio.NewWriter(w)}
}
