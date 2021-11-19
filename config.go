package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Config struct {
	Network               string
	WithdrawalCredentials []string
	BeaconRpcEndpoint     string
	PrometheusPort        int
}

func NewCliConfig() (*Config, error) {
	var network = flag.String("network", "mainnet", "Ethereum 2.0 network mainnet|prater|pyrmont")
	var withdrawalCredentials = flag.String("withdrawal-credentials", "", "Hex withdrawal credentials following the spec, without 0x prefix")
	var beaconRpcEndpoint = flag.String("beacon-rpc-endpoint", "localhost:4000", "Address:Port of a eth2 beacon node endpoint")
	var prometheusPort = flag.Int("prometheus-port", 9500, "Prometheus port to listen to")

	flag.Parse()

	conf := &Config{
		Network:               *network,
		BeaconRpcEndpoint:     *beaconRpcEndpoint,
		PrometheusPort:        *prometheusPort,
		WithdrawalCredentials: strings.Split(*withdrawalCredentials, ","),
	}
	logConfig(conf)
	return conf, nil
}

func logConfig(cfg *Config) {
	log.WithFields(log.Fields{
		"BeaconRpcEndpoint":     cfg.BeaconRpcEndpoint,
		"WithdrawalCredentials": cfg.WithdrawalCredentials,
		"Network":               cfg.Network,
		"PrometheusPort":        cfg.PrometheusPort,
	}).Info("Cli Config:")
}
