package main

import (
	"net/http"

	"github.com/EladLeev/kafka-config-metrics/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Init config and logger
	cfg := util.InitConfig()
	util.InitLog(cfg)
	log.Infof(" \\ʕ ◔ ϖ ◔ ʔ/ ## Kafka-Config-Metrics-Exporter ## \\ʕ ◔ ϖ ◔ ʔ/\n")
	log.Debugf("cfg: %+v", cfg)

	// Init Prometheus metrics
	util.InitProm()

	// Exspose Prometheus endpoint
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
