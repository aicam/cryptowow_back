package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusParams struct {
	Counters map[string]prometheus.Counter
	Gauges   map[string]prometheus.Gauge
}

func main() {

	//myfile, e := os.Create("log.txt")
	//if e != nil {
	//	log.Fatal(e)
	//}
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetOutput(myfile)
	//logrus.WithFields(logrus.Fields{"test": 2, "test_s": "dsada"}).Info("detail test")
	//pp := PrometheusParams{}
	//pp.Counters = make(map[string]prometheus.Counter)
	//pp.Gauges = make(map[string]prometheus.Gauge)
	//pp.Counters["opc"] = promauto.NewCounter(prometheus.CounterOpts{
	//	Name: "opc",
	//	Help: "Test Counter",
	//})
	//pp.Gauges["opg"] = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "opg",
	//	Help: "Test Gauge",
	//})
	//pp.Counters["opc"].Inc()
	//pp.Gauges["opg"].Set(22.0)
	//http.Handle("/metrics", promhttp.Handler())
	//http.ListenAndServe(":2112", nil)
}
