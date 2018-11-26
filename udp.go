package main

import (
	"bytes"
	"net"

	plog "github.com/prometheus/common/log"
)

func listenUDP(log plog.Logger) {
	log.Info("listening for stats UDP on port " + *udpListenAddress)
	serverAddr, err := net.ResolveUDPAddr("udp", *udpListenAddress)
	if err != nil {
		log.Error(err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Error(err)
	}
	defer serverConn.Close()

	buf := make([]byte, 8192)

	for {
		n, _, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			log.Error("Error reading from UDP: ", err)
			continue
		}
		udpPacketCount.Inc()
		if *debug {
			log.Debugf("new udp package: %s", string(buf[0:n]))
		}
		delta, err := NewDelta(bytes.NewBuffer(buf[0:n]))
		if err != nil {
			log.Error("Error creating delta: ", err)
			continue
		}

		err = delta.Apply()
		if err != nil {
			log.Error("Error applying delta: ", err)
		}
	}
}
