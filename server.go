package nrpc

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/pprof"
	"reflect"
)

type ServerOptions struct {
	Addr   string
	Logger *zerolog.Logger
}

type Server struct {
	logger zerolog.Logger
	s      *http.Server
	mux    *http.ServeMux
	hcs    *HealthChecks
}

// Register register a rpc object with default name
func (s *Server) Register(r interface{}) {
	t := reflect.TypeOf(r)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	s.RegisterName(t.Name(), r)
}

// RegisterName register a rpc object with given name
func (s *Server) RegisterName(name string, r interface{}) {
	// add health check
	if hc, ok := r.(HealthCheck); ok {
		s.hcs.Add(hc)
	}
	// extract and add handlers
	hs := ExtractHandlers(name, r)
	for m, h := range hs {
		h.logger = s.logger.With().Str("service", name).Str("method", m).Logger()
		s.mux.Handle("/"+name+"/"+m, h)
	}
}

func NewServer(opts ServerOptions) *Server {
	if opts.Addr == "" {
		opts.Addr = ":3000"
	}
	if opts.Logger == nil {
		opts.Logger = &log.Logger
	}
	mux := http.NewServeMux()
	hcs := &HealthChecks{}
	// mount pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// mount prometheus
	mux.Handle("/metrics", promhttp.Handler())
	// mount healthz
	mux.Handle("/healthz", hcs)
	return &Server{
		s: &http.Server{
			Addr:    opts.Addr,
			Handler: mux,
		},
		hcs:    hcs,
		mux:    mux,
		logger: opts.Logger.With().Str("topic", "nrpc-server").Logger(),
	}
}

func (s *Server) Start(ech chan error) {
	go func() {
		if ech != nil {
			ech <- s.s.ListenAndServe()
		} else {
			_ = s.s.ListenAndServe()
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
