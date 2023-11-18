package util

import (
	"strconv"
	"sync"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector define all metrics as GaugeVec
// The "metrics" are configuration parameters so it's better to represent them as Gauge metrics.
type Collector struct {
	sync.Mutex
	minInsync                       *prometheus.GaugeVec
	segmentJitterMs                 *prometheus.GaugeVec
	flushMs                         *prometheus.GaugeVec
	logRetentionMs                  *prometheus.GaugeVec
	segmentBytes                    *prometheus.GaugeVec
	flushMessages                   *prometheus.GaugeVec
	fileDeleteDelayMs               *prometheus.GaugeVec
	maxCompactionLagMs              *prometheus.GaugeVec
	maxMessageBytes                 *prometheus.GaugeVec
	minCompactionLagMs              *prometheus.GaugeVec
	indexIntervalBytes              *prometheus.GaugeVec
	retentionBytes                  *prometheus.GaugeVec
	deleteRetentionMs               *prometheus.GaugeVec
	segmentMs                       *prometheus.GaugeVec
	messageTimestampDifferenceMaxMs *prometheus.GaugeVec
	segmentIndexBytes               *prometheus.GaugeVec
}

// NewCollector for all metrics
func NewCollector() *Collector {
	return &Collector{
		minInsync: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "min_insync_replicas",
				Help:      "Specifies the minimum number of replicas that must acknowledge a write for the write to be considered successful.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		segmentJitterMs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "segment_jitter_ms",
				Help:      "The maximum random jitter subtracted from the scheduled segment roll time to avoid thundering herds of segment rolling",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		flushMs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "flush_ms",
				Help:      "A time interval at which we will force an fsync of data written to the log.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		logRetentionMs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "log_retention_ms",
				Help:      "The maximum time we will retain a log before we will discard old log segments to free up space if we are using the delete retention policy.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		segmentBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "segment_bytes",
				Help:      "his configuration controls the segment file size for the log.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		flushMessages: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "flush_messages",
				Help:      "An interval at which we will force an fsync of data written to the log.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		fileDeleteDelayMs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "file_delete_delay_ms",
				Help:      "The time to wait before deleting a file from the filesystem.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		maxCompactionLagMs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "max_compaction_lag_ms",
				Help:      "The maximum time a message will remain ineligible for compaction in the log.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		maxMessageBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "max_message_bytes",
				Help:      "The largest record batch size allowed by Kafka.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		minCompactionLagMs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "min_compaction_lag_ms",
				Help:      "The minimum time a message will remain uncompacted in the log. Only applicable for logs that are being compacted.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		indexIntervalBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "index_interval_bytes",
				Help:      "How frequently Kafka adds an index entry to its offset index.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		retentionBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "retention_bytes",
				Help:      "The maximum size a partition (which consists of log segments) can grow to before we will discard old log segments to free up space if we are using the delete retention policy",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		deleteRetentionMs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "delete_retention_ms",
				Help:      "The amount of time to retain delete tombstone markers for log compacted topics.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		segmentMs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "segment_ms",
				Help:      "This configuration controls the period of time after which Kafka will force the log to roll even if the segment file isn't full to ensure that retention can delete or compact old data.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		messageTimestampDifferenceMaxMs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "message_timestamp_difference_max_ms",
				Help:      "The maximum difference allowed between the timestamp when a broker receives a message and the timestamp specified in the message.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
		segmentIndexBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "segment_index_bytes",
				Help:      "This configuration controls the size of the index that maps offsets to file positions. We preallocate this index file and shrink it only after log rolls. You generally should not need to change this setting.",
				Namespace: "topic",
			},
			[]string{"topic"},
		),
	}
}

func resetMetrics(c *Collector) {
	c.minInsync.Reset()
	c.segmentJitterMs.Reset()
	c.flushMs.Reset()
	c.logRetentionMs.Reset()
	c.segmentBytes.Reset()
	c.flushMessages.Reset()
	c.fileDeleteDelayMs.Reset()
	c.maxCompactionLagMs.Reset()
	c.maxMessageBytes.Reset()
	c.minCompactionLagMs.Reset()
	c.indexIntervalBytes.Reset()
	c.retentionBytes.Reset()
	c.deleteRetentionMs.Reset()
	c.segmentMs.Reset()
	c.messageTimestampDifferenceMaxMs.Reset()
	c.segmentIndexBytes.Reset()
}

// Collect implements Collector
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	// Prevent execution on 2 Collectors at once
	c.Lock()
	defer c.Unlock()

	// Reset metrics
	resetMetrics(c)

	// Pull from each cluster on cfg.Clusters
	/*
		"Metrics should only be pulled from the application when Prometheus scrapes them,
				exporters should not perform scrapes based on their own timers.
					That is, all scrapes should be synchronous."
		https://prometheus.io/docs/instrumenting/writing_exporters/#scheduling
	*/
	for cluster := range Configuration.Clusters {
		PullConfigs(Configuration, cluster, c)
	}

	// Collect
	c.minInsync.Collect(ch)
	c.segmentJitterMs.Collect(ch)
	c.flushMs.Collect(ch)
	c.logRetentionMs.Collect(ch)
	c.segmentBytes.Collect(ch)
	c.flushMessages.Collect(ch)
	c.fileDeleteDelayMs.Collect(ch)
	c.maxCompactionLagMs.Collect(ch)
	c.maxMessageBytes.Collect(ch)
	c.minCompactionLagMs.Collect(ch)
	c.indexIntervalBytes.Collect(ch)
	c.retentionBytes.Collect(ch)
	c.deleteRetentionMs.Collect(ch)
	c.segmentMs.Collect(ch)
	c.messageTimestampDifferenceMaxMs.Collect(ch)
	c.segmentIndexBytes.Collect(ch)
}

// Describe implements Collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.minInsync.Describe(ch)
	c.segmentJitterMs.Describe(ch)
	c.flushMs.Describe(ch)
	c.logRetentionMs.Describe(ch)
	c.segmentBytes.Describe(ch)
	c.flushMessages.Describe(ch)
	c.fileDeleteDelayMs.Describe(ch)
	c.maxCompactionLagMs.Describe(ch)
	c.maxMessageBytes.Describe(ch)
	c.minCompactionLagMs.Describe(ch)
	c.indexIntervalBytes.Describe(ch)
	c.retentionBytes.Describe(ch)
	c.deleteRetentionMs.Describe(ch)
	c.segmentMs.Describe(ch)
	c.messageTimestampDifferenceMaxMs.Describe(ch)
	c.segmentIndexBytes.Describe(ch)
}

// InitProm register the Prometheus metrics
func InitProm() {
	prometheus.MustRegister(NewCollector())
}

// RegisterMetrics populate the metrics with a topic label
func RegisterMetrics(topicConfig []sarama.ConfigEntry, topicName string, c *Collector) {
	for _, v := range topicConfig {
		switch v.Name {
		case "min.insync.replicas":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.minInsync.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "segment.jitter.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.segmentJitterMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "retention.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.logRetentionMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "segment.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.segmentBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "flush.messages":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.flushMessages.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "file.delete.delay.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.fileDeleteDelayMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "max.compaction.lag.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.maxCompactionLagMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "max.message.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.maxMessageBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "min.compaction.lag.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.minCompactionLagMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "index.interval.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.indexIntervalBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "retention.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.retentionBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "delete.retention.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.deleteRetentionMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "segment.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.segmentMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "message.timestamp.difference.max.ms":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.messageTimestampDifferenceMaxMs.With(prometheus.Labels{"topic": topicName}).Set(m)
		case "segment.index.bytes":
			m, _ := strconv.ParseFloat(v.Value, 64)
			c.segmentIndexBytes.With(prometheus.Labels{"topic": topicName}).Set(m)
		}
	}

}
