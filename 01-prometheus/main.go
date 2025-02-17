package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var opsProcessed *prometheus.CounterVec

var opsDuration *prometheus.HistogramVec

func setMeter() {
	opsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_request_total",
			Help: "Number of total request",
		},
		[]string{"URL", "status"},
	)

	opsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "app_request_duration",
			Help: "Duration of a request",
		},
		[]string{"URL", "status"},
	)

	prometheus.Register(opsProcessed)

	prometheus.Register(opsDuration)
}

func main() {
	setMeter()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respCode := 200

		defer func() {
			opsDuration.WithLabelValues("/", fmt.Sprint(respCode)).Observe(time.Since(start).Seconds())
		}()

		failQuery := r.URL.Query().Get("fail")

		if failQuery != "" {
			respCode = 500

			isFail, err := strconv.ParseBool(failQuery)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error parse query"))

				log.Println("error parse query")

				opsProcessed.WithLabelValues("/", fmt.Sprint(http.StatusInternalServerError)).Inc()

				return
			}

			if isFail {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("encounter fail"))

				log.Println("encounter fail")

				opsProcessed.WithLabelValues("/", fmt.Sprint(http.StatusInternalServerError)).Inc()

				return
			}
		}

		w.Write([]byte("Hello, World!"))

		opsProcessed.WithLabelValues("/", fmt.Sprint(http.StatusOK)).Inc()
	})
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on :9320")
	if err := http.ListenAndServe(":9320", nil); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
