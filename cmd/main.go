package main

import (
	"ecs_exporter/collector"
	"ecs_exporter/config"
	"ecs_exporter/logger"
	"ecs_exporter/token"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sirupsen/logrus"
)

func main() {
	// 1. 读取配置
	config.LoadConfig("config/config.yaml")

	// 2. 初始化日志
	logger.Init(config.Cfg.Log)
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
	logrus.Info("Starting server at :9100")
	if err := http.ListenAndServe(":9100", nil); err != nil {
		logrus.Fatalf("server error: %v", err)
	}
}
