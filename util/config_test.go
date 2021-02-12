package util

import (
	"testing"
)

// Validate Confioguration
func TestInitConfig(t *testing.T) {
	cfg := InitConfig()

	if len(cfg.Clusters) == 0 {
		t.Error("Unable to find clusters")
	}

	if cfg.Global.Port <= 0 {
		t.Errorf("Invalid port %v", cfg.Global.Port)
	}
}
