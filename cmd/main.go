package main

import (
	"context"
	"ecs_exporter/collector"
	"ecs_exporter/config"
	"ecs_exporter/logger"
	"ecs_exporter/token"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sirupsen/logrus"
)

func main() {
	// 1. 读取配置
	config.LoadConfig("config/config.yaml")

	// 2. 初始化日志
	logger.Init(config.Cfg.Log)

	// 确保日志目录存在
	if config.Cfg.Log.EnableFile {
		logDir := filepath.Dir(config.Cfg.Log.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.Errorf("Failed to create log directory %s: %v", logDir, err)
		}
	}

	logrus.Info("Logger initialized")

	// 3. 初始化 Token 管理
	tm := token.NewTokenManager(&config.Cfg.API)
	tm.Start()
	defer tm.Stop()

	// 4. 注册 Prometheus 采集器
	ecsCollector := collector.NewECSCollector(&config.Cfg.API, tm)
	prometheus.MustRegister(ecsCollector)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 5. 启动 HTTP 服务，暴露 /metrics
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr: ":9100",
	}

	// 在 goroutine 中启动服务器
	go func() {
		logrus.Info("Starting server at :9100")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("server error: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	// 上下文设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
}
