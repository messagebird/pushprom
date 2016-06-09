package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricType is a string that uniquelly identifies a Prometheus metric
type MetricType string

const (
	// COUNTER is a Prometheus metric type that only goes up
	COUNTER MetricType = "counter"
	// GAUGE is a Prometheus metric type that goes up or down
	GAUGE MetricType = "gauge"
	// HISTOGRAM is a Prometheus metric type
	HISTOGRAM MetricType = "histogram"
	// SUMMARY is a Prometheus metric type
	SUMMARY MetricType = "summary"
)

// Delta defines the pushprom message format. It represent a change on a Prometheus metric. It implements a simplistic rpc.
type Delta struct {
	Type   MetricType        `json:"type"`
	Name   string            `json:"name"`
	Help   string            `json:"help"`
	Method string            `json:"method"`
	Value  float64           `json:"value"`
	Labels prometheus.Labels `json:"labels"`
}

// NewDelta creates a new Delta from the json contents of the io.Reader
func NewDelta(reader io.Reader) (*Delta, error) {
	dec := json.NewDecoder(reader)
	delta := Delta{}
	if err := dec.Decode(&delta); err != nil {
		return nil, err
	}
	return &delta, nil
}

// LabelNames returns the names(keys) of all labels
func (delta Delta) LabelNames() []string {
	names := []string{}
	for k := range delta.Labels {
		names = append(names, k)
	}
	return names
}

func (delta Delta) applyOnCounter() error {
	var metric prometheus.Counter
	opts := prometheus.CounterOpts{
		Name: delta.Name,
		Help: delta.Help,
	}
	if len(delta.Labels) > 0 {
		vec := prometheus.NewCounterVec(opts, delta.LabelNames())

		registered, err := prometheus.RegisterOrGet(vec)
		if err != nil {
			return err
		}
		if registered != nil {
			vec = registered.(*prometheus.CounterVec)
		}

		metric = vec.With(delta.Labels)
	} else {
		metric = prometheus.NewCounter(opts)

		registered, err := prometheus.RegisterOrGet(metric)
		if err != nil {
			return err
		}
		if registered != nil {
			metric = registered.(prometheus.Counter)
		}
	}

	switch delta.Method {
	case "set":
		metric.Set(delta.Value)
	case "inc":
		metric.Inc()
	case "add":
		metric.Add(delta.Value)
	default:
		return fmt.Errorf("method %s is not implemented yet", delta.Method)
	}

	return nil
}

func (delta Delta) applyOnGauge() error {
	var metric prometheus.Gauge
	opts := prometheus.GaugeOpts{
		Name: delta.Name,
		Help: delta.Help,
	}
	if len(delta.Labels) > 0 {
		vec := prometheus.NewGaugeVec(opts, delta.LabelNames())

		registered, err := prometheus.RegisterOrGet(vec)
		if err != nil {
			return err
		}
		if registered != nil {
			vec = registered.(*prometheus.GaugeVec)
		}

		metric = vec.With(delta.Labels)
	} else {
		metric = prometheus.NewGauge(opts)

		registered, err := prometheus.RegisterOrGet(metric)
		if err != nil {
			return err
		}
		if registered != nil {
			metric = registered.(prometheus.Gauge)
		}
	}

	switch delta.Method {
	case "set":
		metric.Set(delta.Value)
	case "inc":
		metric.Inc()
	case "dec":
		metric.Dec()
	case "add":
		metric.Add(delta.Value)
	case "sub":
		metric.Sub(delta.Value)
	default:
		return fmt.Errorf("method %s is not implemented yet", delta.Method)
	}

	return nil
}

func (delta Delta) applyOnHistogram() error {
	var metric prometheus.Histogram
	opts := prometheus.HistogramOpts{
		Name: delta.Name,
		Help: delta.Help,
	}
	if len(delta.Labels) > 0 {
		vec := prometheus.NewHistogramVec(opts, delta.LabelNames())

		registered, err := prometheus.RegisterOrGet(vec)
		if err != nil {
			return err
		}
		if registered != nil {
			vec = registered.(*prometheus.HistogramVec)
		}

		metric = vec.With(delta.Labels)
	} else {
		metric = prometheus.NewHistogram(opts)

		registered, err := prometheus.RegisterOrGet(metric)
		if err != nil {
			return err
		}
		if registered != nil {
			metric = registered.(prometheus.Histogram)
		}
	}

	if delta.Method == "observe" {
		metric.Observe(delta.Value)
	} else {
		return fmt.Errorf("method %s is not implemented yet", delta.Method)
	}

	return nil
}

func (delta Delta) applyOnSummary() error {
	var metric prometheus.Summary
	opts := prometheus.SummaryOpts{
		Name: delta.Name,
		Help: delta.Help,
	}
	if len(delta.Labels) > 0 {
		vec := prometheus.NewSummaryVec(opts, delta.LabelNames())

		registered, err := prometheus.RegisterOrGet(vec)
		if err != nil {
			return err
		}
		if registered != nil {
			vec = registered.(*prometheus.SummaryVec)
		}

		metric = vec.With(delta.Labels)
	} else {
		metric = prometheus.NewSummary(opts)

		registered, err := prometheus.RegisterOrGet(metric)
		if err != nil {
			return err
		}
		if registered != nil {
			metric = registered.(prometheus.Summary)
		}
	}

	if delta.Method == "observe" {
		metric.Observe(delta.Value)
	} else {
		return fmt.Errorf("method %s is not implemented yet", delta.Method)
	}

	return nil
}

// Apply applies the delta (calls delta.method(delta.value)) on the correspondent (registered) Prometheus metric
func (delta Delta) Apply() error {
	switch delta.Type {
	case COUNTER:
		return delta.applyOnCounter()
	case GAUGE:
		return delta.applyOnGauge()
	case HISTOGRAM:
		return delta.applyOnHistogram()
	case SUMMARY:
		return delta.applyOnSummary()
	}
	return fmt.Errorf("metric type %s is not implemented yet", delta.Type)
}
