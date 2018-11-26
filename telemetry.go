package main

import "github.com/prometheus/client_golang/prometheus"

var (
	udpPacketCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pushprom_udp_packets_total",
			Help: "The number of packets received.",
		},
	)

	httpRequestCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pushprom_http_requests_total",
			Help: "The number of http requests received.",
		},
	)

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "request_duration_seconds",
			Help: "A histogram of latencies for requests.",
		},
		[]string{"handler", "method"},
	)
)

func init() {
	prometheus.MustRegister(udpPacketCount)
	prometheus.MustRegister(httpRequestCount)
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}
