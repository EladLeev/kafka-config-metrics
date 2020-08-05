package util

import (
	"testing"

	"github.com/Shopify/sarama"
)

func TestParseKafkaVersion(t *testing.T) {
	c := parseKafkaVersion("0.10.1.1")
	if c != sarama.KafkaVersion(sarama.V0_10_1_1) {
		t.Errorf("Something is wrong with Sarama version: %v", c)
	}
}
