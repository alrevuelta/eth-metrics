package metrics

import (
	ethTypes "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	//log "github.com/sirupsen/logrus"
	"github.com/alrevuelta/eth-pools-metrics/schemas"
	"github.com/stretchr/testify/require"
	"testing"
)

// Validators p1-p7 have active duties
var p1 = ToBytes48([]byte{1})
var p2 = ToBytes48([]byte{2})
var p3 = ToBytes48([]byte{3})
var p4 = ToBytes48([]byte{4})
var p5 = ToBytes48([]byte{5})

// Simulate that p6-p7 fail
var p6 = ToBytes48([]byte{6})
var p7 = ToBytes48([]byte{7})

// Assign duties to p1-p7
var duties = &ethpb.DutiesResponse{
	CurrentEpochDuties: []*ethpb.DutiesResponse_Duty{
		{
			ProposerSlots:  []ethTypes.Slot{32000},
			PublicKey:      p1[:],
			ValidatorIndex: 1,
		},
		{
			ProposerSlots:  []ethTypes.Slot{32001},
			PublicKey:      p2[:],
			ValidatorIndex: 2,
		},
		{
			ProposerSlots:  []ethTypes.Slot{32002},
			PublicKey:      p3[:],
			ValidatorIndex: 3,
		},
		{
			ProposerSlots:  []ethTypes.Slot{32003},
			PublicKey:      p4[:],
			ValidatorIndex: 4,
		},
		{
			ProposerSlots:  []ethTypes.Slot{32004},
			PublicKey:      p5[:],
			ValidatorIndex: 5,
		},
		{
			ProposerSlots:  []ethTypes.Slot{32005},
			PublicKey:      p6[:],
			ValidatorIndex: 6,
		},
		{
			ProposerSlots:  []ethTypes.Slot{32006},
			PublicKey:      p7[:],
			ValidatorIndex: 7,
		},
	}}

// Blocks can be:
// BeaconBlockContainer_Phase0Block
// BeaconBlockContainer_AltairBlock
// And soon: BeaconBlockMerge

// Only p1-p5 duties are fulfilled
var blocks = &ethpb.ListBeaconBlocksResponse{
	BlockContainers: []*ethpb.BeaconBlockContainer{
		{
			Block: &ethpb.BeaconBlockContainer_AltairBlock{
				AltairBlock: &ethpb.SignedBeaconBlockAltair{
					Block: &ethpb.BeaconBlockAltair{
						ProposerIndex: 1,
						Slot:          32000,
						Body:          &ethpb.BeaconBlockBodyAltair{Graffiti: []byte("1")}}}}},
		{
			Block: &ethpb.BeaconBlockContainer_AltairBlock{
				AltairBlock: &ethpb.SignedBeaconBlockAltair{
					Block: &ethpb.BeaconBlockAltair{
						ProposerIndex: 2,
						Slot:          32001,
						Body:          &ethpb.BeaconBlockBodyAltair{Graffiti: []byte("2")}}}}},
		{
			Block: &ethpb.BeaconBlockContainer_AltairBlock{
				AltairBlock: &ethpb.SignedBeaconBlockAltair{
					Block: &ethpb.BeaconBlockAltair{
						ProposerIndex: 3,
						Slot:          32002,
						Body:          &ethpb.BeaconBlockBodyAltair{Graffiti: []byte("3")}}}}},
		{
			Block: &ethpb.BeaconBlockContainer_AltairBlock{
				AltairBlock: &ethpb.SignedBeaconBlockAltair{
					Block: &ethpb.BeaconBlockAltair{
						ProposerIndex: 4,
						Slot:          32003,
						Body:          &ethpb.BeaconBlockBodyAltair{Graffiti: []byte("4")}}}}},
		{
			Block: &ethpb.BeaconBlockContainer_AltairBlock{
				AltairBlock: &ethpb.SignedBeaconBlockAltair{
					Block: &ethpb.BeaconBlockAltair{
						ProposerIndex: 5,
						Slot:          32004,
						Body:          &ethpb.BeaconBlockBodyAltair{Graffiti: []byte("5")}}}}},
	},
}

func Test_getProposalDuties(t *testing.T) {
	metrics := getProposalDuties(duties, blocks)

	require.Equal(t, len(metrics.Scheduled), 7)
	require.Equal(t, len(metrics.Proposed), 5)
	require.Equal(t, len(metrics.Missed), 2)

	// Scheduled blocks
	for i := 0; i < 7; i++ {
		require.Equal(t, metrics.Scheduled[i].ValIndex, uint64(i+1))
		require.Equal(t, metrics.Scheduled[i].Slot, ethTypes.Slot(32000+i))
	}

	// Proposed blocks
	for i := 0; i < 5; i++ {
		require.Equal(t, metrics.Proposed[i].ValIndex, uint64(i+1))
		require.Equal(t, metrics.Proposed[i].Slot, ethTypes.Slot(32000+i))
	}

	// Missed blocks
	for i := 0; i < 2; i++ {
		require.Equal(t, metrics.Missed[i].ValIndex, uint64(i+6))
		require.Equal(t, metrics.Missed[i].Slot, ethTypes.Slot(32005+i))
	}
}

func Test_getMissedDuties(t *testing.T) {
	missedDuties := getMissedDuties(
		// Schedulled
		[]schemas.Duty{
			{ValIndex: 1, Slot: ethTypes.Slot(1000)},
			{ValIndex: 2, Slot: ethTypes.Slot(2000)},
			{ValIndex: 3, Slot: ethTypes.Slot(3000)},
			{ValIndex: 4, Slot: ethTypes.Slot(4000)},
		},
		// Proposed
		[]schemas.Duty{
			{ValIndex: 1, Slot: ethTypes.Slot(1000)},
			{ValIndex: 4, Slot: ethTypes.Slot(4000)},
		},
	)

	require.Equal(t, missedDuties[0].ValIndex, uint64(2))
	require.Equal(t, missedDuties[0].Slot, ethTypes.Slot(2000))

	require.Equal(t, missedDuties[1].ValIndex, uint64(3))
	require.Equal(t, missedDuties[1].Slot, ethTypes.Slot(3000))
}
