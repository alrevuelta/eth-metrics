package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alrevuelta/eth-pools-metrics/config"
	"github.com/alrevuelta/eth-pools-metrics/metrics"
	"github.com/alrevuelta/eth-pools-metrics/price"
	"github.com/alrevuelta/eth-pools-metrics/prometheus"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := config.NewCliConfig()
	if err != nil {
		log.Fatal(err)
	}

	logLevel, err := log.ParseLevel(config.Verbosity)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)

	prometheus.Run(config.PrometheusPort)

	metrics, err := metrics.NewMetrics(
		context.Background(),
		config)

	if err != nil {
		log.Fatal(err)
	}

	price, err := price.NewPrice(config.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	go price.Run()
	metrics.Run()

	// Wait for signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	for {
		sig := <-sigCh
		if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == os.Interrupt || sig == os.Kill {
			break
		}
	}

	log.Info("Stopping eth-pools-metrics")
}
