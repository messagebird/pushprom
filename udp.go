package main

import (
	"bytes"
	"context"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/messagebird/pushprom/delta"
)

func listenUDP(wg *sync.WaitGroup, ctx context.Context, stderrLogger *log.Logger, stdoutLogger *log.Logger) {
	stdoutLogger.Println("listening for stats UDP on port " + *udpListenAddress)
	serverAddr, err := net.ResolveUDPAddr("udp", *udpListenAddress)
	if err != nil {
		stderrLogger.Print(err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		stderrLogger.Print(err)
	}

	defer func(serverConn *net.UDPConn, wg *sync.WaitGroup) {
		stdoutLogger.Println("closing incoming UDP port " + *udpListenAddress)
		err := serverConn.Close()
		if err != nil {
			stderrLogger.Print(err)
		}
		defer wg.Done()
		stdoutLogger.Println("udp listener is now offline")
	}(serverConn, wg)

	buf := make([]byte, 8192)

	for {
		// handle cancelled context.
		select {
		case <-ctx.Done():
			stdoutLogger.Println("shutting down udp listener on " + *udpListenAddress)
			return
		default:
		}

		err = serverConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		if err != nil {
			stderrLogger.Print("failed setting read UDP deadline: ", err)
		}

		n, _, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			if !strings.Contains(err.Error(), "i/o timeout") {
				stderrLogger.Print("error reading from UDP: ", err)
			}
			continue
		}
		udpPacketCount.Inc()

		// fmt.Printf("new udp package: %s\n", string(buf[0:n]))

		newDelta, err := delta.NewDelta(bytes.NewBuffer(buf[0:n]))
		if err != nil {
			stderrLogger.Print("Error creating delta: ", err)
			continue
		}

		err = newDelta.Apply()
		if err != nil {
			stderrLogger.Print("Error applying delta: ", err)
		}
	}
}
