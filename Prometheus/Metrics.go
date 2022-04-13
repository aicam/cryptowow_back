package Prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusParams struct {
	Counters   map[string]prometheus.Counter
	Gauges     map[string]prometheus.Gauge
	Histograms map[string]prometheus.Histogram
	Summar     map[string]prometheus.Summary
}

func GetIndexServerPrometheusParams() PrometheusParams {
	pp := PrometheusParams{
		Counters:   make(map[string]prometheus.Counter),
		Gauges:     make(map[string]prometheus.Gauge),
		Histograms: make(map[string]prometheus.Histogram),
		Summar:     make(map[string]prometheus.Summary),
	}
	pp.Counters["hero_selling_counter"] = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hero_selling_counter",
		Help: "Save the number of heros are selling, currently",
	})
	return pp
}
