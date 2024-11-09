package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricDefinition represents a single metric configuration
type MetricDefinition struct {
	Name      string
	Help      string
	ConfigKey string
}

// Collector manages all Kafka configuration metrics
type Collector struct {
	metrics map[string]*prometheus.GaugeVec
	mu      sync.RWMutex
}

// NewCollector creates and initializes a new Collector
func NewCollector() *Collector {
	c := &Collector{
		metrics: make(map[string]*prometheus.GaugeVec),
	}

	definitions := []MetricDefinition{
		{Name: "min_insync_replicas", Help: "Minimum number of replicas that must acknowledge a write", ConfigKey: "min.insync.replicas"},
		{Name: "segment_jitter_ms", Help: "Maximum random jitter for segment rolling", ConfigKey: "segment.jitter.ms"},
		{Name: "flush_ms", Help: "Time interval for forcing fsync of data", ConfigKey: "flush.ms"},
		{Name: "log_retention_ms", Help: "Maximum time to retain a log", ConfigKey: "retention.ms"},
		{Name: "segment_bytes", Help: "Controls the segment file size for the log", ConfigKey: "segment.bytes"},
		{Name: "flush_messages", Help: "Number of messages after which to force an fsync", ConfigKey: "flush.messages"},
		{Name: "file_delete_delay_ms", Help: "Time to wait before deleting a file", ConfigKey: "file.delete.delay.ms"},
		{Name: "max_compaction_lag_ms", Help: "Maximum time a message remains ineligible for compaction", ConfigKey: "max.compaction.lag.ms"},
		{Name: "max_message_bytes", Help: "Maximum size of a message batch", ConfigKey: "max.message.bytes"},
		{Name: "min_compaction_lag_ms", Help: "Minimum time a message remains uncompacted", ConfigKey: "min.compaction.lag.ms"},
		{Name: "index_interval_bytes", Help: "Index entry addition frequency", ConfigKey: "index.interval.bytes"},
		{Name: "retention_bytes", Help: "Maximum size of a partition", ConfigKey: "retention.bytes"},
		{Name: "delete_retention_ms", Help: "Time to retain delete tombstone markers", ConfigKey: "delete.retention.ms"},
		{Name: "segment_ms", Help: "Time after which log will be forced to roll", ConfigKey: "segment.ms"},
		{Name: "message_timestamp_difference_max_ms", Help: "Maximum allowed difference between message timestamp and broker time", ConfigKey: "message.timestamp.difference.max.ms"},
		{Name: "segment_index_bytes", Help: "Size of the index mapping offsets to file positions", ConfigKey: "segment.index.bytes"},
	}

	// Initialize all metrics
	for _, def := range definitions {
		c.metrics[def.Name] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      def.Name,
				Help:      def.Help,
				Namespace: "kafka_topic_config",
			},
			[]string{"topic", "cluster"},
		)
	}

	return c
}

// SetMetric safely sets a metric value
func (c *Collector) SetMetric(name, topic, cluster string, value float64) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if metric, exists := c.metrics[name]; exists {
		metric.With(prometheus.Labels{
			"topic":   topic,
			"cluster": cluster,
		}).Set(value)
	}
}

// Reset safely resets all metrics
func (c *Collector) Reset() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, metric := range c.metrics {
		metric.Reset()
	}
}

// Collect implements prometheus.Collector
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, metric := range c.metrics {
		metric.Collect(ch)
	}
}

// Describe implements prometheus.Collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, metric := range c.metrics {
		metric.Describe(ch)
	}
}

// GetMetricByKey returns the metric name for a given Kafka config key
func (c *Collector) GetMetricByKey(configKey string) string {
	definitions := []MetricDefinition{
		{Name: "min_insync_replicas", ConfigKey: "min.insync.replicas"},
		{Name: "segment_jitter_ms", ConfigKey: "segment.jitter.ms"},
		{Name: "flush_ms", ConfigKey: "flush.ms"},
		{Name: "log_retention_ms", ConfigKey: "retention.ms"},
		{Name: "segment_bytes", ConfigKey: "segment.bytes"},
		{Name: "flush_messages", ConfigKey: "flush.messages"},
		{Name: "file_delete_delay_ms", ConfigKey: "file.delete.delay.ms"},
		{Name: "max_compaction_lag_ms", ConfigKey: "max.compaction.lag.ms"},
		{Name: "max_message_bytes", ConfigKey: "max.message.bytes"},
		{Name: "min_compaction_lag_ms", ConfigKey: "min.compaction.lag.ms"},
		{Name: "index_interval_bytes", ConfigKey: "index.interval.bytes"},
		{Name: "retention_bytes", ConfigKey: "retention.bytes"},
		{Name: "delete_retention_ms", ConfigKey: "delete.retention.ms"},
		{Name: "segment_ms", ConfigKey: "segment.ms"},
		{Name: "message_timestamp_difference_max_ms", ConfigKey: "message.timestamp.difference.max.ms"},
		{Name: "segment_index_bytes", ConfigKey: "segment.index.bytes"},
	}

	for _, def := range definitions {
		if def.ConfigKey == configKey {
			return def.Name
		}
	}
	return ""
}
