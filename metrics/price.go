package metrics

import (
	"github.com/alrevuelta/eth-pools-metrics/prometheus"
	log "github.com/sirupsen/logrus"
	gecko "github.com/superoo7/go-gecko/v3"
	"runtime"
	"time"
)

var cg = gecko.NewClient(nil)
var ids = []string{"ethereum"}
var vc = []string{"usd", "eurr"}

func (a *Metrics) StreamEthPrice() {
	for {
		// Get eth price from coingecko
		sp, err := cg.SimplePrice(ids, vc)
		if err != nil {
			log.Error(err)
		}
		eth := (*sp)["ethereum"]

		logPrice(eth["usd"])
		setPrometheusPrice(eth["usd"])

		// Temporal fix to memory leak. Perhaps having an infinite loop
		// inside a routinne is not a good idea. TODO
		runtime.GC()

		// Every hour
		time.Sleep(60 * 60 * time.Second)
	}
}

func logPrice(price float32) {
	log.Info("Ethereum price in USD: ", price)
}

func setPrometheusPrice(price float32) {
	prometheus.EthereumPriceUsd.Set(float64(price))
}
