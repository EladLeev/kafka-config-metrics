package config

import "time"

// Config represents the application configuration
type Config struct {
	Global   GlobalConfig             `yaml:"global"`
	Log      LogConfig                `yaml:"log"`
	Kafka    KafkaConfig              `yaml:"kafka"`
	Clusters map[string]ClusterConfig `yaml:"clusters"`
}

type GlobalConfig struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type LogConfig struct {
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
}

type KafkaConfig struct {
	DefaultRefreshRate time.Duration `yaml:"defaultRefreshRate"`
	MinKafkaVersion    string        `yaml:"minKafkaVersion"`
	AdminTimeout       time.Duration `yaml:"adminTimeout"`
}

type TLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	CACert     string `yaml:"caCert"`
	ClientCert string `yaml:"clientCert"`
	ClientKey  string `yaml:"clientKey"`
}

type SASLConfig struct {
	Enabled     bool   `yaml:"enabled"`
	UsernameEnv string `yaml:"usernameEnv"`
	PasswordEnv string `yaml:"passwordEnv"`
	Mechanism   string `yaml:"mechanism"`
}

type RuntimeSASLConfig struct {
	Enabled   bool
	Username  string
	Password  string
	Mechanism string
}

type ClusterConfig struct {
	Brokers     []string   `yaml:"brokers"`
	TopicFilter string     `yaml:"topicFilter"`
	TLS         TLSConfig  `yaml:"tls"`
	SASL        SASLConfig `yaml:"sasl"`
	credentials *SASLCredentials
}

type SASLCredentials struct {
	Enabled   bool
	Username  string
	Password  string
	Mechanism string
}
