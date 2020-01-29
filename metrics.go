package nrpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricsRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "nrpc",
		Name:      "requests_total",
		Help:      "The total number of processed requests",
	}, []string{"service", "method", "failed", "solid"})
	metricsRequestDurationSeconds = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem: "nrpc",
		Name:      "request_duration_seconds",
		Help:      "The total duration in seconds of processed requests",
	}, []string{"service", "method"})
	metricsRequestSizeBytes = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem: "nrpc",
		Name:      "request_size_bytes",
		Help:      "The total size in bytes of processed requests",
	}, []string{"service", "method"})
	metricsResponseSizeBytes = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem: "nrpc",
		Name:      "response_size_bytes",
		Help:      "The total size in bytes of processed responses",
	}, []string{"service", "method"})
)
