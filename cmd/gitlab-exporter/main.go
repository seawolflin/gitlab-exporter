package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seawolflin/gitlab-exporter/internal/collector"
	"log"
	"net/http"
)

func main() {
	u := collector.NewUserCollector()
	prometheus.MustRegister(u)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
