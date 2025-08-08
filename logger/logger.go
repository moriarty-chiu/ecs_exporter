package logger

import (
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"ecs_exporter/config"
)

func Init(cfg config.LogConfig) {
	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logrus.Warnf("Invalid log level: %s, using default level: info", cfg.Level)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// 设置日志格式
	var formatter logrus.Formatter

	if cfg.Format == "json" {
		formatter = &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "@level",
				logrus.FieldKeyMsg:   "@message",
			},
		}
	} else {
		formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}
	}

	logrus.SetFormatter(formatter)

	// 如果启用了文件日志，则配置日志轮转
	if cfg.EnableFile {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.Errorf("Failed to create log directory %s: %v", logDir, err)
			return
		}

		// 创建日志轮转配置 - 按日期生成文件名
		logPattern := cfg.Filename + ".%Y%m%d" // 生成格式如: app_20231201.log
		writer, err := rotatelogs.New(
			logPattern,
			rotatelogs.WithLinkName(cfg.Filename), // 软链接到当前日志文件
			rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge)*24*time.Hour), // 日志保留天数
			rotatelogs.WithRotationTime(24*time.Hour),                     // 每天轮转一次
			rotatelogs.WithRotationSize(int64(cfg.MaxSize)*1024*1024),     // 按文件大小轮转(MB)
		)

		if err != nil {
			logrus.Errorf("Failed to initialize log rotator: %v", err)
			return
		}

		// 设置日志钩子，同时写入文件和控制台
		hook := lfshook.NewHook(
			lfshook.WriterMap{
				logrus.DebugLevel: writer,
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.FatalLevel: writer,
				logrus.PanicLevel: writer,
			},
			formatter,
		)

		logrus.AddHook(hook)
	}
}
