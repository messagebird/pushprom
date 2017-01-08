package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewDelta(t *testing.T) {
	expected := &Delta{
		Type:   GAUGE,
		Method: "add",
		Name:   "parachutiste",
		Help:   "des oufs qui chutent du ciel",
		Value:  7.77,
		Labels: map[string]string{
			"kind":  "walker",
			"color": "blue",
		},
	}

	out, err := json.Marshal(expected)
	assert.Nil(t, err)

	buf := bytes.NewBuffer(out)

	result, err := NewDelta(buf)
	assert.Nil(t, err)

	assert.Equal(t, expected, result)
}

func fetchMetrics(t *testing.T) map[string]string {
	result := make(map[string]string)
	ts := httptest.NewServer(prometheus.Handler())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	reader := bufio.NewReader(res.Body)
	err = nil
	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		line = strings.Trim(line, "\n")
		parts := strings.Split(line, " ")
		if (parts[0] != "#") && (len(parts) > 1) {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func readMetric(metrics map[string]string, name string, labels prometheus.Labels) (float64, error) {
	// http_request_duration_microseconds{handler="alerts",quantile="0.5"}
	fullName := name
	if len(labels) > 0 {
		pairs := []string{}
		// TODO SORT GOOD
		for k, v := range labels {
			pairs = append(pairs, k+"=\""+v+"\"")
		}
		fullName += "{" + strings.Join(pairs, ",") + "}"
	}

	stringValue := metrics[fullName]
	return strconv.ParseFloat(stringValue, 64)
}

func TestApplyVector(t *testing.T) {

	delta := &Delta{
		Type:   GAUGE,
		Method: "add",
		Name:   "monster",
		Help:   "evil imaginary enemies",
		Value:  67.92,
		Labels: map[string]string{
			"nose_size_meters": "7",
		},
	}

	// set an initial value
	initialValue := 1072.3
	gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: delta.Name,
		Help: delta.Help,
	}, delta.LabelNames())

	err := prometheus.Register(gaugeVec)
	assert.Nil(t, err)

	gauge := gaugeVec.With(delta.Labels)
	gauge.Set(initialValue)

	// apply the delta
	err = delta.Apply()
	assert.Nil(t, err)

	result, err := readMetric(fetchMetrics(t), delta.Name, delta.Labels)
	if assert.Nil(t, err) {
		// check if the delta value was added to the metric
		expected := initialValue + delta.Value
		assert.Equal(t, expected, result)
	}
}

func TestApplyOne(t *testing.T) {
	delta := &Delta{
		Type:   GAUGE,
		Method: "add",
		Name:   "rappers",
		Help:   "gangsta musicians",
		Value:  27.292,
	}
	// set an initial value
	initialValue := 101.4
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: delta.Name,
		Help: delta.Help,
	})
	err := prometheus.Register(gauge)
	assert.Nil(t, err)

	gauge.Set(initialValue)

	// apply the delta
	err = delta.Apply()
	assert.Nil(t, err)

	result, err := readMetric(fetchMetrics(t), delta.Name, delta.Labels)
	if assert.Nil(t, err) {
		// check if the delta value was added to the metric
		expected := initialValue + delta.Value
		assert.Equal(t, expected, result)
	}
}

func TestApplyWithWrongLabels(t *testing.T) {
	// create a metric with one set of labels then try to update it with another. it should break.
	delta := &Delta{
		Type:   GAUGE,
		Method: "add",
		Name:   "dancers",
		Help:   "people that jump pretty",
		Value:  1.22,
		Labels: map[string]string{
			"nose_size_meters": "187",
		},
	}
	err := delta.Apply()
	assert.Nil(t, err)

	delta.Labels["jumpiness"] = "9"

	err = delta.Apply()
	assert.NotNil(t, err)

	delta.Labels = map[string]string{}
	err = delta.Apply()
	assert.NotNil(t, err)
}

func TestMulti(t *testing.T) {

	tests := []struct {
		delta      Delta
		f          func(a, b float64) float64
		firstValue float64
	}{
		/*{
			Delta{
				Type:   COUNTER,
				Method: "set",
			},
			func(before, value float64) float64 {
				return value
			},
			0,
		},*/
		{
			Delta{
				Type:   COUNTER,
				Method: "inc",
			},
			func(before, value float64) float64 {
				return before + 1
			},
			0,
		},
		{
			Delta{
				Type:   COUNTER,
				Method: "add",
			},
			func(before, value float64) float64 {
				return before + value
			},
			0,
		},
		{
			Delta{
				Type:   GAUGE,
				Method: "add",
			},
			func(before, value float64) float64 {
				return before + value
			},
			0,
		},
		{
			Delta{
				Type:   GAUGE,
				Method: "sub",
			},
			func(before, value float64) float64 {
				return before - value
			},
			0,
		},
		{
			Delta{
				Type:   GAUGE,
				Method: "inc",
			},
			func(before, value float64) float64 {
				return before + 1
			},
			0,
		},
		{
			Delta{
				Type:   GAUGE,
				Method: "dec",
			},
			func(before, value float64) float64 {
				return before - 1
			},
			0,
		},
		{
			Delta{
				Type:   GAUGE,
				Method: "set",
			},
			func(before, value float64) float64 {
				return value
			},
			0,
		},
		{
			Delta{
				Type:   HISTOGRAM,
				Method: "observe",
			},
			func(before, value float64) float64 {
				// here we just test the _sum
				return before + value
			},
			0,
		},
		{
			Delta{
				Type:   SUMMARY,
				Method: "observe",
			},
			func(before, value float64) float64 {
				// here we just test the _sum
				return before + value
			},
			0,
		},
	}

	// dup all tests so we test deltas with labels and without
	for i, test := range tests {
		test.delta.Labels = map[string]string{
			fmt.Sprintf("key_%d", i): fmt.Sprintf("val_%d", i),
		}
		tests = append(tests, test)
	}

	// apply all deltas
	for i := 0; i < len(tests); i++ {
		test := &tests[i]
		// stub data
		test.delta.Value = 67.92 + float64(i)
		test.delta.Name = fmt.Sprintf("test_%d", i)
		test.delta.Help = fmt.Sprintf("help_%d", i)
		// apply Delta
		err := test.delta.Apply()
		assert.Nil(t, err)
	}

	initialValue := 0.0 // all metrics start with value zero
	// get metrics and compare results
	metrics := fetchMetrics(t)

	for i := 0; i < len(tests); i++ {
		test := &tests[i]
		// read result and compare
		metricName := test.delta.Name
		if test.delta.Type == HISTOGRAM || test.delta.Type == SUMMARY {
			metricName += "_sum"
		}
		result, err := readMetric(metrics, metricName, test.delta.Labels)
		if assert.Nil(t, err) {
			test.firstValue = result
			// check if the delta value was added to the metric
			expected := test.f(initialValue, test.delta.Value)
			assert.Equal(t, expected, result)
		}
	}

	// apply all deltas again
	for i := 0; i < len(tests); i++ {
		test := &tests[i]
		err := test.delta.Apply()
		assert.Nil(t, err)
	}

	// get metrics and compare results
	metrics = fetchMetrics(t)
	for i := 0; i < len(tests); i++ {
		test := &tests[i]
		// read result and compare
		metricName := test.delta.Name
		if test.delta.Type == HISTOGRAM || test.delta.Type == SUMMARY {
			metricName += "_sum"
		}
		result, err := readMetric(metrics, metricName, test.delta.Labels)
		if assert.Nil(t, err) {
			// check if the delta value was added to the metric
			expected := test.f(test.firstValue, test.delta.Value)
			assert.Equal(t, expected, result)
		}
	}

}
