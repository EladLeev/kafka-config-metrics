package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/IBM/sarama"

	"github.com/EladLeev/kafka-config-metrics/pkg/config"
)

// Client manages Kafka cluster connections and operations
type Client struct {
	version      sarama.KafkaVersion
	adminTimeout time.Duration
}

// NewClient creates a new Kafka client
func NewClient(kafkaVersion string, adminTimeout time.Duration) (*Client, error) {
	version, err := sarama.ParseKafkaVersion(kafkaVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid Kafka version: %w", err)
	}

	return &Client{
		version:      version,
		adminTimeout: adminTimeout,
	}, nil
}

func (c *Client) createTLSConfig(tlsCfg config.TLSConfig) (*tls.Config, error) {
	if !tlsCfg.Enabled {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Load CA cert if provided
	if tlsCfg.CACert != "" {
		caCert, err := os.ReadFile(tlsCfg.CACert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	// Load client cert/key if provided
	if tlsCfg.ClientCert != "" && tlsCfg.ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(tlsCfg.ClientCert, tlsCfg.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate/key: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

// Connect establishes a connection to a Kafka cluster
func (c *Client) Connect(clusterConfig config.ClusterConfig) (sarama.ClusterAdmin, error) {
	connectConfig := sarama.NewConfig()
	connectConfig.Version = c.version
	connectConfig.Admin.Timeout = c.adminTimeout

	// Configure TLS if enabled
	if clusterConfig.TLS.Enabled {
		tlsConfig, err := c.createTLSConfig(clusterConfig.TLS)
		if err != nil {
			return nil, fmt.Errorf("failed to configure TLS: %w", err)
		}
		connectConfig.Net.TLS.Enable = true
		connectConfig.Net.TLS.Config = tlsConfig
	}

	// Configure SASL if enabled
	if creds := clusterConfig.GetCredentials(); creds != nil && creds.Enabled {
		connectConfig.Net.SASL.Enable = true
		connectConfig.Net.SASL.User = creds.Username
		connectConfig.Net.SASL.Password = creds.Password

		switch creds.Mechanism {
		case "PLAIN":
			connectConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		case "SCRAM-SHA-256":
			connectConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		case "SCRAM-SHA-512":
			connectConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		default:
			return nil, fmt.Errorf("unsupported SASL mechanism: %s", creds.Mechanism)
		}
	}

	admin, err := sarama.NewClusterAdmin(clusterConfig.Brokers, connectConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cluster: %w", err)
	}
	return admin, nil
}

// GetFilteredTopics returns topics matching the filter
func (c *Client) GetFilteredTopics(admin sarama.ClusterAdmin, filterPattern string) ([]string, error) {
	topics, err := admin.ListTopics()
	if err != nil {
		return nil, fmt.Errorf("failed to list topics: %w", err)
	}

	if filterPattern == "" {
		result := make([]string, 0, len(topics))
		for topic := range topics {
			result = append(result, topic)
		}
		return result, nil
	}

	re, err := regexp.Compile(filterPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid topic filter pattern: %w", err)
	}

	var filtered []string
	for topic := range topics {
		if !re.MatchString(topic) {
			filtered = append(filtered, topic)
		}
	}
	return filtered, nil
}

// GetTopicConfig retrieves configuration for a specific topic
func (c *Client) GetTopicConfig(admin sarama.ClusterAdmin, topic string) ([]sarama.ConfigEntry, error) {
	resource := sarama.ConfigResource{
		Type: sarama.TopicResource,
		Name: topic,
	}

	config, err := admin.DescribeConfig(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to describe config for topic %s: %w", topic, err)
	}
	return config, nil
}
