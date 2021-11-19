package thegraph

import (
	//"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetDepositedKeys_Mainnet(t *testing.T) {

	// kraken withdrawal credentials, contains thousands of deposits
	theGraph, err := NewThegraph("mainnet", []string{"004f58172d06b6d54c015d688511ad5656450933aff85dac123cd09410a0825c"})
	require.NoError(t, err)

	pubKeysDeposited, err := theGraph.GetDepositedKeys()
	require.NoError(t, err)

	log.Info("Length of deposited keys: ", len(pubKeysDeposited))

	// TODO: Run queries until a given block and assert
	require.Equal(t, 1, 1)
}
