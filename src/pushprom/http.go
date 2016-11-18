package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
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

func listenHTTP() {
	log.Printf("exposing metrics on http://" + *httpListenAddress + "/metrics\n")
	http.Handle("/metrics", prometheus.Handler())

	log.Println("listening for stats on http://" + *httpListenAddress)
	http.HandleFunc("/", httpHandler)
	log.Fatal(http.ListenAndServe(*httpListenAddress, nil))
}
