package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"ecs_exporter/config"
)

func Init(cfg config.LogConfig) {
	level, err := logrus.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		logrus.Fatalf("failed to create log dir: %v", err)
	}

	logFile := filepath.Join(cfg.Dir, cfg.File)
	rotateLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	logrus.SetOutput(io.MultiWriter(os.Stdout, rotateLogger))
}
