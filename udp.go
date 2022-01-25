package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/messagebird/pushprom/delta"

)

func listenUDP(ctx context.Context) {
	fmt.Println("listening for stats UDP on port " + *udpListenAddress)
	serverAddr, err := net.ResolveUDPAddr("udp", *udpListenAddress)
	if err != nil {
		log.Print(err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Print(err)
	}

	defer func(serverConn *net.UDPConn) {
		err := serverConn.Close()
		if err != nil {
			log.Print(err)
		}
	}(serverConn)

	buf := make([]byte, 8192)

	for {
		// handle cancelled context.
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, _, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			log.Print("Error reading from UDP: ", err)
			continue
		}
		udpPacketCount.Inc()

		fmt.Printf("new udp package: %s\n", string(buf[0:n]))

		newDelta, err := delta.NewDelta(bytes.NewBuffer(buf[0:n]))
		if err != nil {
			log.Print("Error creating delta: ", err)
			continue
		}

		err = newDelta.Apply()
		if err != nil {
			log.Print("Error applying delta: ", err)
		}
	}
}
