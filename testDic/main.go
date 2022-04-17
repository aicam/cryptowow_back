package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

type PrometheusParams struct {
	Counters map[string]prometheus.Counter
	Gauges   map[string]prometheus.Gauge
}

func fuckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
	//DB := database.DbSqlMigration("aicam:021021ali@tcp(127.0.0.1:3306)/cryptowow?charset=utf8mb4&parseTime=True")
	//s := ArenaService.Service{
	//	DB:      DB,
	//	Redis:   rdb,
	//	Context: context.Background(),
	//	PP:      Prometheus.PrometheusParams{},
	//}
	//

	myfile, e := os.Create("log.txt")
	if e != nil {
		log.Fatal(e)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(myfile)
	logrus.WithFields(logrus.Fields{"test": 2, "test_s": "dsada"}).Error("detail test")
	fmt.Println("asdas")
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
