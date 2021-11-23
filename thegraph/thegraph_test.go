package thegraph

import (
	//"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetDepositedKeysByWithCred_Mainnet(t *testing.T) {

	theGraph, err := NewThegraph("mainnet", []string{
		// some of kraken withdrawal credentials
		"004f58172d06b6d54c015d688511ad5656450933aff85dac123cd09410a0825c",
		// some of lido withdrawal credentials
		"010000000000000000000000B9D7934878B5FB9610B3FE8A5E441E8FAD7E293F"},
		// empty wallet addresses
		[]string{})
	require.NoError(t, err)

	pubKeysDeposited, err := theGraph.GetAllDepositedKeys()
	require.NoError(t, err)

	log.Info("Length of deposited keys: ", len(pubKeysDeposited))

	// TODO: Run queries until a given block and assert
	require.Equal(t, 1, 1)
}

func TestGetDepositedKeysByAddress_Mainnet(t *testing.T) {
	theGraph, err := NewThegraph("mainnet",
		[]string{
			"004f58172d06b6d54c015d688511ad5656450933aff85dac123cd09410a0825c",
		},
		[]string{
			// kraken 1 depositor address
			"0xa40dfee99e1c85dc97fdc594b16a460717838703",
			// kraken 2 depositor address
			"0x631c2d8d0d7a80824e602a79800a98d93e909a9e",
		})
	require.NoError(t, err)
	pubKeysDeposited, err := theGraph.GetAllDepositedKeys()
	require.NoError(t, err)

	// TODO: Writte an assert once query up to x block is implemented
	log.Info("Length of deposited keys: ", len(pubKeysDeposited))
}

func TestRemoveDuplicates(t *testing.T) {
	input := make([][]byte, 0)
	input = append(input, []byte("1"))
	input = append(input, []byte("1"))
	input = append(input, []byte("2"))
	input = append(input, []byte("3"))
	input = append(input, []byte("3"))
	input = append(input, []byte("3"))
	input = append(input, []byte("3"))
	input = append(input, []byte("4"))

	clean := RemoveDuplicates(input)

	require.Equal(t, len(clean), 4)
	require.Equal(t, clean[0], input[0])
	require.Equal(t, clean[1], input[2])
	require.Equal(t, clean[2], input[3])
	require.Equal(t, clean[3], input[7])
}
