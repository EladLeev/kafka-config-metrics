package util

import (
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

// InitProm register the Prometheus metrics
func InitProm() {
	prometheus.MustRegister(minInsync)
	prometheus.MustRegister(segmentJitterMs)
	prometheus.MustRegister(flushMs)
	prometheus.MustRegister(logRetentionMs)
	prometheus.MustRegister(segmentBytes)
	prometheus.MustRegister(flushMessages)
	prometheus.MustRegister(fileDeleteDelayMs)
	prometheus.MustRegister(maxCompactionLagMs)
	prometheus.MustRegister(maxMessageBytes)
	prometheus.MustRegister(minCompactionLagMs)
	prometheus.MustRegister(indexIntervalBytes)
	prometheus.MustRegister(retentionBytes)
	prometheus.MustRegister(deleteRetentionMs)
	prometheus.MustRegister(segmentMs)
	prometheus.MustRegister(messageTimestampDifferenceMaxMs)
	prometheus.MustRegister(segmentIndexBytes)
}

// RegisterMetrics populate the metrics with a topic label
func RegisterMetrics(topicConfig []sarama.ConfigEntry, topicName string) {
	for _, v := range topicConfig {
		switch v.Name {
		case "min.insync.replicas":
			m, _ := strconv.ParseFloat(v.Value, 64)
			minInsync.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "segment.jitter.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			segmentJitterMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "retention.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			logRetentionMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "segment.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			segmentBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "flush.messages":
			m, _ := strconv.ParseFloat(v.Value, 64)
			flushMessages.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "file.delete.delay.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			fileDeleteDelayMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "max.compaction.lag.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			maxCompactionLagMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "max.message.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			maxMessageBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "min.compaction.lag.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			minCompactionLagMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "index.interval.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			indexIntervalBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "retention.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			retentionBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "delete.retention.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			deleteRetentionMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "segment.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			segmentMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "message.timestamp.difference.max.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			messageTimestampDifferenceMaxMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "segment.index.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			segmentIndexBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		}
	}

}
