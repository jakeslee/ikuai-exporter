package main

import (
	"github.com/jakeslee/ikuai"
	"github.com/jakeslee/ikuai-exporter/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	i := ikuai.NewIKuai("http://10.10.1.253", "test", "test123")

	registry := prometheus.NewRegistry()

	registry.MustRegister(pkg.NewIKuaiExporter(i))

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))

	log.Printf("exporter started at :9090")

	log.Fatal(http.ListenAndServe(":9090", nil))
}
