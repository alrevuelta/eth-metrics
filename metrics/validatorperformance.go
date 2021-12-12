package metrics

import (
	"context"
	"encoding/hex"
	"eth-pools-metrics/prometheus" // TODO: Set github prefix when released
	"github.com/pkg/errors"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/v2/config/params"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v2/time/slots"
	log "github.com/sirupsen/logrus"
	"time"
)

func (a *Metrics) StreamValidatorPerformance() {
	for {
		time.Sleep(2 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Fetch needed data to run the metrics
		newData, err := a.FetchValidatorPerformance(ctx)
		if err != nil {
			log.WithError(err).Warn("Failed to fetch metrics data")
			continue
		}

		if !newData {
			continue
		}

		// Calculate the metrics
		nOfTotalVotes, nOfIncorrectSource, nOfIncorrectTarget, nOfIncorrectHead := a.getIncorrectAttestations()
		nOfValsWithDecreasedBalance, nOfValidators := a.getNumOfBalanceDecreasedVals()
		balanceDecreasedPercent := (float64(nOfValsWithDecreasedBalance) / float64(nOfValidators)) * 100

		// Log the metrics
		logEpochSlot := log.WithFields(log.Fields{
			"Epoch": a.Epoch,
			"Slot":  a.Slot,
		})

		logEpochSlot.WithFields(log.Fields{
			"nOfTotalVotes":      nOfTotalVotes,
			"nOfIncorrectSource": nOfIncorrectSource,
			"nOfIncorrectTarget": nOfIncorrectTarget,
			"nOfIncorrectHead":   nOfIncorrectHead,
		}).Info("Incorrect voting:")

		logEpochSlot.WithFields(log.Fields{
			"ActiveValidators":    len(a.activeKeys),
			"DepositedValidators": len(a.depositedKeys),
			"SlashedValidators":   "TODO",
			"ExitingValidators":   "TODO",
			"OtherStates":         "TODO",
		}).Info("Validator Status:")

		logEpochSlot.WithFields(log.Fields{
			"PercentIncorrectSource": (float64(nOfIncorrectSource) / float64(nOfTotalVotes)) * 100,
			"PercentIncorrectTarget": (float64(nOfIncorrectTarget) / float64(nOfTotalVotes)) * 100,
			"PercentIncorrectHead":   (float64(nOfIncorrectHead) / float64(nOfTotalVotes)) * 100,
		}).Info("Incorrect voting percents:")

		logEpochSlot.WithFields(log.Fields{
			"nOfValidators":               len(a.activeKeys),
			"nOfValsWithDecreasedBalance": nOfValsWithDecreasedBalance,
			"balanceDecreasedPercent":     balanceDecreasedPercent,
		}).Info("Balance decreased:")

		// Update prometheus metrics
		prometheus.NOfTotalVotes.Set(float64(nOfTotalVotes))
		prometheus.NOfIncorrectSource.Set(float64(nOfIncorrectSource))
		prometheus.NOfIncorrectTarget.Set(float64(nOfIncorrectTarget))
		prometheus.NOfIncorrectHead.Set(float64(nOfIncorrectHead))
		prometheus.BalanceDecreasedPercent.Set(balanceDecreasedPercent)
	}
}

//Fetches data from the beacon chain for a given set of validators. Note
//that not all request accepts the epoch as input, so this function takes
//care of synching with the beacon so that all fetched data refers to the same
//epoch
func (a *Metrics) FetchValidatorPerformance(ctx context.Context) (bool, error) {
	head, err := GetChainHead(ctx, a.beaconChainClient)
	if err != nil {
		return false, errors.Wrap(err, "error getting chain head")
	}

	// Run metrics in already completed epochs
	metricsEpoch := uint64(head.HeadEpoch) - 1
	metricsSlot := uint64(head.HeadSlot)

	log.Info("Slot: ", ethTypes.Slot(metricsSlot)%params.BeaconConfig().SlotsPerEpoch)

	if a.depositedKeys == nil {
		log.Warn("No active keys to get vals performance")
		time.Sleep(30 * time.Second)
		return false, nil
	}

	// Wait until the last slot to ensure all attestations are included
	if a.Epoch >= metricsEpoch || !slots.IsEpochEnd(head.HeadSlot) {
		return false, nil
	}

	slotTime, err := slots.ToTime(uint64(a.genesisSeconds), ethTypes.Slot(head.HeadSlot+1))

	if err != nil {
		return false, errors.Wrap(err, "could not get next slot time")
	}

	// Set as deadline the begining of the first slot of the next epoch
	ctx, cancel := context.WithDeadline(ctx, slotTime)
	defer cancel()

	a.Epoch = metricsEpoch
	a.Slot = metricsSlot

	log.WithFields(log.Fields{
		"Epoch": metricsEpoch,
		"Slot":  metricsSlot,
		// zero-indexed
		"SlotInEpoch": ethTypes.Slot(metricsSlot) % params.BeaconConfig().SlotsPerEpoch,
	}).Info("Fetching new validators info")

	req := &ethpb.ValidatorPerformanceRequest{
		PublicKeys: a.activeKeys,
	}

	valsPerformance, err := a.beaconChainClient.GetValidatorPerformance(ctx, req)
	if err != nil {
		return false, errors.Wrap(err, "could not get validator performance from beacon client")
	}

	a.valsPerformance = valsPerformance

	for i := range valsPerformance.MissingValidators {
		log.WithFields(log.Fields{
			"Epoch":   a.Epoch,
			"Address": hex.EncodeToString(a.valsPerformance.MissingValidators[i]),
		}).Warn("Validator performance not found in beacon chain")
	}

	log.Info("Remaining time for next slot: ", ctx)

	return true, nil
}

// Gets the total number of votes and the incorrect ones
// The source is the attestation itself
// https://pintail.xyz/posts/validator-rewards-in-practice/?s=03#attestation-efficiency
func (a *Metrics) getIncorrectAttestations() (uint64, uint64, uint64, uint64) {
	nOfIncorrectSource := uint64(0)
	nOfIncorrectTarget := uint64(0)
	nOfIncorrectHead := uint64(0)
	for i := range a.valsPerformance.PublicKeys {
		nOfIncorrectSource += BoolToUint64(!a.valsPerformance.CorrectlyVotedSource[i])
		nOfIncorrectTarget += BoolToUint64(!a.valsPerformance.CorrectlyVotedTarget[i])
		nOfIncorrectHead += BoolToUint64(!a.valsPerformance.CorrectlyVotedHead[i])
		// since missing source is the most severe, log it
		if !a.valsPerformance.CorrectlyVotedSource[i] {
			log.Info("Key that missed the attestation: ", hex.EncodeToString(a.valsPerformance.PublicKeys[i]), "--", a.valsPerformance.CorrectlyVotedSource[i], "--", a.valsPerformance.BalancesAfterEpochTransition[i])
		}
	}

	// Each validator contains three votes: source, target and head
	nOfTotalVotes := uint64(len(a.valsPerformance.PublicKeys)) * 3

	return nOfTotalVotes, nOfIncorrectSource, nOfIncorrectTarget, nOfIncorrectHead
}

// Gets the total number of validators and the ones that decreased in value
func (a *Metrics) getNumOfBalanceDecreasedVals() (uint64, uint64) {
	nOfValsWithDecreasedBalance := uint64(0)
	for i := range a.valsPerformance.PublicKeys {
		if a.valsPerformance.BalancesAfterEpochTransition[i] < a.valsPerformance.BalancesBeforeEpochTransition[i] {
			log.Info("Key with decr balance: ", hex.EncodeToString(a.valsPerformance.PublicKeys[i]), "--", a.valsPerformance.BalancesBeforeEpochTransition[i], "--", a.valsPerformance.BalancesAfterEpochTransition[i])
			nOfValsWithDecreasedBalance++
		}
	}
	nOfValidators := uint64(len(a.valsPerformance.PublicKeys))

	return nOfValsWithDecreasedBalance, nOfValidators
}
