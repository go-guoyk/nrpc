package nrpc

import (
	"bytes"
	"net/url"
	"sort"
	"strings"
)

type Metadata map[string]string

func ParseMetadata(str []byte) (m Metadata, err error) {
	m = make(Metadata)
	for len(str) > 0 {
		kv := str
		if i := bytes.Index(str, []byte{';'}); i > 0 {
			kv, str = bytes.TrimSpace(str[:i]), str[i+1:]
		} else {
			str = nil
		}
		if i := bytes.Index(kv, []byte{'='}); i > 0 {
			var key, val string
			if val, err = url.QueryUnescape(string(bytes.TrimSpace(kv[i+1:]))); err != nil {
				return
			}
			key = strings.ToLower(string(bytes.TrimSpace(kv[:i])))
			m[key] = val
		} else {
			continue
		}
	}
	return
}

func (m Metadata) Get(key string) string {
	key = strings.ToLower(key)
	if m == nil {
		return ""
	}
	return m[key]
}

func (m Metadata) Set(key string, val string) {
	key = strings.ToLower(key)
	if m == nil {
		return
	} else {
		m[key] = val
	}
}

func (m Metadata) Encode() []byte {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	buf := &bytes.Buffer{}
	for _, k := range ks {
		if buf.Len() > 0 {
			buf.WriteByte(';')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(m[k]))
	}
	return buf.Bytes()
}
