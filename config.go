package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Network               string
	WithdrawalCredentials []string
	FromAddress           []string
	BeaconRpcEndpoint     string
	PrometheusPort        int
}

// custom implementation to allow providing the same flag multiple times
// --flag=value1 --flag=value2
type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func NewCliConfig() (*Config, error) {
	var fromAddress arrayFlags
	var withdrawalCredentials arrayFlags

	flag.Var(&withdrawalCredentials, "withdrawal-credentials", "Withdrawal credentials used in the deposits. Can be used multiple times")
	flag.Var(&fromAddress, "from-address", "Wallet addresses used to deposit. Can be used multiple times")

	var network = flag.String("network", "mainnet", "Ethereum 2.0 network mainnet|prater|pyrmont")
	var beaconRpcEndpoint = flag.String("beacon-rpc-endpoint", "localhost:4000", "Address:Port of a eth2 beacon node endpoint")
	var prometheusPort = flag.Int("prometheus-port", 9500, "Prometheus port to listen to")

	flag.Parse()

	conf := &Config{
		Network:               *network,
		BeaconRpcEndpoint:     *beaconRpcEndpoint,
		PrometheusPort:        *prometheusPort,
		WithdrawalCredentials: withdrawalCredentials,
		FromAddress:           fromAddress,
	}
	logConfig(conf)
	return conf, nil
}

func logConfig(cfg *Config) {
	log.WithFields(log.Fields{
		"BeaconRpcEndpoint":     cfg.BeaconRpcEndpoint,
		"WithdrawalCredentials": cfg.WithdrawalCredentials,
		"FromAddress":           cfg.FromAddress,
		"Network":               cfg.Network,
		"PrometheusPort":        cfg.PrometheusPort,
	}).Info("Cli Config:")
}
