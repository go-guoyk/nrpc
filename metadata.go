package nrpc

import "net/url"

type Metadata map[string]string

func ParseMetadata(str string) (m Metadata, err error) {
	// TODO: don't use url.Values
	m = make(Metadata)
	var vs url.Values
	if vs, err = url.ParseQuery(str); err != nil {
		return
	}
	for k := range vs {
		m[k] = vs.Get(k)
	}
	return
}

func (m Metadata) Get(key string) string {
	if m == nil {
		return ""
	}
	return m[key]
}

func (m Metadata) Set(key string, val string) {
	if m == nil {
		return
	} else {
		m[key] = val
	}
}

func (m Metadata) Encode() string {
	// TODO: don't use url.Values
	if m == nil {
		return ""
	}
	vs := url.Values{}
	for k, v := range m {
		vs.Set(k, v)
	}
	return vs.Encode()
}
