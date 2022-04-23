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

func GetServerPrometheusParams() PrometheusParams {
	pp := PrometheusParams{
		Counters:   make(map[string]prometheus.Counter),
		Gauges:     make(map[string]prometheus.Gauge),
		Histograms: make(map[string]prometheus.Histogram),
		Summar:     make(map[string]prometheus.Summary),
	}
	// counters
	pp.Counters["hero_selling_counter"] = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hero_selling_counter",
		Help: "Save the number of heros are selling, currently",
	})
	pp.Counters["bet_system_invite_operation_counter"] = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bet_system_invite_operation_counter",
		Help: "Total number of invitations sent for arena bet match",
	})
	pp.Counters["bet_system_accept_operation_counter"] = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bet_system_accept_operation_counter",
		Help: "Total number of accepted invitations",
	})
	pp.Counters["bet_system_start_game_operation_counter"] = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bet_system_start_game_operation_counter",
		Help: "Total number of start game requests",
	})
	pp.Counters["bet_system_match_counter"] = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bet_system_match_counter",
		Help: "Total number of matches started in bet system",
	})
	pp.Counters["bet_system_match_finished"] = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bet_system_match_finished",
		Help: "Total number of matches finished",
	})
	// gauges
	pp.Gauges["bet_system_invite_operation_in_progress"] = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bet_system_invite_operation_in_progress",
		Help: "Number of invitations sent and has not accepted yet",
	})
	pp.Gauges["bet_system_accept_operation_in_progress"] = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bet_system_accept_operation_in_progress",
		Help: "Number of invitations accepted and has not played yet",
	})
	pp.Gauges["bet_system_match_in_progress"] = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bet_system_match_in_progress",
		Help: "Number of matches are in progress currently",
	})
	// histograms
	pp.Histograms["bet_system_invitation_request_response_duration"] = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "bet_system_invitation_request_response_duration",
		Help: "Response delay in invitation",
	})
	return pp
}
