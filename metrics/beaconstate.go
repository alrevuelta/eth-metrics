package metrics

import (
	"bytes"
	"context"
	"math/big"
	"strconv"
	"time"

	//"github.com/alrevuelta/eth-pools-metrics/prometheus"
	"github.com/alrevuelta/eth-pools-metrics/postgresql"
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
}

func NewBeaconState(eth2Endpoint string, pg *postgresql.Postgresql, fromAddresses []string) (*BeaconState, error) {
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

		currentBeaconState, err := p.GetBeaconState(currentEpoch)
		if err != nil {
			log.Error("Error fetching beacon state:", err)
			continue
		}

		pubKeysDeposited, err := p.pg.GetKeysByFromAddresses(p.fromAddresses)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Info("len of deposited:", len(pubKeysDeposited))

		validatorIndexes := GetIndexesFromKeys(pubKeysDeposited, currentBeaconState)

		log.Info("len indexes:", len(validatorIndexes))

		source, target, head := GetParticipation(
			validatorIndexes,
			currentBeaconState)

		log.Info("source participation:", source)
		log.Info("target participation:", target)
		log.Info("head participation:", head)

		currentBalance, effectiveBalance := GetTotalBalanceAndEffective(validatorIndexes, currentBeaconState)
		log.Info("currentBalance:", currentBalance)
		log.Info("effectiveBalance:", effectiveBalance)
		rewards := big.NewInt(0).Sub(currentBalance, effectiveBalance)
		log.Info("rewards:", rewards)

		// TODO: Get validator indexes that missed source

		if prevBeaconState == nil {
			prevBeaconState = currentBeaconState
			prevEpoch = currentEpoch
			continue
		}

		lessBalanceIndexes, earnedBalance, lostBalance := GetValidatorsWithLessBalance(
			validatorIndexes,
			prevBeaconState,
			currentBeaconState)
		log.Info("validators with less balance", lessBalanceIndexes)
		log.Info("earnedBalance", earnedBalance)
		log.Info("lostBalance", lostBalance)

		prevBalance, _ := GetTotalBalanceAndEffective(validatorIndexes, prevBeaconState)
		delta := big.NewInt(0).Sub(currentBalance, prevBalance)

		log.Info("prevBalance:", prevBalance)
		log.Info("overall delta balance:", delta)

		prevBeaconState = currentBeaconState
		prevEpoch = currentEpoch
	}
}

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
	beaconState *spec.VersionedBeaconState) []uint64 {

	indexes := make([]uint64, 0)

	// TODO: Naive searching approach
	for index, beaconStateKey := range beaconState.Altair.Validators {
		for _, key := range validatorKeys {
			if bytes.Compare(beaconStateKey.PublicKey[:], key) == 0 {
				indexes = append(indexes, uint64(index))
				break
			}
		}
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

func GetParticipation(
	validatorIndexes []uint64,
	beaconState *spec.VersionedBeaconState) (uint64, uint64, uint64) {

	// See spec: from LSB to MSB: source, target, head.
	// https://github.com/ethereum/consensus-specs/blob/master/specs/altair/beacon-chain.md#participation-flag-indices

	var nCorrectSource, nCorrectTarget, nCorrectHead uint64

	for _, valIndx := range validatorIndexes {
		// TODO: Dont know why but Infura returns 0 for all CurrentEpochAttestations

		epochAttestations := beaconState.Altair.PreviousEpochParticipation[valIndx]
		if isBitSet(uint8(epochAttestations), 0) {
			nCorrectSource++
		}
		if isBitSet(uint8(epochAttestations), 1) {
			nCorrectTarget++
		}
		if isBitSet(uint8(epochAttestations), 2) {
			nCorrectHead++
		}
	}
	return nCorrectSource, nCorrectTarget, nCorrectHead
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

func logMetrics(todo string) {
	log.Info("TODO: ", todo)
}

func setPrometheusMetrics() {
	// TODO:
}
