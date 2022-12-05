package price

import (
	"time"

	"github.com/alrevuelta/eth-pools-metrics/config"
	"github.com/alrevuelta/eth-pools-metrics/postgresql"
	"github.com/alrevuelta/eth-pools-metrics/prometheus"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	gecko "github.com/superoo7/go-gecko/v3"
)

var vc = []string{"usd", "eurr"}

type Price struct {
	postgresql *postgresql.Postgresql
	coingecko  *gecko.Client
}

func NewPrice(postgresEndpoint string) (*Price, error) {

	cg := gecko.NewClient(nil)

	var pg *postgresql.Postgresql
	var err error
	if postgresEndpoint != "" {
		pg, err = postgresql.New(postgresEndpoint)
		if err != nil {
			return nil, errors.Wrap(err, "could not create postgresql")
		}
		err = pg.CreateEthPriceTable()
		if err != nil {
			return nil, errors.Wrap(err, "error creating pool table to store data")
		}
	}

	return &Price{
		postgresql: pg,
		coingecko:  cg,
	}, nil
}

func (p *Price) GetEthPrice() {
	id := ""
	if config.Network == "mainnet" {
		id = "ethereum"
	} else if config.Network == "gnosis" {
		id = "gnosis"
	} else {
		log.Fatal("Network not supported: ", config.Network)
	}

	sp, err := p.coingecko.SimplePrice([]string{id}, vc)
	if err != nil {
		log.Error(err)
	}

	eth := (*sp)[id]
	ethPriceUsd := eth["usd"]

	logPrice(ethPriceUsd)
	setPrometheusPrice(ethPriceUsd)

	if p.postgresql != nil {
		err := p.postgresql.StoreEthPrice(ethPriceUsd)
		if err != nil {
			log.Error(err)
		}
	}
}

func (p *Price) Run() {
	todoSetAsFlag := 30 * time.Minute
	ticker := time.NewTicker(todoSetAsFlag)
	for ; true; <-ticker.C {
		p.GetEthPrice()
	}
}

func logPrice(price float32) {
	log.Info("Ethereum price in USD: ", price)
}

func setPrometheusPrice(price float32) {
	prometheus.EthereumPriceUsd.Set(float64(price))
}
