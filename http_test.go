package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/messagebird/pushprom/delta"
	"github.com/messagebird/pushprom/metrics"
	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
)

func TestHTTP(t *testing.T) {

	delta := &delta.Delta{
		Type:   delta.GAUGE,
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

	ms, err := metrics.Fetch()
	assert.Nil(t, err)
	result, err := metrics.Read(ms, delta.Name, delta.Labels)
	if assert.Nil(t, err) {
		assert.Equal(t, delta.Value, result)
	}
}
