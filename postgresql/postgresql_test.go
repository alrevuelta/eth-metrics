package postgresql

import (
	"github.com/alrevuelta/eth-pools-metrics/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func Test_TODO(t *testing.T) {

	// Create mock test

	postgresql := New()
	err := postgresql.StoreProposalDuties()

	validatorPerformance := metrics.ValidatorPerformanceMetrics{}

	validatorPerformance.Epoch = 10
	validatorPerformance.NOfTotalVotes = 0
	validatorPerformance.NOfIncorrectSource = 1
	validatorPerformance.NOfIncorrectTarget = 2
	validatorPerformance.NOfIncorrectHead = 3
	validatorPerformance.NOfValidatingKeys = 4
	validatorPerformance.NOfValsWithLessBalance = 5
	validatorPerformance.EarnedBalance = big.NewInt(10)
	validatorPerformance.LosedBalance = big.NewInt(20)

	err = postgresql.StoreValidatorPerformance(validatorPerformance)

	log.Info("errror:", err)

	require.Equal(t, 1, 1)
	log.Info("done")
}
