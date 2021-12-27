package main

import (
	"context"
	"github.com/alrevuelta/eth-pools-metrics/metrics"
	"github.com/alrevuelta/eth-pools-metrics/prometheus"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// TODO: Bump automatically with -ldflags
// go build -v -ldflags="-X 'main.ReleaseVersion=x.y.z'"
var ReleaseVersion = "0.0.3"

func main() {
	config, err := NewCliConfig()
	if err != nil {
		log.Fatal("Error creating cli config", err)
	}

	prometheus.Run(config.PrometheusPort)

	metrics, err := metrics.NewMetrics(
		context.Background(),
		config.BeaconRpcEndpoint,
		config.Network,
		config.WithdrawalCredentials,
		config.FromAddress)

	if err != nil {
		log.Fatal("Error creating new metrics: ", err)
	}

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
