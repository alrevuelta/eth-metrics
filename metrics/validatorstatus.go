package metrics

import (
	"context"
	"eth-pools-metrics/prometheus" // TODO: Set github prefix when released
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"time"
)

type StatusCount struct {
	// custom field
	Validating uint64

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
		statusCount := getStatusCount(context.Background(), valsStatus)
		logStatus(statusCount)
		setPrometheusStatusMetrics(statusCount)

		time.Sleep(6 * 60 * time.Second)
	}
}

func setPrometheusStatusMetrics(statusCount StatusCount) {
	prometheus.NOfValidatingValidators.Set(float64(statusCount.Validating))
	prometheus.NOfUnkownValidators.Set(float64(statusCount.Unknown))
	prometheus.NOfDepositedValidators.Set(float64(statusCount.Deposited))
	prometheus.NOfPendingValidators.Set(float64(statusCount.Pending))
	prometheus.NOfActiveValidators.Set(float64(statusCount.Active))
	prometheus.NOfExitingValidators.Set(float64(statusCount.Exiting))
	prometheus.NOfSlashingValidators.Set(float64(statusCount.Slashing))
	prometheus.NOfExitedValidators.Set(float64(statusCount.Exited))
	prometheus.NOfInvalidValidators.Set(float64(statusCount.Invalid))
	prometheus.NOfPartiallyDepositedValidators.Set(float64(statusCount.PartiallyDeposited))
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

func getStatusCount(
	ctx context.Context,
	statusResponse *ethpb.MultipleValidatorStatusResponse) StatusCount {

	statusCount := StatusCount{}
	for i := range statusResponse.PublicKeys {
		status := statusResponse.Statuses[i].Status
		if isKeyValidating(status) {
			statusCount.Validating++
		} else if status == ethpb.ValidatorStatus_UNKNOWN_STATUS {
			statusCount.Unknown++
		} else if status == ethpb.ValidatorStatus_DEPOSITED {
			statusCount.Deposited++
		} else if status == ethpb.ValidatorStatus_PENDING {
			statusCount.Pending++
		} else if status == ethpb.ValidatorStatus_ACTIVE {
			statusCount.Active++
		} else if status == ethpb.ValidatorStatus_EXITING {
			statusCount.Exiting++
		} else if status == ethpb.ValidatorStatus_SLASHING {
			statusCount.Slashing++
		} else if status == ethpb.ValidatorStatus_EXITED {
			statusCount.Exited++
		} else if status == ethpb.ValidatorStatus_INVALID {
			statusCount.Invalid++
		} else if status == ethpb.ValidatorStatus_PARTIALLY_DEPOSITED {
			statusCount.PartiallyDeposited++
		}
	}
	return statusCount
}

func logStatus(statusCount StatusCount) {
	log.WithFields(log.Fields{
		"Validating":         statusCount.Validating,
		"Unknown":            statusCount.Unknown,
		"Deposited":          statusCount.Deposited,
		"Pending":            statusCount.Pending,
		"Active":             statusCount.Active,
		"Exiting":            statusCount.Exiting,
		"Slashing":           statusCount.Slashing,
		"Exited":             statusCount.Exited,
		"Invalid":            statusCount.Invalid,
		"PartiallyDeposited": statusCount.PartiallyDeposited,
	}).Info("Validator Status Count:")
}
