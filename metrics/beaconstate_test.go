package metrics

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"
)

func Test_TODO(t *testing.T) {
	eth2Endpoint := "localhost:5051"
	bs, err := NewBeaconState(eth2Endpoint)
	if err != nil {
		log.Fatal(err)
	}

	bs.GetBeaconState()
	require.Equal(t, 1, 1)
}
