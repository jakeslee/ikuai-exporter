/*
Copyright © 2026 Jakes Lee
*/
package cmd

import (
	"net/http"
	"strings"

	"github.com/jakeslee/ikuai"
	"github.com/jakeslee/ikuai-exporter/cmd/options"
	"github.com/jakeslee/ikuai-exporter/pkg"
	"github.com/jakeslee/ikuai-exporter/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var opts = options.NewServerOptions()

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run metrics endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		level, err := logrus.ParseLevel(opts.Level)
		if err != nil {
			return err
		}
		logrus.SetLevel(level)

		i := ikuai.NewIKuai(strings.TrimSpace(opts.URL), opts.Username, opts.Password, opts.InsecureSkip, true)

		if level >= logrus.DebugLevel {
			i.Debug()
		}

		registry := prometheus.NewRegistry()
		registry.MustRegister(pkg.NewIKuaiExporter(i))

		http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))

		logrus.Infof("iKuai exporter %v started on :9090", version.Version())
		logrus.Fatal(http.ListenAndServe(":9090", nil))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().StringVar(&opts.URL, "url", opts.URL, "iKuai URL")
	serverCmd.Flags().StringVarP(&opts.Username, "username", "u", opts.Username, "iKuai username")
	serverCmd.Flags().StringVarP(&opts.Password, "password", "p", opts.Password, "The password for the user on iKuai")
	serverCmd.Flags().BoolVar(&opts.InsecureSkip, "insecure-skip", opts.InsecureSkip, "Skip iKuai certificate verification")
	serverCmd.Flags().StringVarP(&opts.Level, "level", "l", opts.Level, "Log level")

	viper.BindEnv("url", "IK_URL")
	viper.BindEnv("username", "IK_USER")
	viper.BindEnv("password", "IK_PWD")
}
