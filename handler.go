package nrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.guoyk.net/trackid"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	typeContext = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeError   = reflect.TypeOf((*error)(nil)).Elem()
)

const (
	charsetUTF8                    = "charset=utf-8"
	mimeTextPlain                  = "text/plain"
	mimeTextPlainCharsetUTF8       = mimeTextPlain + "; " + charsetUTF8
	mimeApplicationJSON            = "application/json"
	mimeApplicationJSONCharsetUTF8 = mimeApplicationJSON + "; " + charsetUTF8

	headerCorrelationID = "X-Correlation-Id"
	headerContentType   = "Content-Type"
	headerContentLength = "Content-Length"
)

type CountableReader struct {
	Reader io.Reader
	Total  int64
}

func (c *CountableReader) Read(p []byte) (n int, err error) {
	n, err = c.Reader.Read(p)
	atomic.AddInt64(&c.Total, int64(n))
	return
}

type Handler struct {
	svc string
	mtd string
	tgt interface{}
	fn  reflect.Value
	in  reflect.Type
}

func checkRPCFunc(t reflect.Type) (in reflect.Type, ok bool) {
	if t.NumIn() == 2 {
		if !typeContext.AssignableTo(t.In(1)) {
			return
		}
	} else if t.NumIn() == 3 {
		if !typeContext.AssignableTo(t.In(1)) {
			return
		}
		t1 := t.In(2)
		if t1.Kind() != reflect.Ptr {
			return
		}
		if t1.Elem().Kind() != reflect.Struct {
			return
		}
		in = t1.Elem()
	} else {
		return
	}
	if t.NumOut() == 1 {
		if !t.Out(0).AssignableTo(typeError) {
			return
		}
	} else if t.NumOut() == 2 {
		if t.Out(0).Kind() != reflect.Struct {
			return
		}
		if !t.Out(1).AssignableTo(typeError) {
			return
		}
	} else {
		return
	}
	ok = true
	return
}

// ExtractHandlers create a Map of *Handler based on receiver's methods
// supported signatures:
//  - Method1(ctx context.Context) (err error)
//  - Method2(ctx context.Context, in *SomeStruct1) (err error)
//  - Method3(ctx context.Context, in *SomeStruct1) (out SomeStruct2, err error)
//  - Method4(ctx context.Context) (out SomeStruct2, err error)
func ExtractHandlers(name string, tgt interface{}) map[string]*Handler {
	ret := map[string]*Handler{}
	t := reflect.TypeOf(tgt)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if in, ok := checkRPCFunc(m.Type); ok {
			ret[m.Name] = &Handler{svc: name, mtd: m.Name, tgt: tgt, fn: m.Func, in: in}
		}
	}
	return ret
}

func (h *Handler) meterRequestSize(size int64) {
	metricsRequestSizeBytes.WithLabelValues(
		h.svc,
		h.mtd,
	).Observe(float64(size))
}

func (h *Handler) meterRequestsTotal(failed, solid bool) {
	metricsRequestsTotal.WithLabelValues(
		h.svc,
		h.mtd,
		strconv.FormatBool(failed),
		strconv.FormatBool(solid),
	).Inc()
}

func (h *Handler) meterResponseSize(size int) {
	metricsResponseSizeBytes.WithLabelValues(
		h.svc,
		h.mtd,
	).Observe(float64(size))
}

func (h *Handler) meterRequestDuration() func() {
	start := time.Now()
	return func() {
		seconds := float64(time.Since(start)/time.Millisecond) / float64(1000)
		metricsRequestDurationSeconds.WithLabelValues(
			h.svc,
			h.mtd,
		).Observe(seconds)
	}
}

func (h *Handler) sendError(ctx context.Context, rw http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	solid := IsSolid(err)
	if solid {
		code = http.StatusBadRequest
	}
	buf := []byte(err.Error())
	rw.Header().Set(headerContentType, mimeTextPlainCharsetUTF8)
	rw.Header().Set(headerContentLength, strconv.Itoa(len(buf)))
	rw.WriteHeader(code)
	_, _ = rw.Write(buf)

	h.meterResponseSize(len(buf))
	h.meterRequestsTotal(true, solid)

	log.Error().Err(err).Str(
		"topic", "nrpc_access",
	).Str(
		"service", h.svc,
	).Str(
		"method", h.mtd,
	).Bool(
		"failed", true,
	).Bool(
		"solid", solid,
	).Str(
		"crid", trackid.Get(ctx),
	).Msg("")
}

func (h *Handler) sendBody(ctx context.Context, rw http.ResponseWriter, body interface{}) {
	if body == nil {
		rw.WriteHeader(http.StatusOK)
	} else {
		if buf, err := json.Marshal(body); err != nil {
			h.sendError(ctx, rw, err)
			return // return early prevent double metrics
		} else {
			rw.Header().Set(headerContentType, mimeApplicationJSONCharsetUTF8)
			rw.Header().Set(headerContentLength, strconv.Itoa(len(buf)))
			_, _ = rw.Write(buf)

			h.meterResponseSize(len(buf))
		}
	}

	h.meterRequestsTotal(false, false)

	log.Info().Str(
		"topic", "nrpc_access",
	).Str(
		"service", h.svc,
	).Str(
		"method", h.mtd,
	).Bool(
		"failed", false,
	).Bool(
		"solid", false,
	).Str(
		"crid", trackid.Get(ctx),
	).Msg("")
}

func (h *Handler) sendValues(ctx context.Context, rw http.ResponseWriter, rets []reflect.Value) {
	var err error
	var out interface{}
	if len(rets) == 1 {
		if !rets[0].IsNil() {
			err = rets[0].Interface().(error)
		}
	} else {
		out = rets[0].Interface()
		if !rets[1].IsNil() {
			err = rets[1].Interface().(error)
		}
	}
	if err != nil {
		h.sendError(ctx, rw, err)
	} else {
		h.sendBody(ctx, rw, out)
	}
}

func (h *Handler) buildArgs(ctx context.Context, req *http.Request) (args []reflect.Value, err error) {
	args = []reflect.Value{reflect.ValueOf(h.tgt), reflect.ValueOf(ctx)}
	if h.in != nil {
		v := reflect.New(h.in).Interface()
		if req.Method == http.MethodGet {
			dec := form.NewDecoder()
			dec.SetTagName("query")
			if err = dec.Decode(v, req.URL.Query()); err != nil {
				err = Solid(err)
				return
			}
		} else if req.Method == http.MethodPost {
			cr := &CountableReader{Reader: req.Body}
			dec := json.NewDecoder(cr)
			if err = dec.Decode(v); err != nil {
				err = Solid(err)
				return
			}
			h.meterRequestSize(cr.Total)
		} else {
			err = Solid(fmt.Errorf("invalid http method: %s", req.Method))
			return
		}
		// defaults
		if err = defaults.Set(v); err != nil {
			err = Solid(err)
			return
		}
		// validate
		val := validator.New()
		if err = val.Struct(v); err != nil {
			err = Solid(err)
			return
		}
		args = append(args, reflect.ValueOf(v))
	}
	return
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// duration
	commit := h.meterRequestDuration()
	defer commit()

	// correlation id
	ctx := trackid.Set(req.Context(), req.Header.Get(headerCorrelationID))
	rw.Header().Set(headerCorrelationID, trackid.Get(ctx))

	args, err := h.buildArgs(ctx, req)
	if err != nil {
		h.sendError(ctx, rw, err)
		return
	}

	rets := h.fn.Call(args)

	h.sendValues(ctx, rw, rets)
}
