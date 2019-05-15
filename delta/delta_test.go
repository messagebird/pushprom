package delta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/messagebird/pushprom/metrics"
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

	ms, err := metrics.Fetch()
	assert.Nil(t, err)
	result, err := metrics.Read(ms, delta.Name, delta.Labels)
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

	ms, err := metrics.Fetch()
	assert.Nil(t, err)
	result, err := metrics.Read(ms, delta.Name, delta.Labels)
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
	ms, err := metrics.Fetch()
	assert.Nil(t, err)

	for i := 0; i < len(tests); i++ {
		test := &tests[i]
		// read result and compare
		metricName := test.delta.Name
		if test.delta.Type == HISTOGRAM || test.delta.Type == SUMMARY {
			metricName += "_sum"
		}
		result, err := metrics.Read(ms, metricName, test.delta.Labels)
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
	ms, err = metrics.Fetch()
	assert.Nil(t, err)
	for i := 0; i < len(tests); i++ {
		test := &tests[i]
		// read result and compare
		metricName := test.delta.Name
		if test.delta.Type == HISTOGRAM || test.delta.Type == SUMMARY {
			metricName += "_sum"
		}
		result, err := metrics.Read(ms, metricName, test.delta.Labels)
		if assert.Nil(t, err) {
			// check if the delta value was added to the metric
			expected := test.f(test.firstValue, test.delta.Value)
			assert.Equal(t, expected, result)
		}
	}

}
