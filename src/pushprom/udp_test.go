package main

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUDP(t *testing.T) {

	*udpListenAddress = ":3007"
	go listenUDP()

	// wait for it to "boot"
	time.Sleep(time.Millisecond * 1000)

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1"+*udpListenAddress)
	assert.Nil(t, err)

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	assert.Nil(t, err)

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	assert.Nil(t, err)

	defer conn.Close()

	delta := &Delta{
		Type:   GAUGE,
		Method: "set",
		Name:   "tree_width",
		Help:   "the width in meters of the tree",
		Value:  2.2,
	}

	buf, _ := json.Marshal(delta)

	_, err = conn.Write(buf)
	assert.Nil(t, err)

	time.Sleep(time.Millisecond * 500)

	metrics := fetchMetrics(t)
	result, err := readMetric(metrics, delta.Name, delta.Labels)
	if assert.Nil(t, err) {
		assert.Equal(t, delta.Value, result)
	}
}
