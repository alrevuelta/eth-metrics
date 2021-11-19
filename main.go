package main

import (
	"context"
	"eth-pools-metrics/metrics"    // TODO: Set github prefix when released
	"eth-pools-metrics/prometheus" // TODO: Set github prefix when released
	log "github.com/sirupsen/logrus"
	"time"
)

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
		config.WithdrawalCredentials)

	if err != nil {
		log.Fatal("Error creating new metrics: ", err)
	}

	metrics.Run()

	// Loop forever
	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
