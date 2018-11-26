package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	plog "github.com/prometheus/common/log"
)

type httpHandler struct {
	log plog.Logger
}

func (httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defer r.Body.Close()
		httpRequestCount.Inc()

		delta, err := NewDelta(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		if err = delta.Apply(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
		} else {
			fmt.Fprintf(w, "ok")
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method not allowed")
	}
}

func listenHTTP(log plog.Logger) {
	log.Warnf("exposing metrics on http://" + *httpListenAddress + "/metrics\n")
	http.Handle("/metrics", promhttp.Handler())

	log.Warn("listening for stats on http://" + *httpListenAddress)

	// Instrument the handlers with all the metrics, injecting the "handler"
	// label by currying.
	postHandler := promhttp.InstrumentHandlerDuration(httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": "post"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, httpHandler{
			log: log,
		}),
	)

	http.HandleFunc("/", postHandler)

	log.Fatal(http.ListenAndServe(*httpListenAddress, nil))
}
