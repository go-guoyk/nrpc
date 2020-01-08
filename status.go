package nrpc

const (
	StatusOK                = "ok"
	StatusErrInternal       = "err_internal"
	StatusErrNotFound       = "err_not_found"
	StatusErrNotImplemented = "err_not_implemented"
)

var (
	DefaultMessages = map[string]string{
		StatusOK:                "ok",
		StatusErrInternal:       "internal",
		StatusErrNotFound:       "not found",
		StatusErrNotImplemented: "not implemented",
	}
)
