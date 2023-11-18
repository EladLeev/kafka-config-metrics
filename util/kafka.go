package util

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

func parseKafkaVersion(kafkaVersion string) sarama.KafkaVersion {
	version, err := sarama.ParseKafkaVersion(kafkaVersion)
	if err != nil {
		log.Panicf("Unknown Kafka Version:\n%v", err)
	}
	return version
}

// OpenConnection to one of a given broker using supplied version
func OpenConnection(kafkaVersion string, clusterAddr []string, adminTimeout int) sarama.ClusterAdmin {
	config := sarama.NewConfig()
	config.Version = parseKafkaVersion(kafkaVersion)
	config.Admin.Timeout = time.Duration(adminTimeout)

	clusterAdmin, err := sarama.NewClusterAdmin(clusterAddr, config)
	if err != nil {
		log.Errorf("Unable to connect to cluster: %v.\n%v", clusterAddr, err)
	}
	return clusterAdmin
}

// describeTopicConfig Gets topicName and ClusterAdmin interface and return []ConfigEntry
func describeTopicConfig(topicName string, clusterAdmin sarama.ClusterAdmin, wg *sync.WaitGroup, collector *Collector) {
	defer wg.Done()
	resource := sarama.ConfigResource{
		Type: sarama.TopicResource,
		Name: topicName,
	}

	r, err := clusterAdmin.DescribeConfig(resource)
	if err != nil {
		log.Errorf("Unable to DescribeConfig:\n%v", err)
		return
	}
	RegisterMetrics(r, topicName, collector)
	//log.Debug(r)
	// for _, v := range r {
	// 	if s, err := strconv.Atoi(v.Value); err == nil {
	// 		// Register metrics
	// 		collector.minCompactionLagMs.With(prometheus.Labels{"topic": topicName}).Set(float64(s))
	// 	}
	// }
	return
}

// listTopics on a given cluster
func listTopics(clusterAdmin sarama.ClusterAdmin) map[string]sarama.TopicDetail {
	r, err := clusterAdmin.ListTopics()
	if err != nil {
		//TODO
		log.Errorf("Unable to ListTopics:\n%v", err)
	}
	return r
}

func filterTopic(filterRegex string) *regexp.Regexp {
	re, err := regexp.Compile(filterRegex)
	if err != nil {
		fmt.Print(err)
	}
	return re
}

// PullConfigs opens a connection to a given cluster, filter topics based on the config file
// and then run DescribeTopicConfig to populate the metrics
func PullConfigs(cfg TomlConfig, clusterName string, collector *Collector) {
	clusterBrokers := cfg.Clusters[clusterName].Brokers
	clusterAdmin := OpenConnection(cfg.Kafka.MinKafkaVersion, clusterBrokers, cfg.Kafka.AdminTimeout)
	log.Debugf("KafkaVersion: %v | BrokerList: %v", cfg.Kafka.MinKafkaVersion, clusterBrokers)
	defer clusterAdmin.Close()

	// Pull topic list and topic configs
	var wg sync.WaitGroup
	var re *regexp.Regexp
	topicList := listTopics(clusterAdmin)
	// If `topicfilter` config is not empty, compile Regex
	if cfg.Clusters[clusterName].TopicFilter != "" {
		re = filterTopic(cfg.Clusters[clusterName].TopicFilter)
	}
	for topic := range topicList {
		if (re != nil) && re.MatchString(topic) {
			continue
		}
		wg.Add(1)
		go describeTopicConfig(topic, clusterAdmin, &wg, collector)
	}
	wg.Wait()
}
