package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/messagebird/pushprom/delta"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type httpHandler struct {
}

func (httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	if r.Method == "POST" {
		httpRequestCount.Inc()

		newDelta, err := delta.NewDelta(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		if err = newDelta.Apply(); err != nil {
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

func listenHTTP(wg *sync.WaitGroup, ctx context.Context, stderrLogger *log.Logger, stdoutLogger *log.Logger) {
	stdoutLogger.Println("exposing metrics on http://" + *httpListenAddress + "/metrics")

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	// Instrument the handlers with all the metrics, injecting the "handler"
	// label by currying.
	postHandler := promhttp.InstrumentHandlerDuration(httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": "post"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, httpHandler{}),
	)

	mux.HandleFunc("/", postHandler)

	server := http.Server{
		// see https://ieftimov.com/posts/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		Addr:              *httpListenAddress,
		Handler:           mux,
	}

	go func(wg *sync.WaitGroup) {
		<-ctx.Done()

		stdoutLogger.Println("shutting down http listener on " + *httpListenAddress)

		servCtx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		err := server.Shutdown(servCtx)
		if err != nil {
			stderrLogger.Fatalf("http listener failed to shutdown gracefully: %v", err)
		}
		cancel()
		defer wg.Done()
		stdoutLogger.Println("http listener is now offline")
	}(wg)

	stdoutLogger.Println("listening for stats HTTP on http://" + *httpListenAddress)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		stderrLogger.Fatalf("Failed to ListenAndServe: %v", err)
	}
}
