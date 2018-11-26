package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
)

func TestHTTP(t *testing.T) {

	delta := &Delta{
		Type:   GAUGE,
		Method: "set",
		Name:   "tree_size",
		Help:   "the size in meters of the tree",
		Value:  8.2,
	}

	out, _ := json.Marshal(delta)
	buf := bytes.NewBuffer(out)

	ts := httptest.NewServer(http.Handler(httpHandler{
		log: log.NewLogger(os.Stdout),
	}))
	defer ts.Close()

	res, err := http.Post(ts.URL, "application/json", buf)
	if err != nil {
		log.Fatal(err)
	}
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)

	metrics := fetchMetrics(t)
	result, err := readMetric(metrics, delta.Name, delta.Labels)
	if assert.Nil(t, err) {
		assert.Equal(t, delta.Value, result)
	}
}
