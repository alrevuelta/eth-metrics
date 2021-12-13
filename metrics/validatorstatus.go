package metrics

import (
	"context"
	"eth-pools-metrics/prometheus" // TODO: Set github prefix when released
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"time"
)

// TODO: Handle race condition
func (a *Metrics) StreamValidatorStatus() {
	for {
		if a.depositedKeys == nil {
			log.Warn("No depositedKeys are available")
			time.Sleep(10 * time.Second)
			continue
		}

		// Get the status of all the validators
		valsStatus, err := a.validatorClient.MultipleValidatorStatus(context.Background(), &ethpb.MultipleValidatorStatusRequest{
			PublicKeys: a.depositedKeys,
		})
		if err != nil {
			log.Error(err)
			time.Sleep(10 * time.Second)
			continue
		}

		// TODO: Get other status

		// Get the performance of the filtered validators
		activeKeys := FilterActiveValidators(valsStatus)
		a.activeKeys = activeKeys
		log.Info("Active validators: ", len(activeKeys))
		if len(a.activeKeys) == 0 {
			log.Error(err)
			time.Sleep(10 * time.Second)
			continue
		}

		prometheus.NOfValidators.Set(float64(len(a.activeKeys)))
		// TODO: Other status (slashed, etc)
		log.WithFields(log.Fields{
			"ActiveValidators":  len(a.activeKeys),
			"SlashedValidators": "TODO",
			"ExitingValidators": "TODO",
			"OtherStates":       "TODO",
		}).Info("Validator Status:")

		time.Sleep(6 * 60 * time.Second)
	}
}
