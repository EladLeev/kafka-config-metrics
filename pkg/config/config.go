package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

func (c *ClusterConfig) GetCredentials() *SASLCredentials {
	return c.credentials
}

func (c *ClusterConfig) LoadSASLCredentials(clusterName string) error {
	if !c.SASL.Enabled {
		return nil
	}

	username := os.Getenv(c.SASL.UsernameEnv)
	if username == "" {
		return fmt.Errorf("environment variable %s for cluster %s SASL username not set",
			c.SASL.UsernameEnv, clusterName)
	}

	password := os.Getenv(c.SASL.PasswordEnv)
	if password == "" {
		return fmt.Errorf("environment variable %s for cluster %s SASL password not set",
			c.SASL.PasswordEnv, clusterName)
	}

	c.credentials = &SASLCredentials{
		Enabled:   true,
		Username:  username,
		Password:  password,
		Mechanism: c.SASL.Mechanism,
	}

	return nil
}

// LoadConfig loads and validates the configuration
func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Load credentials for each cluster
	for name, cluster := range config.Clusters {
		if err := cluster.LoadSASLCredentials(name); err != nil {
			return nil, err
		}
		config.Clusters[name] = cluster
	}

	// Convert seconds to duration
	config.Global.Timeout *= time.Second
	config.Kafka.AdminTimeout *= time.Second
	config.Kafka.DefaultRefreshRate *= time.Second

	// Set defaults
	if config.Global.Port == "" {
		config.Global.Port = ":9899"
	}
	if config.Global.Timeout == 0 {
		config.Global.Timeout = 3 * time.Second
	}
	if config.Kafka.AdminTimeout == 0 {
		config.Kafka.AdminTimeout = 5 * time.Second
	}
	if config.Kafka.DefaultRefreshRate == 0 {
		config.Kafka.DefaultRefreshRate = 60 * time.Second
	}

	return &config, nil
}
