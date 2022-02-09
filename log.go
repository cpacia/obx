package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"obx/net"
	"obx/repo"
	"path"
	"strings"
)

var LogLevelMap = map[string]zapcore.Level{
	"debug":     zap.DebugLevel,
	"info":      zap.InfoLevel,
	"warning":   zap.WarnLevel,
	"error":     zap.ErrorLevel,
	"alert":     zap.DPanicLevel,
	"critical":  zap.PanicLevel,
	"emergency": zap.FatalLevel,
}

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}

func setupLogging(logDir, level string, testnet bool) error {
	var cfg zap.Config
	if testnet {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	logLevel, ok := LogLevelMap[strings.ToLower(level)]
	if !ok {
		return errors.New("invalid log level")
	}
	cfg.Encoding = "console"
	cfg.Level = zap.NewAtomicLevelAt(logLevel)

	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + logLevelSeverity[level] + "]")
	}
	cfg.EncoderConfig.EncodeLevel = customLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var (
		logger *zap.Logger
		err    error
	)
	if logDir != "" {
		logRotator := &lumberjack.Logger{
			Filename:   path.Join(logDir, repo.DefaultLogFilename),
			MaxSize:    10, // Megabytes
			MaxBackups: 3,
			MaxAge:     30, // Days
		}

		lumberjackZapHook := func(e zapcore.Entry) error {
			logRotator.Write([]byte(fmt.Sprintf("%+v\n", e)))
			return nil
		}

		logger, err = cfg.Build(zap.Hooks(lumberjackZapHook))
		if err != nil {
			return err
		}
	} else {
		logger, err = cfg.Build()
		if err != nil {
			return err
		}
	}
	zap.ReplaceGlobals(logger)

	log = zap.S()
	repo.UpdateLogger()
	net.UpdateLogger()
	return nil
}
