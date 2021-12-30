package metrics

import (
	"encoding/hex"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var key1 = ToBytes48([]byte{1})
var key2 = ToBytes48([]byte{2})
var key3 = ToBytes48([]byte{3})

var valsPerformance = &ethpb.ValidatorPerformanceResponse{
	PublicKeys:                    [][]byte{key1[:], key2[:], key3[:]},
	CurrentEffectiveBalances:      []uint64{32 * 1e9, 32 * 1e9, 32 * 1e9},
	InclusionSlots:                []ethTypes.Slot{0, 0, 0},
	InclusionDistances:            []ethTypes.Slot{0, 0, 0},
	CorrectlyVotedSource:          []bool{false, false, true},
	CorrectlyVotedTarget:          []bool{false, false, true},
	CorrectlyVotedHead:            []bool{false, true, true},
	BalancesBeforeEpochTransition: []uint64{100, 110, 130},
	BalancesAfterEpochTransition:  []uint64{90, 120, 150},
	MissingValidators:             [][]byte{},
}

func Test_getBalanceMetrics(t *testing.T) {
	nOfValsWithDecreasedBalance, nOfValidators, earned, losed, decreasedKeys := getBalanceMetrics(valsPerformance)

	require.Equal(t, nOfValsWithDecreasedBalance, uint64(1))
	require.Equal(t, nOfValidators, uint64(3))
	require.Equal(t, earned, big.NewInt(30))
	require.Equal(t, losed, big.NewInt(-10))
	require.Equal(t, losed, big.NewInt(-10))
	require.Equal(t, decreasedKeys, []string{hex.EncodeToString(key1[:])})
}

func Test_getAttestationMetrics(t *testing.T) {
	nOfTotalVotes, nOfIncorrectSource, nOfIncorrectTarget, nOfIncorrectHead, missedKeys := getAttestationMetrics(valsPerformance)

	require.Equal(t, nOfTotalVotes, uint64(9))
	require.Equal(t, nOfIncorrectSource, uint64(2))
	require.Equal(t, nOfIncorrectTarget, uint64(2))
	require.Equal(t, nOfIncorrectHead, uint64(1))
	require.Equal(t, missedKeys, []string{hex.EncodeToString(key1[:])})
}
