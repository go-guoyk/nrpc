package nrpc

import "os"

const (
	MetadataKeyTrackId  = "track_id"
	MetadataKeyHostname = "hostname"
)

var (
	hostname string
)

func init() {
	hostname, _ = os.Hostname()
}
