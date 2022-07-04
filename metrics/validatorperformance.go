package metrics

import (
	"context"
	"encoding/hex"
	"math/big"
	"runtime"
	"time"

	"github.com/alrevuelta/eth-pools-metrics/prometheus"
	"github.com/alrevuelta/eth-pools-metrics/schemas"
	"github.com/pkg/errors"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/v2/config/params"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v2/time/slots"
	log "github.com/sirupsen/logrus"
)

func (a *Metrics) StreamValidatorPerformance() {
	for {
		time.Sleep(2 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Fetch needed data to run the metrics
		valsPerformance, newData, epoch, err := a.FetchValidatorPerformance(ctx)
		if err != nil {
			log.WithError(err).Warn("Failed to fetch metrics data")
			continue
		}

		if !newData {
			continue
		}

		metrics := getValidatorPerformanceMetrics(valsPerformance)
		metrics.Epoch = epoch
		time, err := a.EpochToTime(epoch)
		if err != nil {
			log.Error(err)
		}
		metrics.Time = time
		//metrics.PoolName = a.PoolName

		logValidatorPerformance(metrics)
		setPrometheusValidatorPerformance(metrics)

		if a.postgresql != nil {
			err = a.postgresql.StoreValidatorPerformance(metrics)
			if err != nil {
				log.Error(err)
			}
		}

		// Temporal fix to memory leak. Perhaps having an infinite loop
		// inside a routinne is not a good idea. TODO
		runtime.GC()
	}
}

// Gets the total number of votes and the incorrect ones
// The source is the attestation itself
// https://pintail.xyz/posts/validator-rewards-in-practice/?s=03#attestation-efficiency
func getAttestationMetrics(valsPerformance *ethpb.ValidatorPerformanceResponse) (uint64, uint64, uint64, uint64, []string) {
	// TODO: Investigate InactivityScores
	nOfIncorrectSource := uint64(0)
	nOfIncorrectTarget := uint64(0)
	nOfIncorrectHead := uint64(0)
	missedAttestationsKeys := make([]string, 0)
	for i := range valsPerformance.PublicKeys {
		nOfIncorrectSource += BoolToUint64(!valsPerformance.CorrectlyVotedSource[i])
		nOfIncorrectTarget += BoolToUint64(!valsPerformance.CorrectlyVotedTarget[i])
		nOfIncorrectHead += BoolToUint64(!valsPerformance.CorrectlyVotedHead[i])
		if !valsPerformance.CorrectlyVotedSource[i] {
			missedAttestationsKeys = append(missedAttestationsKeys, hex.EncodeToString(valsPerformance.PublicKeys[i]))
		}
	}

	// Each validator contains three votes: source, target and head
	nOfTotalVotes := uint64(len(valsPerformance.PublicKeys)) * 3

	return nOfTotalVotes, nOfIncorrectSource, nOfIncorrectTarget, nOfIncorrectHead, missedAttestationsKeys
}

// Get metrics on balances in the epoch transition
func getBalanceMetrics(valsPerformance *ethpb.ValidatorPerformanceResponse) (
	uint64,
	uint64,
	*big.Int,
	*big.Int,
	[]string) {

	lostBalanceKeys := make([]string, 0)

	nOfValsWithDecreasedBalance := uint64(0)
	earnedBalance := big.NewInt(0)
	losedBalance := big.NewInt(0)
	for i := range valsPerformance.PublicKeys {
		delta := big.NewInt(0).Sub(big.NewInt(0).SetUint64(valsPerformance.BalancesAfterEpochTransition[i]), big.NewInt(0).SetUint64(valsPerformance.BalancesBeforeEpochTransition[i]))
		if delta.Cmp(big.NewInt(0)) == -1 {
			lostBalanceKeys = append(lostBalanceKeys, hex.EncodeToString(valsPerformance.PublicKeys[i]))
			nOfValsWithDecreasedBalance++
			losedBalance.Add(losedBalance, delta)
		} else {
			earnedBalance.Add(earnedBalance, delta)
		}
	}
	nOfValidators := uint64(len(valsPerformance.PublicKeys))

	return nOfValsWithDecreasedBalance, nOfValidators, earnedBalance, losedBalance, lostBalanceKeys
}

func getValidatorPerformanceMetrics(valsPerformance *ethpb.ValidatorPerformanceResponse) schemas.ValidatorPerformanceMetrics {
	metrics := schemas.ValidatorPerformanceMetrics{}

	nOfTotalVotes, nOfIncorrectSource, nOfIncorrectTarget, nOfIncorrectHead, missedKeys := getAttestationMetrics(valsPerformance)
	nOfValsWithDecreasedBalance, nOfValidators, earned, losed, lostKeys := getBalanceMetrics(valsPerformance)

	metrics.NOfTotalVotes = nOfTotalVotes
	metrics.NOfIncorrectSource = nOfIncorrectSource
	metrics.NOfIncorrectTarget = nOfIncorrectTarget
	metrics.NOfIncorrectHead = nOfIncorrectHead
	metrics.NOfValidatingKeys = nOfValidators
	metrics.NOfValsWithLessBalance = nOfValsWithDecreasedBalance
	metrics.EarnedBalance = earned
	metrics.LosedBalance = losed
	metrics.MissedAttestationsKeys = missedKeys
	metrics.LostBalanceKeys = lostKeys

	return metrics
}

func logValidatorPerformance(metrics schemas.ValidatorPerformanceMetrics) {
	balanceDecreasedPercent := (float64(metrics.NOfValsWithLessBalance) / float64(metrics.NOfValidatingKeys)) * 100

	logEpochSlot := log.WithFields(log.Fields{
		"Epoch": metrics.Epoch,
	})

	logEpochSlot.WithFields(log.Fields{
		"nOfTotalVotes":      metrics.NOfTotalVotes,
		"nOfIncorrectSource": metrics.NOfIncorrectSource,
		"nOfIncorrectTarget": metrics.NOfIncorrectTarget,
		"nOfIncorrectHead":   metrics.NOfIncorrectHead,
	}).Info("Incorrect voting:")

	logEpochSlot.WithFields(log.Fields{
		"PercentIncorrectSource": (float64(metrics.NOfIncorrectSource) / float64(metrics.NOfTotalVotes)) * 100,
		"PercentIncorrectTarget": (float64(metrics.NOfIncorrectTarget) / float64(metrics.NOfTotalVotes)) * 100,
		"PercentIncorrectHead":   (float64(metrics.NOfIncorrectHead) / float64(metrics.NOfTotalVotes)) * 100,
	}).Info("Incorrect voting percents:")

	logEpochSlot.WithFields(log.Fields{
		"nOfValidators":               metrics.NOfValidatingKeys,
		"nOfValsWithDecreasedBalance": metrics.NOfValsWithLessBalance,
		"balanceDecreasedPercent":     balanceDecreasedPercent,
		"earnedBalance":               metrics.EarnedBalance,
		"losedBalance":                metrics.LosedBalance,
	}).Info("Balance decreased:")

	for _, v := range metrics.MissedAttestationsKeys {
		logEpochSlot.WithFields(log.Fields{
			"ValidadorKey": v,
		}).Info("Validator missed attestation")
	}
	for _, v := range metrics.LostBalanceKeys {
		logEpochSlot.WithFields(log.Fields{
			"ValidadorKey": v,
		}).Info("Validator with less inter-epoch balance")
	}
}

func setPrometheusValidatorPerformance(metrics schemas.ValidatorPerformanceMetrics) {
	prometheus.NOfTotalVotes.Set(float64(metrics.NOfTotalVotes))
	prometheus.NOfIncorrectSource.Set(float64(metrics.NOfIncorrectSource))
	prometheus.NOfIncorrectTarget.Set(float64(metrics.NOfIncorrectTarget))
	prometheus.NOfIncorrectHead.Set(float64(metrics.NOfIncorrectHead))
	prometheus.EarnedAmountInEpoch.Set(float64(metrics.EarnedBalance.Int64()))
	prometheus.LosedAmountInEpoch.Set(float64(metrics.LosedBalance.Int64()))

	// TODO: Deprecate this, send the raw number
	balanceDecreasedPercent := (float64(metrics.NOfValsWithLessBalance) / float64(metrics.NOfValidatingKeys)) * 100
	prometheus.BalanceDecreasedPercent.Set(balanceDecreasedPercent)

	for _, v := range metrics.MissedAttestationsKeys {
		prometheus.MissedAttestationsKeys.WithLabelValues(v).Inc()
	}

	for _, v := range metrics.LostBalanceKeys {
		prometheus.LessBalanceKeys.WithLabelValues(v).Inc()
	}
}

// Fetches data from the beacon chain for a given set of validators. Note
// that not all request accepts the epoch as input, so this function takes
// care of synching with the beacon so that all fetched data refers to the same
// epoch
func (a *Metrics) FetchValidatorPerformance(ctx context.Context) (*ethpb.ValidatorPerformanceResponse, bool, uint64, error) {
	head, err := GetChainHead(ctx, a.beaconChainClient)
	if err != nil {
		return nil, false, 0, errors.Wrap(err, "error getting chain head")
	}

	metricsEpoch := uint64(head.HeadEpoch)
	metricsSlot := uint64(head.HeadSlot)

	//log.Info("Slot: ", ethTypes.Slot(metricsSlot)%params.BeaconConfig().SlotsPerEpoch)

	if a.validatingKeys == nil {
		log.Warn("No active keys to get vals performance")
		time.Sleep(30 * time.Second)
		return nil, false, 0, nil
	}

	// Wait until the last slot to ensure all attestations are included
	if a.Epoch >= metricsEpoch || !slots.IsEpochEnd(head.HeadSlot) {
		return nil, false, 0, nil
	}

	slotTime, err := slots.ToTime(uint64(a.genesisSeconds), ethTypes.Slot(head.HeadSlot+1))

	if err != nil {
		return nil, false, 0, errors.Wrap(err, "could not get next slot time")
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
		PublicKeys: a.validatingKeys,
	}

	valsPerformance, err := a.beaconChainClient.GetValidatorPerformance(ctx, req)
	if err != nil {
		return nil, false, 0, errors.Wrap(err, "could not get validator performance from beacon client")
	}

	for i := range valsPerformance.MissingValidators {
		log.WithFields(log.Fields{
			"Epoch":   a.Epoch,
			"Address": hex.EncodeToString(valsPerformance.MissingValidators[i]),
		}).Warn("Validator performance not found in beacon chain")
	}

	log.Info("Remaining time for next slot: ", ctx)

	return valsPerformance, true, metricsEpoch, nil
}
