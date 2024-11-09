package main

import (
	"github.com/sirupsen/logrus"

	"github.com/EladLeev/kafka-config-metrics/pkg/config"
)

func initLogger(cfg config.LogConfig) {
	if cfg.Level == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
	}

	if cfg.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}
