package metrics

import (
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var key = ToBytes48([]byte{1})

var balances = []*ethpb.ValidatorBalances_Balance{
	{
		PublicKey: key[:], Index: 0, Balance: uint64(33_000_000_000), Status: "ACTIVE",
	},
	{
		PublicKey: key[:], Index: 1, Balance: uint64(32_500_000_000), Status: "ACTIVE",
	},
	{
		PublicKey: key[:], Index: 2, Balance: uint64(34_700_000_000), Status: "ACTIVE",
	},
	{
		PublicKey: key[:], Index: 3, Balance: uint64(35_000_000_001), Status: "ACTIVE",
	},
	{
		PublicKey: key[:], Index: 4, Balance: uint64(31_000_000_000), Status: "EXITED",
	},
}

func Test_getRewardsFromBalances(t *testing.T) {
	cumulativeRewards, totalDeposits := getRewardsFromBalances(balances)

	require.Equal(t, cumulativeRewards, big.NewInt(0).SetUint64(6_200_000_001))
	require.Equal(t, totalDeposits, big.NewInt(0).SetUint64(5*32*1_000_000_000))
}
