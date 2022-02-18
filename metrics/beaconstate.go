package metrics

import (
	"context"
	"encoding/hex"
	"math/big"
	"strconv"
	"time"

	//"github.com/alrevuelta/eth-pools-metrics/prometheus"
	"github.com/alrevuelta/eth-pools-metrics/pools"
	"github.com/alrevuelta/eth-pools-metrics/postgresql"
	"github.com/alrevuelta/eth-pools-metrics/prometheus"
	"github.com/alrevuelta/eth-pools-metrics/schemas"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/rs/zerolog"

	log "github.com/sirupsen/logrus"
)

type BeaconState struct {
	httpClient    *http.Service
	eth2Endpoint  string
	pg            *postgresql.Postgresql
	fromAddresses []string
	poolNames     []string
}

func NewBeaconState(eth2Endpoint string, pg *postgresql.Postgresql, fromAddresses []string, poolNames []string) (*BeaconState, error) {
	client, err := http.New(context.Background(),
		http.WithTimeout(60*time.Second),
		http.WithAddress(eth2Endpoint),
		http.WithLogLevel(zerolog.WarnLevel),
	)
	if err != nil {
		return nil, err
	}

	httpClient := client.(*http.Service)

	return &BeaconState{
		httpClient:    httpClient,
		eth2Endpoint:  eth2Endpoint,
		pg:            pg,
		fromAddresses: fromAddresses,
		poolNames:     poolNames,
	}, nil
}

func (p *BeaconState) Run() {
	var prevEpoch uint64 = uint64(0)
	var prevBeaconState *spec.VersionedBeaconState = nil

	for {
		// Before doing anything, check if we are in the next epoch
		headSlot, err := p.httpClient.NodeSyncing(context.Background())
		if err != nil {
			log.Error("Could not get node sync status:", err)
			continue
		}

		if headSlot.IsSyncing {
			log.Error("Node is not in sync")
			continue
		}
		// TODO: Don't hardcode 32
		// Floor division
		// Go 1 epoch behind head
		currentEpoch := uint64(headSlot.HeadSlot)/uint64(32) - 1

		if prevEpoch >= currentEpoch {
			// do nothing
			continue
		}

		// TODO: Retry once if fails
		currentBeaconState, err := p.GetBeaconState(currentEpoch)
		if err != nil {
			log.Error("Error fetching beacon state:", err)
			continue
		}

		// if no prev beacon state is known, fetch it
		if prevBeaconState == nil {
			prevBeaconState, err = p.GetBeaconState(currentEpoch - 1)
			// TODO: Retry
			if err != nil {
				log.Error(err)
				continue
			}
		}

		for _, poolName := range p.poolNames {
			poolAddressList := pools.PoolsAddresses[poolName]
			log.Info("The pool:", poolName, " keys are: ", poolAddressList)
			pubKeysDeposited, err := p.pg.GetKeysByFromAddresses(poolAddressList)
			if err != nil {
				log.Error(err)
				continue
			}

			valKeyToIndex := PopulateKeysToIndexesMap(currentBeaconState)
			validatorIndexes := GetIndexesFromKeys(pubKeysDeposited, valKeyToIndex)

			metrics := PopulateParticipationAndBalance(
				validatorIndexes,
				currentBeaconState,
				prevBeaconState)

			logMetrics(metrics, poolName)
			setPrometheusMetrics(metrics, poolName)
		}

		prevBeaconState = currentBeaconState
		prevEpoch = currentEpoch
	}
}

func PopulateKeysToIndexesMap(beaconState *spec.VersionedBeaconState) map[string]uint64 {
	// TODO: Naive approach. Reset the map every time
	valKeyToIndex := make(map[string]uint64, 0)
	for index, beaconStateKey := range beaconState.Altair.Validators {
		valKeyToIndex[hex.EncodeToString(beaconStateKey.PublicKey[:])] = uint64(index)
	}
	return valKeyToIndex
}

// TODO: Skip validators that are not active yet
func PopulateParticipationAndBalance(
	validatorIndexes []uint64,
	beaconState *spec.VersionedBeaconState,
	prevBeaconState *spec.VersionedBeaconState) schemas.ValidatorPerformanceMetrics {

	metrics := schemas.ValidatorPerformanceMetrics{
		EarnedBalance:    big.NewInt(0),
		LosedBalance:     big.NewInt(0),
		TotalBalance:     big.NewInt(0),
		EffectiveBalance: big.NewInt(0),
		TotalRewards:     big.NewInt(0),
	}

	nOfIncorrectSource, nOfIncorrectTarget, nOfIncorrectHead, indexesMissedAtt := GetParticipation(
		validatorIndexes,
		beaconState)

	currentBalance, effectiveBalance := GetTotalBalanceAndEffective(validatorIndexes, beaconState)
	rewards := big.NewInt(0).Sub(currentBalance, effectiveBalance)

	lessBalanceIndexes, earnedBalance, lostBalance := GetValidatorsWithLessBalance(
		validatorIndexes,
		prevBeaconState,
		beaconState)

	metrics.IndexesLessBalance = lessBalanceIndexes
	metrics.EarnedBalance = earnedBalance
	metrics.LosedBalance = lostBalance

	// TODO: Don't hardcode 32
	metrics.Epoch = beaconState.Altair.Slot / 32

	metrics.NOfTotalVotes = uint64(len(validatorIndexes)) * 3
	metrics.NOfIncorrectSource = nOfIncorrectSource
	metrics.NOfIncorrectTarget = nOfIncorrectTarget
	metrics.NOfIncorrectHead = nOfIncorrectHead
	metrics.NOfValidatingKeys = uint64(len(validatorIndexes))
	//metrics.NOfValsWithLessBalance = nOfValsWithDecreasedBalance
	//metrics.EarnedBalance = earned
	//metrics.LosedBalance = losed
	metrics.IndexesMissedAtt = indexesMissedAtt
	//metrics.LostBalanceKeys = lostKeys
	metrics.TotalBalance = currentBalance
	metrics.EffectiveBalance = effectiveBalance
	metrics.TotalRewards = rewards

	return metrics
}

// TODO: Get slashed validators

func (p *BeaconState) GetBeaconState(epoch uint64) (*spec.VersionedBeaconState, error) {
	slotStr := strconv.FormatUint(epoch*32, 10)
	beaconState, err := p.httpClient.BeaconState(
		context.Background(),
		slotStr)
	if err != nil {
		return nil, err
	}
	log.Info("Got beacon state for epoch:", beaconState.Altair.Slot/32)
	return beaconState, nil
}

func GetTotalBalanceAndEffective(
	validatorIndexes []uint64,
	beaconState *spec.VersionedBeaconState) (*big.Int, *big.Int) {
	totalBalances := big.NewInt(0).SetUint64(0)
	effectiveBalance := big.NewInt(0).SetUint64(0)
	for _, valIdx := range validatorIndexes {
		valBalance := big.NewInt(0).SetUint64(beaconState.Altair.Balances[valIdx])
		valEffBalance := big.NewInt(0).SetUint64(uint64(beaconState.Altair.Validators[valIdx].EffectiveBalance))
		totalBalances.Add(totalBalances, valBalance)
		effectiveBalance.Add(effectiveBalance, valEffBalance)
	}
	return totalBalances, effectiveBalance
}

func GetIndexesFromKeys(
	validatorKeys [][]byte,
	valKeyToIndex map[string]uint64) []uint64 {

	indexes := make([]uint64, 0)

	// Use global prepopulated map
	for _, key := range validatorKeys {
		indexes = append(indexes, valKeyToIndex[hex.EncodeToString(key)])
	}

	return indexes
}

func GetValidatorsWithLessBalance(
	validatorIndexes []uint64,
	prevBeaconState *spec.VersionedBeaconState,
	currentBeaconState *spec.VersionedBeaconState) ([]uint64, *big.Int, *big.Int) {

	indexesWithLessBalance := make([]uint64, 0)
	earnedBalance := big.NewInt(0)
	lostBalance := big.NewInt(0)

	for _, valIdx := range validatorIndexes {
		// handle if there was a new validator index not register in the prev state
		if valIdx >= uint64(len(prevBeaconState.Altair.Balances)) {
			continue
		}

		prevEpochValBalance := big.NewInt(0).SetUint64(prevBeaconState.Altair.Balances[valIdx])
		currentEpochValBalance := big.NewInt(0).SetUint64(currentBeaconState.Altair.Balances[valIdx])
		delta := big.NewInt(0).Sub(currentEpochValBalance, prevEpochValBalance)

		if delta.Cmp(big.NewInt(0)) == -1 {
			indexesWithLessBalance = append(indexesWithLessBalance, valIdx)
			lostBalance.Add(lostBalance, delta)
		} else {
			earnedBalance.Add(earnedBalance, delta)
		}
	}

	return indexesWithLessBalance, earnedBalance, lostBalance
}

// See spec: from LSB to MSB: source, target, head.
// https://github.com/ethereum/consensus-specs/blob/master/specs/altair/beacon-chain.md#participation-flag-indices
func GetParticipation(
	validatorIndexes []uint64,
	beaconState *spec.VersionedBeaconState) (uint64, uint64, uint64, []uint64) {

	indexesMissedAtt := make([]uint64, 0)

	var nIncorrectSource, nIncorrectTarget, nIncorrectHead uint64

	for _, valIndx := range validatorIndexes {
		// Ignore slashed validators
		if beaconState.Altair.Validators[valIndx].Slashed {
			continue
		}
		beaconStateEpoch := beaconState.Altair.Slot / 32
		// Ignore not yet active validators
		// TODO: Test this
		if uint64(beaconState.Altair.Validators[valIndx].ActivationEpoch) > beaconStateEpoch {
			//log.Warn("index: ", valIndx, " is not active yet")
			continue
		}

		// TODO: Dont know why but Infura returns 0 for all CurrentEpochAttestations
		epochAttestations := beaconState.Altair.PreviousEpochParticipation[valIndx]
		if !isBitSet(uint8(epochAttestations), 0) {
			nIncorrectSource++
			indexesMissedAtt = append(indexesMissedAtt, valIndx)
		}
		if !isBitSet(uint8(epochAttestations), 1) {
			nIncorrectTarget++
		}
		if !isBitSet(uint8(epochAttestations), 2) {
			nIncorrectHead++
		}
	}
	return nIncorrectSource, nIncorrectTarget, nIncorrectHead, indexesMissedAtt
}

func GetInactivityScores(
	validatorIndexes []uint64,
	beaconState *spec.VersionedBeaconState) []uint64 {
	inactivityScores := make([]uint64, 0)
	for _, valIdx := range validatorIndexes {
		inactivityScores = append(inactivityScores, beaconState.Altair.InactivityScores[valIdx])
	}
	return inactivityScores
}

// Check if bit n (0..7) is set where 0 is the LSB in little endian
func isBitSet(input uint8, n int) bool {
	return (input & (1 << n)) > uint8(0)
}

func logMetrics(
	metrics schemas.ValidatorPerformanceMetrics,
	poolName string) {
	balanceDecreasedPercent := (float64(len(metrics.IndexesLessBalance)) / float64(metrics.NOfValidatingKeys)) * 100

	log.WithFields(log.Fields{
		"PoolName":                    poolName,
		"Epoch":                       metrics.Epoch,
		"nOfTotalVotes":               metrics.NOfTotalVotes,
		"nOfIncorrectSource":          metrics.NOfIncorrectSource,
		"nOfIncorrectTarget":          metrics.NOfIncorrectTarget,
		"nOfIncorrectHead":            metrics.NOfIncorrectHead,
		"nOfValidators":               metrics.NOfValidatingKeys,
		"PercentIncorrectSource":      (float64(metrics.NOfIncorrectSource) / float64(metrics.NOfTotalVotes)) * 100,
		"PercentIncorrectTarget":      (float64(metrics.NOfIncorrectTarget) / float64(metrics.NOfTotalVotes)) * 100,
		"PercentIncorrectHead":        (float64(metrics.NOfIncorrectHead) / float64(metrics.NOfTotalVotes)) * 100,
		"nOfValsWithDecreasedBalance": len(metrics.IndexesLessBalance),
		"balanceDecreasedPercent":     balanceDecreasedPercent,
		"epochEarnedBalance":          metrics.EarnedBalance,
		"epochLostBalance":            metrics.LosedBalance,
		"totalBalance":                metrics.TotalBalance,
		"effectiveBalance":            metrics.EffectiveBalance,
		"totalRewards":                metrics.TotalRewards,
		"ValidadorKeyMissedAtt":       metrics.IndexesMissedAtt,
		"ValidadorKeyLessBalance":     metrics.IndexesLessBalance,
	}).Info(poolName + " Stats:")
}

func setPrometheusMetrics(
	metrics schemas.ValidatorPerformanceMetrics,
	poolName string) {

	prometheus.TotalBalanceMetrics.WithLabelValues(
		poolName).Set(float64(metrics.TotalBalance.Int64()))

	prometheus.ActiveValidatorsMetrics.WithLabelValues(
		poolName).Set(float64(metrics.NOfValidatingKeys))

	prometheus.IncorrectSourceMetrics.WithLabelValues(
		poolName).Set(float64(metrics.NOfIncorrectSource))

	prometheus.IncorrectTargetMetrics.WithLabelValues(
		poolName).Set(float64(metrics.NOfIncorrectTarget))

	prometheus.IncorrectHeadMetrics.WithLabelValues(
		poolName).Set(float64(metrics.NOfIncorrectHead))

	prometheus.NOfTotalVotes.Set(float64(metrics.NOfTotalVotes))
	prometheus.NOfIncorrectSource.Set(float64(metrics.NOfIncorrectSource))
	prometheus.NOfIncorrectTarget.Set(float64(metrics.NOfIncorrectTarget))
	prometheus.NOfIncorrectHead.Set(float64(metrics.NOfIncorrectHead))
	prometheus.EarnedAmountInEpoch.Set(float64(metrics.EarnedBalance.Int64()))
	prometheus.LosedAmountInEpoch.Set(float64(metrics.LosedBalance.Int64()))
	prometheus.CumulativeRewards.Set(float64(metrics.TotalRewards.Int64()))
	prometheus.TotalBalance.Set(float64(metrics.TotalBalance.Int64()))
	prometheus.EffectiveBalance.Set(float64(metrics.EffectiveBalance.Int64()))

	// TODO: Deprecate this, send the raw number
	balanceDecreasedPercent := (float64(metrics.NOfValsWithLessBalance) / float64(metrics.NOfValidatingKeys)) * 100
	prometheus.BalanceDecreasedPercent.Set(balanceDecreasedPercent)

	/* TODO:
	for _, v := range metrics.MissedAttestationsKeys {
		prometheus.MissedAttestationsKeys.WithLabelValues(v).Inc()
	}

	for _, v := range metrics.LostBalanceKeys {
		prometheus.LessBalanceKeys.WithLabelValues(v).Inc()
	}*/
}
