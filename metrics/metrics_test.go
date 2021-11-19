package metrics

import (
	"context"
	"eth-pools-metrics/thegraph" // TODO: Set github prefix when released
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

// Note that these tests require a beacon node running in localhost:4000
// using mainnet.

func Test_Rewards(t *testing.T) {

	// dummy unused address
	withCredList := []string{"00ac6c3bab80be36913be57d5c741aea2d9a44dea1a8b1310b53dec04313e880"}

	metrics, err := NewMetrics(context.Background(), "localhost:4000", "mainnet", withCredList)
	if err != nil {
		log.Fatal(err)
	}

	// random validator
	key0, _ := hexutil.Decode("0xb98648a278ed97ee7adcea09f90da8ce4149efe8d8e3a1799f656fd0312e60eac5acf638bf119417d031802212d38729")
	key1, _ := hexutil.Decode("0x8637abf3d5491267db571ce323ea8d4ab5f9bef9bfb8cf76b48384870e57c094b4524ee366d0f4eda330360da1b8674a")

	metrics.activeKeys = [][]byte{key0, key1}

	epoch := uint64(73614)
	rewards, deposits, err := metrics.GetRewards(context.Background(), epoch)

	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, rewards, big.NewInt(4501492695))
	require.Equal(t, deposits, big.NewInt(64000000000))
}

func Test_Rewards_LargeQuery(t *testing.T) {

	// kraken 1 withdrawal credentials
	withCredList := []string{"00c324cd61bb014032ccdb61cb4f0cf827e654868c7e3858e4e488d70eb43608"}

	metrics, err := NewMetrics(context.Background(), "localhost:4000", "mainnet", withCredList)
	if err != nil {
		log.Fatal(err)
	}

	theGraph, err := thegraph.NewThegraph("mainnet", withCredList)
	require.NoError(t, err)

	pubKeysDeposited, err := theGraph.GetDepositedKeys()
	if err != nil {
		log.Fatal(err)
	}

	metrics.activeKeys = pubKeysDeposited

	epoch := uint64(73614)
	log.Info("Getting rewards")
	rewards, deposits, err := metrics.GetRewards(context.Background(), epoch)

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Rewards: ", rewards)
	log.Info("Deposits: ", deposits)

	require.Equal(t, rewards, big.NewInt(1170102023276))
	require.Equal(t, deposits, big.NewInt(16320000000000))
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

func Test_ParalelGetDuties(t *testing.T) {

	// kraken 1
	withCredList := []string{"004f58172d06b6d54c015d688511ad5656450933aff85dac123cd09410a0825c"}

	// stakefish
	//withCredList := []string{"00ac6c3bab80be36913be57d5c741aea2d9a44dea1a8b1310b53dec04313e880"}

	theGraph, err := thegraph.NewThegraph("mainnet", withCredList)
	require.NoError(t, err)

	pubKeysDeposited, err := theGraph.GetDepositedKeys()

	log.Info("pubKeysDeposited: ", len(pubKeysDeposited))

	metrics, err := NewMetrics(context.Background(), "localhost:4000", "mainnet", withCredList)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Use activeKeys instead
	metrics.depositedKeys = pubKeysDeposited

	chunkSize := 2000
	epoch := uint64(73614)

	duties, err := metrics.ParalelGetDuties(context.Background(), &ethpb.DutiesRequest{
		Epoch:      ethTypes.Epoch(uint64(epoch)),
		PublicKeys: metrics.depositedKeys,
	}, chunkSize)

	if err != nil {
		log.Fatal(err)
	}

	// Duties of that set of validators for the given epoch above
	numOfDuties := 3
	foundDuties := 0

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
