package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

func listenHTTP(ctx context.Context) {
	log.Println("exposing metrics on http://" + *httpListenAddress + "/metrics\n")

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

	server := http.Server{Addr: *httpListenAddress, Handler: mux}

	go func() {
		<-ctx.Done()

		servCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		err := server.Shutdown(servCtx)
		if err != nil {
			log.Print(err)
		}
		cancel()
	}()

	log.Println("listening for stats on http://" + *httpListenAddress)
	log.Fatal(server.ListenAndServe())
}
