package main

import (
	"bytes"
	"log"
	"net"
)

func listenUDP() {
	log.Println("listening for stats UDP on port " + *udpListenAddress)
	serverAddr, err := net.ResolveUDPAddr("udp", *udpListenAddress)
	if err != nil {
		log.Println("Error: ", err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Println("Error: ", err)
	}
	defer serverConn.Close()

	buf := make([]byte, 8192)

	for {
		n, _, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error: ", err)
			return
		}
		udpPacketCount.Inc()
		if *debug {
			log.Printf("new udp package: %s", string(buf[0:n]))
		}
		delta, err := NewDelta(bytes.NewBuffer(buf[0:n]))
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		err = delta.Apply()
		if err != nil {
			log.Println("Error: ", err)
		}

	}
}
