package postgresql

import (
	"context"
	"github.com/alrevuelta/eth-pools-metrics/schemas"
	"github.com/jackc/pgx/v4"
	//log "github.com/sirupsen/logrus"
)

// If a new field is added, the table has to be manually reset
var createPoolsMetricsTable = `
CREATE TABLE IF NOT EXISTS t_pools_metrics_summary (
	 f_epoch BIGINT,
	 f_pool TEXT,
	 f_epoch_timestamp TIMESTAMPTZ NOT NULL,

	 f_n_total_votes BIGINT,
	 f_n_incorrect_source BIGINT,
	 f_n_incorrect_target BIGINT,
	 f_n_incorrect_head BIGINT,
	 f_n_validating_keys BIGINT,
	 f_n_valitadors_with_less_balace BIGINT,
	 f_epoch_earned_balance BIGINT,
	 f_epoch_lost_balace BIGINT,

	 f_n_scheduled_blocks BIGINT,
	 f_n_proposed_blocks BIGINT,

	 PRIMARY KEY (f_epoch, f_pool)
);
`

// TODO: Store price
//f_eth_price_usd BIGINT,

// TODO: Add missing
// MissedAttestationsKeys []string
// LostBalanceKeys        []string
var insertValidatorPerformance = `
INSERT INTO t_pools_metrics_summary(
	f_epoch,
	f_pool,
	f_epoch_timestamp,
	f_n_total_votes,
	f_n_incorrect_source,
	f_n_incorrect_target,
	f_n_incorrect_head,
	f_n_validating_keys,
	f_n_valitadors_with_less_balace,
	f_epoch_earned_balance,
	f_epoch_lost_balace)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (f_epoch, f_pool)
DO UPDATE SET
   f_epoch_timestamp=EXCLUDED.f_epoch_timestamp,
   f_n_total_votes=EXCLUDED.f_n_total_votes,
	 f_n_incorrect_source=EXCLUDED.f_n_incorrect_source,
	 f_n_incorrect_target=EXCLUDED.f_n_incorrect_target,
	 f_n_incorrect_head=EXCLUDED.f_n_incorrect_head,
	 f_n_validating_keys=EXCLUDED.f_n_validating_keys,
	 f_n_valitadors_with_less_balace=EXCLUDED.f_n_valitadors_with_less_balace,
	 f_epoch_earned_balance=EXCLUDED.f_epoch_earned_balance,
	 f_epoch_lost_balace=EXCLUDED.f_epoch_lost_balace
`

// TODO: Add f_epoch_timestamp
var insertProposalDuties = `
INSERT INTO t_pools_metrics_summary(
	f_epoch,
	f_n_scheduled_blocks,
	f_n_proposed_blocks)
VALUES ($1, $2, $3)
ON CONFLICT (f_epoch)
DO UPDATE SET
   f_n_scheduled_blocks=EXCLUDED.f_n_scheduled_blocks,
	 f_n_proposed_blocks=EXCLUDED.f_n_proposed_blocks
`

type Postgresql struct {
	postgresql *pgx.Conn
	PoolName   string
}

// postgresql://user:password@netloc:port/dbname
func New(postgresEndpoint string, poolName string) (*Postgresql, error) {
	conn, err := pgx.Connect(context.Background(), postgresEndpoint)

	if err != nil {
		return nil, err
	}

	return &Postgresql{
		postgresql: conn,
		PoolName:   poolName,
	}, nil
}

func (a *Postgresql) CreateTable() error {
	if _, err := a.postgresql.Exec(
		context.Background(),
		createPoolsMetricsTable); err != nil {
		return err
	}
	return nil
}

func (a *Postgresql) StoreProposalDuties(epoch uint64, scheduledBlocks uint64, proposedBlocks uint64) error {
	_, err := a.postgresql.Exec(
		context.Background(),
		insertProposalDuties,
		epoch,
		scheduledBlocks,
		proposedBlocks)

	if err != nil {
		return err
	}
	return nil
}

func (a *Postgresql) StoreValidatorPerformance(validatorPerformance schemas.ValidatorPerformanceMetrics) error {
	_, err := a.postgresql.Exec(
		context.Background(),
		insertValidatorPerformance,
		validatorPerformance.Epoch,
		validatorPerformance.PoolName,
		validatorPerformance.Time,
		validatorPerformance.NOfTotalVotes,
		validatorPerformance.NOfIncorrectSource,
		validatorPerformance.NOfIncorrectTarget,
		validatorPerformance.NOfIncorrectHead,
		validatorPerformance.NOfValidatingKeys,
		validatorPerformance.NOfValsWithLessBalance,
		validatorPerformance.EarnedBalance.Int64(),
		validatorPerformance.LosedBalance.Int64())

	if err != nil {
		return err
	}
	return nil
}
