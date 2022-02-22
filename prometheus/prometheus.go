package prometheus

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(port int) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}

// TODO: Add the pool before each name

var (
	// TODO: For all validato states, use a GaugeVec with key (state): value (amount)
	NOfUnkownValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_unknown_validators",
			Help:      "Number of unknown validators among all deposited ones",
		},
	)

	NOfDepositedValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_deposited_validators",
			Help:      "Number of deposited validators for the selected from_address/with_cred",
		},
	)

	NOfPendingValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_pending_validators",
			Help:      "Number of pending of activation validators among all deposited ones",
		},
	)

	NOfActiveValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_active_validators",
			Help:      "Number of active validators among all deposited ones",
		},
	)

	NOfExitingValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_exiting_validators",
			Help:      "Number of exiting validators among all deposited ones",
		},
	)

	NOfSlashingValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_slashing_validators",
			Help:      "Number of slashing validators among all deposited ones",
		},
	)

	NOfExitedValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_exited_validators",
			Help:      "Number of exited validators among all deposited ones",
		},
	)

	NOfInvalidValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_invalid_validators",
			Help:      "Number of invalid validators among all deposited ones",
		},
	)

	NOfPartiallyDepositedValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_partiallydeposited_validators",
			Help:      "Number of partially deposited validators among all deposited ones",
		},
	)

	NOfValidatingValidators = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_validating_validators",
			Help:      "Number of validating validators with duties among all deposited ones",
		},
	)

	NOfTotalVotes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_total_votes",
			Help:      "Number of votes for all validators in a given epoch",
		},
	)

	NOfIncorrectSource = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_incorrect_source",
			Help:      "Number of incorrect source votes for all validators in a given epoch",
		},
	)

	NOfIncorrectTarget = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_incorrect_target",
			Help:      "Number of incorrect target votes for all validators in a given epoch",
		},
	)

	NOfIncorrectHead = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_incorrect_head",
			Help:      "Number of incorrect head votes for all validators in a given epoch",
		},
	)

	NOfAttestations = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_attestations",
			Help:      "Number of produced attestations in a given epoch",
		},
	)

	NOfScheduledBlocks = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_scheduled_blocks",
			Help:      "Number of scheduled block proposals in a given epoch",
		},
	)

	NOfProposedBlocks = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "number_proposed_blocks",
			Help:      "Number of proposed blocks in a given epoch",
		},
	)

	AvgIncDistance = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "avg_inc_distance",
			Help:      "Average inclussion distance of all active validators in a given epoch",
		},
	)

	BalanceDecreasedPercent = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "balance_decreased_percent",
			Help:      "Percent of validators that decreased in balance in a given epoch",
		},
	)

	DepositedAmount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "recognized_deposited_amount",
			Help:      "Deposited amount in gwei for the set of validators",
		},
	)

	CumulativeRewards = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "cumulative_rewards",
			Help:      "Cumulative rewards for all validators",
		},
	)

	TotalBalance = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "total_balance_gwei",
			Help:      "Total balance for all validators",
		},
	)

	EffectiveBalance = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "effective_balance_gwei",
			Help:      "Total effective balance for all validators",
		},
	)

	EarnedAmountInEpoch = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "earned_amount_in_epoch",
			Help:      "Earned amount in gwei in the previous epoch transition",
		},
	)

	LosedAmountInEpoch = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "losed_amount_in_epoch",
			Help:      "Losed amount in gwei in the previous epoch transition",
		},
	)

	EthereumPriceUsd = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "eth_price_usd",
			Help:      "Ethereum price in usd",
		},
	)

	MissedAttestationsKeys = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "epoch_missed_attestations_keys",
			Help:      "List of keys and the number of attestations that were missed (since startup)",
		},
		[]string{
			"validatorKey",
		},
	)

	LessBalanceKeys = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "epoch_less_balance_keys",
			Help:      "List of keys and the times its balance decreased (since startup)",
		},
		[]string{
			"validatorKey",
		},
	)

	ProposedBlocks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "proposed_blocks_indexes",
			Help:      "Validator indexes that proposed blocks in a given epoch",
		},
		[]string{
			"epoch",
			"index",
		},
	)

	MissedBlocks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "missed_blocks_indexes",
			Help:      "Validator indexes that missed a block proposal in a given epoch",
		},
		[]string{
			"epoch",
			"index",
		},
	)

	// Code above here will be deprecated

	TotalBalanceMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "total_balance_metrics",
			Help:      "",
		},
		[]string{
			"pool",
		},
	)

	ActiveValidatorsMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "active_validators_metrics",
			Help:      "",
		},
		[]string{
			"pool",
		},
	)

	IncorrectSourceMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "incorrect_source_metrics",
			Help:      "",
		},
		[]string{
			"pool",
		},
	)

	IncorrectTargetMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "incorrect_target_metrics",
			Help:      "",
		},
		[]string{
			"pool",
		},
	)

	IncorrectHeadMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "incorrect_head_metrics",
			Help:      "",
		},
		[]string{
			"pool",
		},
	)

	EpochEarnedAmountMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "epoch_earned_amount_metrics",
			Help:      "",
		},
		[]string{
			"pool",
		},
	)

	EpochLostAmountMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "epoch_lost_amount_metrics",
			Help:      "",
		},
		[]string{
			"pool",
		},
	)

	DeltaEpochBalanceMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "validators",
			Name:      "delta_epoch_balance_metrics",
			Help:      "",
		},
		[]string{
			"pool",
		},
	)

	// TODO: Add remaining time for next slot, to monitor perfomance issues
)
