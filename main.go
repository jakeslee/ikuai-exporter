package main

import (
	"github.com/alexflint/go-arg"
	"github.com/jakeslee/ikuai"
	"github.com/jakeslee/ikuai-exporter/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type Config struct {
	Ikuai         string `arg:"env:IK_URL" help:"iKuai URL" default:"http://10.0.1.253"`
	IkuaiUsername string `arg:"env:IK_USER" help:"iKuai username" default:"test"`
	IkuaiPassword string `arg:"env:IK_PWD" help:"iKuai password" default:"test123"`
	Debug         bool   `arg:"env:DEBUG" help:"iKuai 开启 debug 日志" default:"false"`
	InsecureSkip  bool   `arg:"env:SKIP_TLS_VERIFY" help:"是否跳过 iKuai 证书验证" default:"true"`
}

var (
	version   string
	buildTime string
)

func main() {
	config := &Config{}
	arg.MustParse(config)

	i := ikuai.NewIKuai(config.Ikuai, config.IkuaiUsername, config.IkuaiPassword, config.InsecureSkip, true)

	if config.Debug {
		i.Debug()
	}

	registry := prometheus.NewRegistry()

	registry.MustRegister(pkg.NewIKuaiExporter(i))

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))

	log.Printf("exporter %v started at :9090", version)

	log.Fatal(http.ListenAndServe(":9090", nil))
}
