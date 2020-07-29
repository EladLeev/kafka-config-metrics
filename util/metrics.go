package util

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Define all the metrics as Gauge
// The "metrics" are configuration parameters so it's better to represent them as Gauge metrics.
// TODO: Map it better
var (
	minInsync = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "min_insync_replicas",
			Help:      "Specifies the minimum number of replicas that must acknowledge a write for the write to be considered successful.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	segmentJitterMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "segment_jitter_ms",
			Help:      "The maximum random jitter subtracted from the scheduled segment roll time to avoid thundering herds of segment rolling",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	flushMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "flush_ms",
			Help:      "A time interval at which we will force an fsync of data written to the log.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	logRetentionMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "log_retention_ms",
			Help:      "The maximum time we will retain a log before we will discard old log segments to free up space if we are using the delete retention policy.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	segmentBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "segment_bytes",
			Help:      "his configuration controls the segment file size for the log.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	flushMessages = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "flush_messages",
			Help:      "An interval at which we will force an fsync of data written to the log.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	fileDeleteDelayMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "file_delete_delay_ms",
			Help:      "The time to wait before deleting a file from the filesystem.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	maxCompactionLagMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "max_compaction_lag_ms",
			Help:      "The maximum time a message will remain ineligible for compaction in the log.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	maxMessageBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "max_message_bytes",
			Help:      "The largest record batch size allowed by Kafka.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	minCompactionLagMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "min_compaction_lag_ms",
			Help:      "The minimum time a message will remain uncompacted in the log. Only applicable for logs that are being compacted.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	indexIntervalBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "index_interval_bytes",
			Help:      "How frequently Kafka adds an index entry to its offset index.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	retentionBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "retention_bytes",
			Help:      "The maximum size a partition (which consists of log segments) can grow to before we will discard old log segments to free up space if we are using the delete retention policy",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	deleteRetentionMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "delete_retention_ms",
			Help:      "The amount of time to retain delete tombstone markers for log compacted topics.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	segmentMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "segment_ms",
			Help:      "This configuration controls the period of time after which Kafka will force the log to roll even if the segment file isn't full to ensure that retention can delete or compact old data.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	messageTimestampDifferenceMaxMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "message_timestamp_difference_max_ms",
			Help:      "The maximum difference allowed between the timestamp when a broker receives a message and the timestamp specified in the message.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
	segmentIndexBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "segment_index_bytes",
			Help:      "This configuration controls the size of the index that maps offsets to file positions. We preallocate this index file and shrink it only after log rolls. You generally should not need to change this setting.",
			Namespace: "topic",
		},
		[]string{"topic"},
	)
)
