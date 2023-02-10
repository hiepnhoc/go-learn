package main

import (
	"2margin.vn/account-service/config"
	"2margin.vn/account-service/pkg/logger"
	"2margin.vn/account-service/pkg/utils"
	"log"
	"os"
)

func main() {

	log.Println("Starting account service microservice")
	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config failed : %v", err)
	}

	appLogger := logger.NewAppLogger(cfg)
	appLogger.InitLogger()

	appLogger.Infof(
		"Service: %s, AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		cfg.ServiceName,
		cfg.Server.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Mode,
		cfg.Server.SSL,
	)

}
