package main

import (
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"github.com/prometheus/common/log"

	"github.com/stretchr/testify/assert"
)

var conn *net.UDPConn

func TestMain(m *testing.M) {
	*udpListenAddress = ":3007"

	logger := log.NewNopLogger()

	go listenUDP(logger)

	// wait for it to "boot"
	time.Sleep(time.Millisecond * 1000)

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1"+*udpListenAddress)
	if err != nil {
		log.Fatal("Could not resolve server UPD address")
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("Could not resolve local UDP address")
	}

	conn, err = net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		log.Fatal("Could not establish UDP connection")
	}

	defer conn.Close()

	os.Exit(m.Run())
}

func TestUDP(t *testing.T) {
	delta := &Delta{
		Type:   GAUGE,
		Method: "set",
		Name:   "tree_width",
		Help:   "the width in meters of the tree",
		Value:  2.2,
	}

	buf, _ := json.Marshal(delta)

	_, err := conn.Write(buf)
	assert.Nil(t, err)

	time.Sleep(time.Millisecond * 500)

	metrics := fetchMetrics(t)
	result, err := readMetric(metrics, delta.Name, delta.Labels)
	if assert.Nil(t, err) {
		assert.Equal(t, delta.Value, result)
	}
}

func TestIncorrectJson(t *testing.T) {
	// First, let's write the correct value
	oldDelta := &Delta{
		Type:   GAUGE,
		Method: "set",
		Name:   "tree_width",
		Help:   "the width in meters of the tree",
		Value:  2.2,
	}
	buf, _ := json.Marshal(oldDelta)

	_, err := conn.Write(buf)
	assert.Nil(t, err)

	oldMetrics := fetchMetrics(t)
	oldResult, err := readMetric(oldMetrics, oldDelta.Name, oldDelta.Labels)
	if assert.Nil(t, err) {
		assert.Equal(t, oldDelta.Value, oldResult)
	}

	// Now, let's write the incorrect value
	buf = []byte("hello")

	_, err = conn.Write(buf)
	assert.Nil(t, err)

	// Last, let's write a new value
	newDelta := &Delta{
		Type:   GAUGE,
		Method: "set",
		Name:   "tree_width",
		Help:   "the width in meters of the tree",
		Value:  2.5,
	}
	buf, _ = json.Marshal(newDelta)

	_, err = conn.Write(buf)
	assert.Nil(t, err)

	newMetrics := fetchMetrics(t)
	newResult, err := readMetric(newMetrics, newDelta.Name, newDelta.Labels)
	if assert.Nil(t, err) {
		assert.Equal(t, newDelta.Value, newResult)
	}
}
