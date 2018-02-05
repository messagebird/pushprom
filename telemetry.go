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
)

func init() {
	prometheus.MustRegister(udpPacketCount)
	prometheus.MustRegister(httpRequestCount)
}
