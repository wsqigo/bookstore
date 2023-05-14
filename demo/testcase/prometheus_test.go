package testcase

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

// 文档：https://cjting.me/2017/03/12/use-prometheus-to-monitor-server/

var httpRequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "http request count",
	}, []string{"endpoint"})

var httpRequestDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "http_request_duration",
		Help: "http request duration",
	},
	[]string{"endpoint"},
)

func init() {
	prometheus.MustRegister(httpRequestCount)
	prometheus.MustRegister(httpRequestDuration)
}

func TestPrometheus(t *testing.T) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/test", handler)
	go func() {
		http.ListenAndServe(":8888", nil)
	}()

	startClient()
	doneChan := make(chan struct{})
	<-doneChan
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	path := r.URL.Path
	httpRequestCount.WithLabelValues(path).Inc()

	n := rand.Intn(100)
	if n >= 95 {
		time.Sleep(100 * time.Millisecond)
	} else {
		time.Sleep(50 * time.Millisecond)
	}

	elapsed := float64(time.Since(start) / time.Millisecond)
	httpRequestDuration.WithLabelValues(path).Observe(elapsed)
}

func startClient() {
	sleepTime := 1000

	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		for {
			<-ticker.C
			sleepTime = 200
			<-time.After(30 * time.Second)
			sleepTime = 1000
		}
	}()

	for i := 0; i < 100; i++ {
		go func() {
			sendRequest()
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		}()
	}
}

func sendRequest() {
	resp, err := http.Get("http://localhost:8888/test")
	if err != nil {
		log.Println(err)
		return
	}
	resp.Body.Close()
}
