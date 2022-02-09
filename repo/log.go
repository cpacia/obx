package repo

import "go.uber.org/zap"

var log = zap.S()

func UpdateLogger() {
	log = zap.S()
}
