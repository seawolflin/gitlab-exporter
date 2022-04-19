package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"log"
	"net/http"
)

// 匿名导入，为了执行init方法
import (
	_ "github.com/seawolflin/gitlab-exporter/internal/collector"
	_ "github.com/seawolflin/gitlab-exporter/internal/db"
	_ "github.com/seawolflin/gitlab-exporter/internal/gitlab"
)

func main() {
	context.GetInstance().Parse()

	initializer.InitAll()

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
