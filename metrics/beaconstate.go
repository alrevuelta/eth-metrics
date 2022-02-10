package metrics

import (
	"context"
	"math/big"
	"time"

	//"github.com/alrevuelta/eth-pools-metrics/prometheus"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/rs/zerolog"

	log "github.com/sirupsen/logrus"
)

type BeaconState struct {
	httpClient   *http.Service
	eth2Endpoint string
}

func NewBeaconState(eth2Endpoint string) (*BeaconState, error) {
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
		httpClient:   httpClient,
		eth2Endpoint: eth2Endpoint,
	}, nil
}

func (p *BeaconState) Run() {
	// TODO: Avoid calling the function twice
	var prevBeaconState *spec.VersionedBeaconState = nil
	currentBeaconState, err := p.GetBeaconState()
	if err != nil {
		log.Error("Error fetching beacon state:", err)
	}

	validatorIndexes := []uint64{1, 2, 3}

	todoSetAsFlagUpdateTimeInSeconds := 60 * 60
	for range time.Tick(time.Second * time.Duration(todoSetAsFlagUpdateTimeInSeconds)) {

		source, target, head := GetParticipation(
			validatorIndexes,
			currentBeaconState)

		log.Info("source participation:", source)
		log.Info("target participation:", target)
		log.Info("head participation:", head)

		// TODO: Get validator indexes that missed source

		if prevBeaconState != nil {
			lessBalanceValidators := GetValidatorsWithLessBalance(
				validatorIndexes,
				prevBeaconState,
				currentBeaconState)
			log.Info("validators with less balance", lessBalanceValidators)

			prevBalance := GetTotalBalance(validatorIndexes, prevBeaconState)
			currentBalance := GetTotalBalance(validatorIndexes, currentBeaconState)
			delta := big.NewInt(0).Sub(currentBalance, prevBalance)

			log.Info("prevBalance:", prevBalance)
			log.Info("currentBalance:", currentBalance)
			log.Info("delta balance:", delta)
		}

		prevBeaconState = currentBeaconState

		// TODO: Avoid repeating this
		currentBeaconState, err = p.GetBeaconState()
		if err != nil {
			log.Error("Error fetching beacon state:", err)
		}
	}
}

func (p *BeaconState) GetBeaconState() (*spec.VersionedBeaconState, error) {
	beaconState, err := p.httpClient.BeaconState(context.Background(), "finalized")
	if err != nil {
		return nil, err
	}
	return beaconState, nil
}

func GetTotalBalance(validatorIndexes []uint64, beaconState *spec.VersionedBeaconState) *big.Int {
	totalBalances := big.NewInt(0).SetUint64(0)
	for _, valIdx := range validatorIndexes {
		valBalance := big.NewInt(0).SetUint64(beaconState.Altair.Balances[valIdx])
		totalBalances.Add(totalBalances, valBalance)
	}
	return totalBalances
}

func GetValidatorsWithLessBalance(
	validatorIndexes []uint64,
	prevBeaconState *spec.VersionedBeaconState,
	currentBeaconState *spec.VersionedBeaconState) []uint64 {

	// TODO:

	indexesWithLessBalance := make([]uint64, 0)

	return indexesWithLessBalance

}

func GetParticipation(
	validatorIndexes []uint64,
	beaconState *spec.VersionedBeaconState) (uint64, uint64, uint64) {

	// See spec: from LSB to MSB: source, target, head.
	// https://github.com/ethereum/consensus-specs/blob/master/specs/altair/beacon-chain.md#participation-flag-indices

	var nCorrectSource, nCorrectTarget, nCorrectHead uint64

	for _, valIndx := range validatorIndexes {
		// TODO: Dont know why but Infura returns 0 for all CurrentEpochAttestations

		// TODO: Not working, wait for: https://github.com/attestantio/go-eth2-client/pull/14
		epochAttestations := beaconState.Altair.PreviousEpochAttestations[valIndx]
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
