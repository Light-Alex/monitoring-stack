// monitoring-stack/main.go
package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	// 记录HTTP请求总数
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// 记录HTTP请求处理时间（秒）
	// Buckets: 自定义分桶，用于记录请求处理时间的分布
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func main() {
	// 创建logs目录
	// 判断是否存在logs目录
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.MkdirAll("./logs", 0755)
	}

	// 创建日志文件
	logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	// 配置日志（输出JSON格式，便于Promtail解析）
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	// 同时输出到文件和控制台
	logger.SetOutput(io.MultiWriter(logFile, os.Stdout))

	// 业务处理函数
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// 模拟处理
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		duration := time.Since(start).Seconds()
		// 记录指标
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		// 记录日志
		logger.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   200,
			"duration": duration,
		}).Info("HTTP request processed")
	})

	// 抛出错误日志
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		duration := time.Since(start).Seconds()
		// 记录指标
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, "500").Inc()
		// 记录日志
		logger.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   500,
			"duration": duration,
		}).Error("HTTP request processed with error")
	})

	// 暴露Prometheus指标端点
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Server starting on :8015...")
	if err := http.ListenAndServe(":8015", nil); err != nil {
		log.Fatal(err)
	}
}
