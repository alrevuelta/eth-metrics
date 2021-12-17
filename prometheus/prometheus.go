package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func Run(port int) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}

// TODO: Add the pool before each name

var (
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

	// TODO: Use to check extra deposits
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

	// TODO: Add remaining time for next slot, to monitor perfomance issues
)
