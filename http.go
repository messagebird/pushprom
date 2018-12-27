package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

func listenHTTP(ctx context.Context, log plog.Logger) {
	log.Infof("exposing metrics on http://" + *httpListenAddress + "/metrics\n")

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	log.Info("listening for stats on http://" + *httpListenAddress)

	// Instrument the handlers with all the metrics, injecting the "handler"
	// label by currying.
	postHandler := promhttp.InstrumentHandlerDuration(httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": "post"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, httpHandler{
			log: log,
		}),
	)

	mux.HandleFunc("/", postHandler)

	server := http.Server{Addr: *httpListenAddress, Handler: mux}

	go func() {
		<-ctx.Done()

		servCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		server.Shutdown(servCtx)
		cancel()
	}()

	log.Fatal(server.ListenAndServe())
}
