package metrics

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	//log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var validator_0 = ToBytes48([]byte{10})
var validator_1 = ToBytes48([]byte{20})
var validator_2 = ToBytes48([]byte{30})
var validator_3 = ToBytes48([]byte{40})

func Test_GetIndexesFromKeys(t *testing.T) {
	beaconState := &spec.VersionedBeaconState{
		Altair: &altair.BeaconState{
			Validators: []*phase0.Validator{
				{
					PublicKey: validator_0,
				},
				{
					PublicKey: validator_1,
				},
				{
					PublicKey: validator_2,
				},
				{
					PublicKey: validator_3,
				},
			},
		},
	}

	inputKeys := [][][]byte{
		{validator_3[:], validator_0[:]},                 // test 1
		{validator_0[:]},                                 // test 2
		{validator_3[:], validator_0[:], validator_1[:]}, // test 3
	}

	expectedIndexes := [][]uint64{
		{3, 0},    // test 1
		{0},       // test 2
		{3, 0, 1}, // test 3
	}

	for test := 0; test < len(inputKeys); test++ {
		indexes := GetIndexesFromKeys(
			inputKeys[test],
			beaconState)
		require.Equal(t, indexes, expectedIndexes[test])
	}
}
