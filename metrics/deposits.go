package metrics

import (
	"eth-pools-metrics/prometheus" // TODO: Set github prefix when released
	log "github.com/sirupsen/logrus"
	"time"
)

// TODO: Temporal solution:
// - TheGraph API calls has some limits, so we can't query in every epoch
// - Race condition with the depositedKeys
// - Fetches the deposits every hour
func (a *Metrics) StreamDeposits() {
	for {
		pubKeysDeposited, err := a.theGraph.GetAllDepositedKeys()
		if err != nil {
			log.Error(err)
			time.Sleep(60 * 10 * time.Second)
			continue
		}
		log.Info("Number of deposited keys: ", len(pubKeysDeposited))
		a.depositedKeys = pubKeysDeposited
		prometheus.NOfDepositedValidators.Set(float64(len(pubKeysDeposited)))
		time.Sleep(60 * 60 * time.Second)
	}
}
