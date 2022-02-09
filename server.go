package main

import (
	"go.uber.org/zap"
	"obx/repo"
)

var log = zap.S()

type Server struct {
}

func main() {
	config, err := repo.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := setupLogging(config.LogDir, config.LogLevel, config.Testnet); err != nil {
		log.Fatal(err)
	}

}
