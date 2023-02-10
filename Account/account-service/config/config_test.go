package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("config-local")
	if err != nil {
		t.Error(err)
	}
}

func TestGetConfig(t *testing.T) {

	config, err := GetConfig("config-local")
	if err != nil {
		t.Error(err)
	}

	if config.ServiceName != "ACCOUNT-SERVICE" {
		t.Errorf("got %s , want %s", config.ServiceName, "ACCOUNT-SERVICE")
	}

	_ = os.Setenv("SERVICENAME", "NEW-ACCOUNT-SERVICE")

	config, err = GetConfig("config-local")
	if err != nil {
		t.Error(err)
	}

	if config.ServiceName != "NEW-ACCOUNT-SERVICE" {
		t.Errorf("got %s , want %s", config.ServiceName, "NEW-ACCOUNT-SERVICE")
	}

}
