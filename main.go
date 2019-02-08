package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/common/log"
)

var (
	logLevel          = flag.String("log-level", "info", "Log level: debug, info (default), warn, error, fatal.")
	udpListenAddress  = flag.String("udp-listen-address", "0.0.0.0:9090", "The address to listen on for udp stats requests.")
	httpListenAddress = flag.String("http-listen-address", "0.0.0.0:9091", "The address to listen on for http stat and telemetry requests.")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())

	logger := log.NewLogger(os.Stdout)

	var err error

	if err = logger.SetLevel(*logLevel); err != nil {
		log.Fatalf(err.Error())
	}

	*udpListenAddress, err = ListenAddress(*udpListenAddress)
	if err != nil {
		log.Fatalf(err.Error())
	}

	*httpListenAddress, err = ListenAddress(*httpListenAddress)
	if err != nil {
		log.Fatalf(err.Error())
	}

	go listenUDP(ctx, logger)
	go listenHTTP(ctx, logger)

	handleSIGTERM(cancel)
}

func handleSIGTERM(cancel func()) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)
	<-sigc
	cancel()
}

// ListenAddress Format a correct listen address
func ListenAddress(s string) (string, error) {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}
