package util

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// InitLog init the logger
func InitLog(cfg TomlConfig) {
	// Check if log level
	if cfg.Log.Level == "debug" {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano, // 2019-04-22T10:27:50.491904208+03:00
		})
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		log.SetLevel(log.InfoLevel)
	}
	// Enable JSON log for easy parsing by Logstash or Splunk
	if cfg.Log.Format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}
