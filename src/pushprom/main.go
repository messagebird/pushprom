package main

import (
	"flag"
	"log"
)

var (
	debug             = flag.Bool("debug", false, "Log debugging messages.")
	udpListenAddress  = flag.String("udp-listen-address", ":9090", "The address to listen on for udp stats requests.")
	httpListenAddress = flag.String("http-listen-address", ":9091", "The address to listen on for http stat and telemetry requests.")
)

func main() {
	flag.Parse()
	// print more info on log. like line number.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go listenUDP()
	listenHTTP()
}
