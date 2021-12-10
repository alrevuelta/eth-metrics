package metrics

import (
	"context"
	"eth-pools-metrics/prometheus" // TODO: Set github prefix when released
	"github.com/pkg/errors"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

func (a *Metrics) StreamRewards() {
	go func() {
		lastEpoch := uint64(0)
		for {
			if a.activeKeys == nil {
				log.Warn("No active keys to calculate the rewards")
				time.Sleep(30 * time.Second)
				continue
			}
			head, err := GetChainHead(context.Background(), a.beaconChainClient)
			if err != nil {
				log.Error("error getting chain head", err)
			}

			if uint64(head.FinalizedEpoch) <= lastEpoch {
				time.Sleep(5 * time.Second)
				continue
			}

			log.Info("Getting rewards for epoch: ", head.FinalizedEpoch)

			cumulativeRewards, depositedAmount, err := a.GetRewards(context.Background(), uint64(head.FinalizedEpoch))
			if err != nil {
				log.Error("could not get rewards and balances", err)
				time.Sleep(30 * time.Second)
				continue
			}

			prometheus.DepositedAmount.Set(float64(depositedAmount.Uint64()))
			prometheus.CumulativeRewards.Set(float64(cumulativeRewards.Uint64()))

			log.WithFields(log.Fields{
				"Epoch":             uint64(head.FinalizedEpoch),
				"DepositedAmount":   depositedAmount.Uint64(),
				"CumulativeRewards": cumulativeRewards.Uint64(),
			}).Info("Rewards/Balances:")
			lastEpoch = uint64(head.FinalizedEpoch)

			// Do not fetch every epoch. For a large number of validators it would be too much
			// TODO: Set as config parameter
			time.Sleep(30 * 60 * time.Second)
		}
	}()
}

func (a *Metrics) GetRewards(ctx context.Context, epoch uint64) (*big.Int, *big.Int, error) {
	// Nov 2021: Balances - Deposits matches the rewards
	// but this may change once withdrawals are enabled
	balances, err := a.GetBalances(ctx, epoch)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get balances from beacon chain")
	}

	totalRewards := big.NewInt(0)
	totalDeposits := big.NewInt(0)
	for _, b := range balances {
		status := ethpb.ValidatorStatus(ethpb.ValidatorStatus_value[b.Status])
		if status == ethpb.ValidatorStatus_ACTIVE ||
			status == ethpb.ValidatorStatus_EXITING ||
			status == ethpb.ValidatorStatus_SLASHING {
			deposit := big.NewInt(0).SetUint64(depositInGigaWei)
			balance := big.NewInt(0).SetUint64(b.Balance)
			reward := big.NewInt(0).Sub(balance, deposit)
			totalRewards.Add(totalRewards, reward)
			totalDeposits.Add(totalDeposits, deposit)
		}
	}
	return totalRewards, totalDeposits, nil
}

func (a *Metrics) GetBalances(ctx context.Context, epoch uint64) ([]*ethpb.ValidatorBalances_Balance, error) {
	request := ethpb.ListValidatorBalancesRequest{
		QueryFilter: &ethpb.ListValidatorBalancesRequest_Epoch{ethTypes.Epoch(epoch)},
		PublicKeys:  a.activeKeys,
	}

	balances := make([]*ethpb.ValidatorBalances_Balance, 0)
	for {
		resp, err := a.beaconChainClient.ListValidatorBalances(ctx, &request)
		if err != nil {
			return nil, errors.Wrap(err, "error getting validator balances")
		}

		balances = append(balances, resp.Balances...)

    // TODO: Add debug traces
    //log.Info("NextPageToken: ", resp.NextPageToken, " length: ", len(resp.Balances))

		if resp.NextPageToken == "" {
			break
		} else {
			request.PageToken = resp.NextPageToken
		}
	}
	return balances, nil
}
