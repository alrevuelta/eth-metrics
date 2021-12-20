package metrics

import (
	"context"
	"github.com/alrevuelta/eth-pools-metrics/prometheus"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"time"
)

type ValidatorStatusMetrics struct {
	// custom field: vals with active duties
	Validating uint64

	// TODO: num of slashed validators
	// note that after slashing->exited

	// maps 1:1 with eth2 spec status
	Unknown            uint64
	Deposited          uint64
	Pending            uint64
	Active             uint64
	Exiting            uint64
	Slashing           uint64
	Exited             uint64
	Invalid            uint64
	PartiallyDeposited uint64
}

// TODO: Handle race condition
func (a *Metrics) StreamValidatorStatus() {
	for {
		if a.depositedKeys == nil {
			log.Warn("No depositedKeys are available")
			time.Sleep(10 * time.Second)
			continue
		}

		// Get the status of all the validators
		valsStatus, err := a.validatorClient.MultipleValidatorStatus(
			context.Background(),
			&ethpb.MultipleValidatorStatusRequest{
				PublicKeys: a.depositedKeys,
			})

		if err != nil {
			log.Error(err)
			time.Sleep(10 * time.Second)
			continue
		}

		// Get validators with active duties
		validatingKeys := filterValidatingValidators(valsStatus)
		a.validatingKeys = validatingKeys

		// For other status we just want the count
		metrics := getValidatorStatusMetrics(context.Background(), valsStatus)
		logValidatorStatus(metrics)
		setPrometheusValidatorStatus(metrics)

		time.Sleep(6 * 60 * time.Second)
	}
}

func setPrometheusValidatorStatus(metrics ValidatorStatusMetrics) {
	prometheus.NOfValidatingValidators.Set(float64(metrics.Validating))
	prometheus.NOfUnkownValidators.Set(float64(metrics.Unknown))
	prometheus.NOfDepositedValidators.Set(float64(metrics.Deposited))
	prometheus.NOfPendingValidators.Set(float64(metrics.Pending))
	prometheus.NOfActiveValidators.Set(float64(metrics.Active))
	prometheus.NOfExitingValidators.Set(float64(metrics.Exiting))
	prometheus.NOfSlashingValidators.Set(float64(metrics.Slashing))
	prometheus.NOfExitedValidators.Set(float64(metrics.Exited))
	prometheus.NOfInvalidValidators.Set(float64(metrics.Invalid))
	prometheus.NOfPartiallyDepositedValidators.Set(float64(metrics.PartiallyDeposited))
}

func filterValidatingValidators(vals *ethpb.MultipleValidatorStatusResponse) [][]byte {
	activeKeys := make([][]byte, 0)
	for i := range vals.PublicKeys {
		if isKeyValidating(vals.Statuses[i].Status) {
			activeKeys = append(activeKeys, vals.PublicKeys[i])
		}
	}
	return activeKeys
}

// Active as in having to fulfill duties
func isKeyValidating(status ethpb.ValidatorStatus) bool {
	if status == ethpb.ValidatorStatus_ACTIVE ||
		status == ethpb.ValidatorStatus_EXITING ||
		status == ethpb.ValidatorStatus_SLASHING {
		return true
	}
	return false
}

func getValidatorStatusMetrics(
	ctx context.Context,
	statusResponse *ethpb.MultipleValidatorStatusResponse) ValidatorStatusMetrics {

	metrics := ValidatorStatusMetrics{}
	for i := range statusResponse.PublicKeys {
		status := statusResponse.Statuses[i].Status
		if isKeyValidating(status) {
			metrics.Validating++
		} else if status == ethpb.ValidatorStatus_UNKNOWN_STATUS {
			metrics.Unknown++
		} else if status == ethpb.ValidatorStatus_DEPOSITED {
			metrics.Deposited++
		} else if status == ethpb.ValidatorStatus_PENDING {
			metrics.Pending++
		} else if status == ethpb.ValidatorStatus_ACTIVE {
			metrics.Active++
		} else if status == ethpb.ValidatorStatus_EXITING {
			metrics.Exiting++
		} else if status == ethpb.ValidatorStatus_SLASHING {
			metrics.Slashing++
		} else if status == ethpb.ValidatorStatus_EXITED {
			metrics.Exited++
		} else if status == ethpb.ValidatorStatus_INVALID {
			metrics.Invalid++
		} else if status == ethpb.ValidatorStatus_PARTIALLY_DEPOSITED {
			metrics.PartiallyDeposited++
		}
	}
	return metrics
}

func logValidatorStatus(metrics ValidatorStatusMetrics) {
	log.WithFields(log.Fields{
		"Validating":         metrics.Validating,
		"Unknown":            metrics.Unknown,
		"Deposited":          metrics.Deposited,
		"Pending":            metrics.Pending,
		"Active":             metrics.Active,
		"Exiting":            metrics.Exiting,
		"Slashing":           metrics.Slashing,
		"Exited":             metrics.Exited,
		"Invalid":            metrics.Invalid,
		"PartiallyDeposited": metrics.PartiallyDeposited,
	}).Info("Validator Status Count:")
}
