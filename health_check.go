package nrpc

import (
	"context"
	"net/http"
	"strconv"
)

type HealthCheck interface {
	HealthCheck(ctx context.Context) error
}

type HealthChecks struct {
	hcs []HealthCheck
}

func (hcs *HealthChecks) Add(hc HealthCheck) {
	hcs.hcs = append(hcs.hcs, hc)
}

func (hcs *HealthChecks) HealthCheck(ctx context.Context) error {
	for _, hc := range hcs.hcs {
		if err := hc.HealthCheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (hcs *HealthChecks) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if err := hcs.HealthCheck(req.Context()); err != nil {
		buf := []byte(err.Error())
		rw.Header().Set(headerContentType, mimeTextPlainCharsetUTF8)
		rw.Header().Set(headerContentLength, strconv.Itoa(len(buf)))
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write(buf)
	} else {
		buf := []byte("OK")
		rw.Header().Set(headerContentType, mimeTextPlainCharsetUTF8)
		rw.Header().Set(headerContentLength, strconv.Itoa(len(buf)))
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(buf)
	}
}
