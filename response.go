package nrpc

import "net/url"

type Response struct {
	Status   string
	Message  string
	Metadata url.Values
}
