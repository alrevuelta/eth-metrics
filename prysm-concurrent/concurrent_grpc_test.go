package prysmconcurrent

import (
	"context"
	//"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	//"math/big"
	//"encoding/hex"
	"testing"
)

func Test_ParalelGetDuties(t *testing.T) {
	pubKeysDeposited := GetMainnet27994Keys()

	require.Equal(t, len(pubKeysDeposited), 27994)

	prysmConcurrent, err := NewPrysmConcurrent(context.Background(), "localhost:4000")
	if err != nil {
		log.Fatal(err)
	}

	chunkSize := 2000
	epoch := uint64(73614)

	duties, err := prysmConcurrent.ParalelGetDuties(context.Background(), &ethpb.DutiesRequest{
		Epoch:      ethTypes.Epoch(uint64(epoch)),
		PublicKeys: pubKeysDeposited,
	}, chunkSize)

	if err != nil {
		log.Fatal(err)
	}

	// Duties of that set of validators for the given epoch above
	numOfDuties := 3
	foundDuties := 0

	require.Equal(t, len(duties.CurrentEpochDuties), 27994)

	// Note that the order is not guaranteed
	for i := range duties.CurrentEpochDuties {
		if duties.CurrentEpochDuties[i].ValidatorIndex == 108777 {
			require.Equal(t, int(duties.CurrentEpochDuties[i].ProposerSlots[0]), 2355658)
			foundDuties++
		}
		if duties.CurrentEpochDuties[i].ValidatorIndex == 65806 {
			require.Equal(t, int(duties.CurrentEpochDuties[i].ProposerSlots[0]), 2355677)
			foundDuties++
		}
		if duties.CurrentEpochDuties[i].ValidatorIndex == 68742 {
			require.Equal(t, int(duties.CurrentEpochDuties[i].ProposerSlots[0]), 2355660)
			foundDuties++
		}
	}
	require.Equal(t, numOfDuties, foundDuties)
}

/*
func TestParalelGetMultipleValidatorStatus(t *testing.T) {

	withCredList := []string{"004f58172d06b6d54c015d688511ad5656450933aff85dac123cd09410a0825c"}

	theGraph, err := thegraph.NewThegraph("mainnet", withCredList)
	require.NoError(t, err)

	pubKeysDeposited, err := theGraph.GetDepositedKeys()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("pubKeysDeposited", len(pubKeysDeposited))

	metrics, err := NewMetrics(context.Background(), "localhost:4000", "mainnet", withCredList)
	if err != nil {
		log.Fatal(err)
	}

	metrics.depositedKeys = pubKeysDeposited

	start := time.Now()
	valsStatus, err := metrics.ParalelGetMultipleValidatorStatus(context.Background(), &ethpb.MultipleValidatorStatusRequest{
		PublicKeys: metrics.depositedKeys,
	})
	log.Info("Elapsed time: ", time.Since(start))

	if err != nil {
		log.Fatal("could not get multiple validator status")
	}

	log.Info("valsStatus ", len(valsStatus.PublicKeys))
}
*/

/*
func TestParalelGetValidatorPerformance(t *testing.T) {

	withCredList := []string{"004f58172d06b6d54c015d688511ad5656450933aff85dac123cd09410a0825c"}

	theGraph, err := thegraph.NewThegraph("mainnet", withCredList)
	require.NoError(t, err)

	pubKeysDeposited, err := theGraph.GetDepositedKeys()

	log.Info("pubKeysDeposited", len(pubKeysDeposited))

	metrics, err := NewMetrics(context.Background(), "localhost:4000", "mainnet", withCredList)
	if err != nil {
		log.Fatal(err)
	}

	metrics.depositedKeys = pubKeysDeposited

	// TODO: Use activeKeys instead?

	start := time.Now()
	valsPerformance, err := metrics.ParalelGetValidatorPerformance(context.Background(), &ethpb.ValidatorPerformanceRequest{
		PublicKeys: metrics.depositedKeys,
	})
	log.Info("Elapsed time: ", time.Since(start))

	if err != nil {
		log.Fatal("could not get validator performance", err)
	}

	log.Info("valsPerformance ", len(valsPerformance.PublicKeys))
}
*/
