package metrics

/* Commenting until fixed
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
*/
