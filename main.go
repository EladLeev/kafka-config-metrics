package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/EladLeev/kafka-config-metrics/pkg/config"
	"github.com/EladLeev/kafka-config-metrics/pkg/kafka"
	"github.com/EladLeev/kafka-config-metrics/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "/opt/kcm/kcm.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	initLogger(cfg.Log)

	// Initialize Kafka client and collector
	kafkaClient, err := kafka.NewClient(cfg.Kafka.MinKafkaVersion, cfg.Kafka.AdminTimeout)
	if err != nil {
		logrus.Fatalf("Failed to create Kafka client: %v", err)
	}

	collector := metrics.NewCollector()
	prometheus.MustRegister(collector)

	// Start metrics collection for each cluster
	for clusterName, clusterConfig := range cfg.Clusters {
		go func(name string, config config.ClusterConfig) {
			for {
				if err := collectMetrics(kafkaClient, collector, name, config); err != nil {
					logrus.Errorf("Failed to collect metrics for cluster %s: %v", name, err)
				}
				time.Sleep(cfg.Kafka.DefaultRefreshRate)
			}
		}(clusterName, clusterConfig)
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/-/healthy", healthHandler)
	mux.HandleFunc("/-/ready", readyHandler)

	server := &http.Server{
		Addr:              cfg.Global.Port,
		Handler:           mux,
		ReadHeaderTimeout: cfg.Global.Timeout,
	}

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logrus.Infof(" \\ʕ ◔ ϖ ◔ ʔ/ ## Kafka-Config-Metrics Exporter ## \\ʕ ◔ ϖ ◔ ʔ/")
		logrus.Infof("Starting server on %s", cfg.Global.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	logrus.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.Errorf("Server shutdown failed: %v", err)
	}
}

func collectMetrics(client *kafka.Client, collector *metrics.Collector, clusterName string, config config.ClusterConfig) error {
	admin, err := client.Connect(config)
	if err != nil {
		return err
	}
	defer admin.Close()

	topics, err := client.GetFilteredTopics(admin, config.TopicFilter)
	if err != nil {
		return err
	}

	collector.Reset()

	for _, topic := range topics {
		topicConfig, err := client.GetTopicConfig(admin, topic)
		if err != nil {
			logrus.Errorf("Failed to get config for topic %s: %v", topic, err)
			continue
		}

		for _, entry := range topicConfig {
			if metricName := collector.GetMetricByKey(entry.Name); metricName != "" {
				if value, err := strconv.ParseFloat(entry.Value, 64); err == nil {
					collector.SetMetric(metricName, topic, clusterName, value)
				}
			}
		}
	}

	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
