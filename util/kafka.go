package util

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Shopify/sarama"
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
func OpenConnection(kafkaVersion string, clusterAddr []string) sarama.ClusterAdmin {
	config := sarama.NewConfig()
	//TODO: Make sure min version
	config.Version = parseKafkaVersion(kafkaVersion)
	//config.Admin.Timeout = time.Duration(cfg.Kafka.AdminTimeout) # TODO

	clusterAdmin, err := sarama.NewClusterAdmin(clusterAddr, config)
	if err != nil {
		log.Errorf("Unable to connect to cluster.\n%v", err)
	}
	return clusterAdmin
}

// DescribeTopicConfig Gets topicName and ClusterAdmin interface and return []ConfigEntry
func DescribeTopicConfig(topicName string, clusterAdmin sarama.ClusterAdmin, wg *sync.WaitGroup) {
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
	fmt.Printf("> Topic: %v\n\n", topicName)
	for _, v := range r {
		if s, err := strconv.Atoi(v.Value); err == nil {
			fmt.Printf("%v: %v \n", v.Name, s)
		}
	}
	fmt.Printf("\n###############################################################\n")
}

// ListTopics on a given cluster
func ListTopics(clusterAdmin sarama.ClusterAdmin) map[string]sarama.TopicDetail {
	r, err := clusterAdmin.ListTopics()
	if err != nil {
		//TODO
		log.Errorf("Unable to ListTopics:\n%v", err)
	}
	return r
}
