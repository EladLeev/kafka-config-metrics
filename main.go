package main

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"github.com/EladLeev/kafka-config-metrics/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func filterTopic(filterRegex string) *regexp.Regexp {
	re, err := regexp.Compile(filterRegex)
	if err != nil {
		fmt.Print(err)
	}
	return re
}

func pullConfigs(cfg util.TomlConfig, clusterName string) {
	clusterBrokers := cfg.Clusters[clusterName].Brokers
	clusterAdmin := util.OpenConnection(cfg.Kafka.MinKafkaVersion, clusterBrokers, cfg.Kafka.AdminTimeout)
	log.Debugf("KafkaVersion: %v | BrokerList: %v", cfg.Kafka.MinKafkaVersion, clusterBrokers)
	defer clusterAdmin.Close()

	// Pull topic list and topic configs
	var wg sync.WaitGroup
	var re *regexp.Regexp
	topicList := util.ListTopics(clusterAdmin)
	// If `topicfilter` config is not empty, compile Regex
	if cfg.Clusters[clusterName].TopicFilter != "" {
		re = filterTopic(cfg.Clusters[clusterName].TopicFilter)
	}

	for topic := range topicList {
		if (re != nil) && re.MatchString(topic) {
			continue
		}
		wg.Add(1)
		go util.DescribeTopicConfig(topic, clusterAdmin, &wg)
	}
	wg.Wait()
}

func main() {
	// Init config and logger
	cfg := util.InitConfig()
	util.InitLog(cfg)
	log.Infof(" \\ʕ ◔ ϖ ◔ ʔ/ ## Kafka-Config-Metrics-Exporter ## \\ʕ ◔ ϖ ◔ ʔ/\n")
	log.Debugf("cfg: %+v", cfg)

	// Init Prometheus metrics
	util.InitProm()

	// Pull from each cluster
	for cluster := range cfg.Clusters {
		pullConfigs(cfg, cluster)
	}
	log.Info("done")

	// Exspose Prometheus endpoint
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
