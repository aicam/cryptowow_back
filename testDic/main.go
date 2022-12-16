package main

import (
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server/WalletService"
	"github.com/prometheus/client_golang/prometheus"
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
	//cnt := context.Background()
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
	DBStruct := database.DbSqlMigration(os.Getenv("MAINMYSQLCONNECTION"))
	WalletService.GetArenaBetTotalDebt(DBStruct, "T6")
	//for i := 0; i < 50; i++ {
	//	rdb.Set(cnt, strconv.Itoa(i + 20), "a", 100 * time.Second)
	//}
	//n := time.Now()
	//time.Sleep(1001 * time.Millisecond)
	//n2 := time.Now()
	//log.Println(n2.Sub(n).Milliseconds())
	//res := rdb.Get(cnt, "v").Val()
	//log.Println(res)
	//DB := database.DbSqlMigration("aicam:021021ali@tcp(127.0.0.1:3306)/cryptowow?charset=utf8mb4&parseTime=True")
	//s := ArenaService.Service{
	//	DB:      DB,
	//	Rdb:   rdb,
	//	Context: context.Background(),
	//	PP:      Prometheus.PrometheusParams{},
	//}
	//

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
