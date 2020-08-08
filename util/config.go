package util

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

var (
	// Configuration is main conig
	Configuration TomlConfig
)

// TomlConfig is the main config struct
type TomlConfig struct {
	Global struct {
		Port string `toml:"port"`
	} `toml:"global"`
	Log struct {
		Format string `toml:"format"`
		Level  string `toml:"level"`
	} `toml:"log"`
	Kafka struct {
		DefaultRefreshRate int    `toml:"default_refresh_rate"`
		MinKafkaVersion    string `toml:"min_kafka_version"`
		AdminTimeout       int    `toml:"admin_timeout"`
	} `toml:"kafka"`
	Clusters map[string]clusters
}

type clusters struct {
	Brokers     []string
	TopicFilter string
}

func init() {
	Configuration = initConfig()
}

// initConfig load the config from path, parse it and return TomlConfig
func initConfig() TomlConfig {
	log.Debug("Config init")
	var config TomlConfig
	if _, err := toml.DecodeFile("/opt/kcm/kcm.toml", &config); err != nil {
		log.Panicf("Unble to load config file.\n%v", err)
	}
	// Set custom defaults
	if config.Kafka.AdminTimeout == 0 {
		config.Kafka.AdminTimeout = 5
	}
	return config
}
