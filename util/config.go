package util

import (
	"flag"

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
		Port    string `toml:"port"`
		Timeout int    `toml:"timeout"`
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
	// Parse the user input for a configuration file
	configPath := flag.String("config", "/opt/kcm/kcm.toml", "Absolute path to the configuration file")
	flag.Parse()
	log.Infof("Starting from config file:%v", *configPath)

	var config TomlConfig
	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		log.Panicf("Unble to load config file.\n%v", err)
	}
	// Set custom defaults
	if config.Kafka.AdminTimeout == 0 {
		config.Kafka.AdminTimeout = 5
	}
	if config.Global.Port == "" {
		config.Global.Port = ":9899"
	}
	if config.Global.Timeout == 0 {
		config.Global.Timeout = 3
	}
	return config
}
