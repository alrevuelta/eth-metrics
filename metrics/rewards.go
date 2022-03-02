package metrics

import (
	"context"
	"math/big"
	"runtime"
	"time"

	"github.com/alrevuelta/eth-pools-metrics/prometheus"
	"github.com/alrevuelta/eth-pools-metrics/schemas"
	"github.com/pkg/errors"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
)

const gigaWei = uint64(1_000_000_000)
const depositInGigaWei = uint64(32) * gigaWei

func (a *Metrics) StreamRewards() {
	lastEpoch := uint64(0)
	for {
		if a.depositedKeys == nil {
			log.Warn("No active keys to calculate the rewards")
			time.Sleep(30 * time.Second)
			continue
		}
		head, err := GetChainHead(context.Background(), a.beaconChainClient)
		if err != nil {
			log.Error("error getting chain head", err)
			continue
		}

		if uint64(head.FinalizedEpoch) <= lastEpoch {
			time.Sleep(5 * time.Second)
			continue
		}

		log.Info("Getting rewards for epoch: ", head.FinalizedEpoch)

		metrics, err := a.GetRewards(context.Background(), uint64(head.FinalizedEpoch))
		if err != nil {
			log.Error("could not get rewards and balances", err)
			time.Sleep(30 * time.Second)
			continue
		}

		logRewards(metrics)
		setPrometheusRewards(metrics)

		lastEpoch = uint64(head.FinalizedEpoch)

		// Temporal fix to memory leak. Perhaps having an infinite loop
		// inside a routinne is not a good idea. TODO
		runtime.GC()

		// Do not fetch every epoch. For a large number of validators it would be too much
		// TODO: Set as config parameter
		time.Sleep(30 * 60 * time.Second)
	}
}

func logRewards(metrics schemas.RewardsMetrics) {
	log.WithFields(log.Fields{
		"Epoch":             metrics.Epoch,
		"DepositedAmount":   metrics.TotalDeposits.Uint64(),
		"CumulativeRewards": metrics.CumulativeRewards.Uint64(),
	}).Info("Rewards/Balances:")
}

func setPrometheusRewards(metrics schemas.RewardsMetrics) {
	prometheus.DepositedAmount.Set(float64(metrics.TotalDeposits.Uint64()))
	//prometheus.CumulativeRewards.Set(float64(metrics.CumulativeRewards.Uint64()))
}

func getRewardsFromBalances(
	balances []*ethpb.ValidatorBalances_Balance) (*big.Int, *big.Int) {
	cumulativeRewards := big.NewInt(0)
	totalDeposits := big.NewInt(0)

	// Nov 2021: Balances - Deposits matches the rewards
	// but this may change once withdrawals are enabled
	for _, b := range balances {
		status := ethpb.ValidatorStatus(ethpb.ValidatorStatus_value[b.Status])
		if isEligibleForRewards(status) {
			deposit := big.NewInt(0).SetUint64(depositInGigaWei)
			balance := big.NewInt(0).SetUint64(b.Balance)
			reward := big.NewInt(0).Sub(balance, deposit)
			cumulativeRewards.Add(cumulativeRewards, reward)
			totalDeposits.Add(totalDeposits, deposit)
		}
	}
	return cumulativeRewards, totalDeposits
}

func (a *Metrics) GetRewards(ctx context.Context, epoch uint64) (schemas.RewardsMetrics, error) {
	balances, err := a.GetBalances(ctx, epoch)
	if err != nil {
		return schemas.RewardsMetrics{}, errors.Wrap(err, "could not get balances from beacon chain")
	}

	cumulativeRewards, totalDeposits := getRewardsFromBalances(balances)
	rewardsMetrics := schemas.RewardsMetrics{
		Epoch:             epoch,
		TotalDeposits:     totalDeposits,
		CumulativeRewards: cumulativeRewards,
	}

	return rewardsMetrics, nil
}

func (a *Metrics) GetBalances(ctx context.Context, epoch uint64) ([]*ethpb.ValidatorBalances_Balance, error) {
	request := ethpb.ListValidatorBalancesRequest{
		QueryFilter: &ethpb.ListValidatorBalancesRequest_Epoch{ethTypes.Epoch(epoch)},
		PublicKeys:  a.validatingKeys,
	}

	balances := make([]*ethpb.ValidatorBalances_Balance, 0)
	for {
		resp, err := a.beaconChainClient.ListValidatorBalances(ctx, &request)
		if err != nil {
			return nil, errors.Wrap(err, "error getting validator balances")
		}

		balances = append(balances, resp.Balances...)

		if resp.NextPageToken == "" {
			break
		} else {
			request.PageToken = resp.NextPageToken
		}
	}
	return balances, nil
}

// validator statuses that may contain rewards
// Dec 2021: once withdrawals are enabled, this has to be revisited
func isEligibleForRewards(status ethpb.ValidatorStatus) bool {
	if status == ethpb.ValidatorStatus_ACTIVE ||
		status == ethpb.ValidatorStatus_EXITING ||
		status == ethpb.ValidatorStatus_SLASHING ||
		status == ethpb.ValidatorStatus_EXITED {
		return true
	}
	return false
}
