package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EladLeev/kafka-config-metrics/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Init config and logger
	util.InitLog(util.Configuration)
	log.Infof(" \\ʕ ◔ ϖ ◔ ʔ/ ## Kafka-Config-Metrics-Exporter ## \\ʕ ◔ ϖ ◔ ʔ/\n")
	log.Debugf("cfg: %+v", util.Configuration)

	// Init Prometheus metrics
	util.InitProm()

	// Exspose Prometheus endpoints
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	server := &http.Server{
		Addr:              util.Configuration.Global.Port,
		ReadHeaderTimeout: time.Duration(util.Configuration.Global.Timeout) * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
