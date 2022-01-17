package postgresql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TODO(t *testing.T) {

	// Create mock test
}

func Test_getDepositsWhereClause(t *testing.T) {
	whereClause := getDepositsWhereClause([]string{"0xkey1", "0xkey2"})
	require.Equal(t,
		"f_eth1_sender = decode('key1', 'hex') or f_eth1_sender = decode('key2', 'hex')",
		whereClause)
}
