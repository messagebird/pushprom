package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	debug             = flag.Bool("debug", false, "Log debugging messages.")
	udpListenAddress  = flag.String("udp-listen-address", "0.0.0.0:9090", "The address to listen on for udp stats requests.")
	httpListenAddress = flag.String("http-listen-address", "0.0.0.0:9091", "The address to listen on for http stat and telemetry requests.")
)

func main() {
	flag.Parse()

	// Print more info on the log, like line number.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error

	*udpListenAddress, err = ListenAddress(*udpListenAddress)
	if err != nil {
		log.Fatalf(err.Error())
	}

	*httpListenAddress, err = ListenAddress(*httpListenAddress)
	if err != nil {
		log.Fatalf(err.Error())
	}

	go listenUDP()
	listenHTTP()
}

// ListenAddress Format a correct listen address
func ListenAddress(s string) (string, error) {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}
