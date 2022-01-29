package metrics

import (
	"testing"

	ethTypes "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"
)

var k1 = ToBytes48([]byte{1})
var k2 = ToBytes48([]byte{2})
var k3 = ToBytes48([]byte{3})
var k4 = ToBytes48([]byte{4})
var k5 = ToBytes48([]byte{5})

var valsStatus = &ethpb.MultipleValidatorStatusResponse{
	PublicKeys: [][]byte{k1[:], k2[:], k3[:], k4[:], k5[:]},
	Statuses: []*ethpb.ValidatorStatusResponse{
		{Status: ethpb.ValidatorStatus_UNKNOWN_STATUS},
		{Status: ethpb.ValidatorStatus_EXITING},
		{Status: ethpb.ValidatorStatus_ACTIVE},
		{Status: ethpb.ValidatorStatus_ACTIVE},
		{Status: ethpb.ValidatorStatus_EXITED}},
	Indices: []ethTypes.ValidatorIndex{1, 2, 3, 4, 5},
}

func Test_getValidatorStatusMetrics(t *testing.T) {
	statusMetrics := getValidatorStatusMetrics(valsStatus)

	require.Equal(t, statusMetrics.Validating, uint64(3))
	require.Equal(t, statusMetrics.Unknown, uint64(1))
	require.Equal(t, statusMetrics.Deposited, uint64(0))
	require.Equal(t, statusMetrics.Pending, uint64(0))
	require.Equal(t, statusMetrics.Active, uint64(2))
	require.Equal(t, statusMetrics.Exiting, uint64(1))
	require.Equal(t, statusMetrics.Slashing, uint64(0))
	require.Equal(t, statusMetrics.Exited, uint64(1))
	require.Equal(t, statusMetrics.Invalid, uint64(0))
	require.Equal(t, statusMetrics.PartiallyDeposited, uint64(0))

	// TODO: Test slashed once implemented
}
