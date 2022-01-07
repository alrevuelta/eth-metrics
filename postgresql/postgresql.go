package postgresql

import (
	"context"
	"github.com/alrevuelta/eth-pools-metrics/schemas"
	"github.com/jackc/pgx/v4"
	//log "github.com/sirupsen/logrus"
	"strings"
)

var poolNamePlaceholder = "POOLNAMEPLACEHOLDER"

// If a new field is added, the table has to be manually reset
var createPool = `
CREATE TABLE IF NOT EXISTS t_POOLNAMEPLACEHOLDER (
	 f_pool TEXT,
	 f_epoch BIGINT PRIMARY KEY,

	 f_n_total_votes BIGINT,
	 f_n_incorrect_source BIGINT,
	 f_n_incorrect_target BIGINT,
	 f_n_incorrect_head BIGINT,
	 f_n_validating_keys BIGINT,
	 f_n_valitadors_with_less_balace BIGINT,
	 f_epoch_earned_balance BIGINT,
	 f_epoch_lost_balace BIGINT,

	 f_eth_price_usd BIGINT,

	 f_n_scheduled_blocks BIGINT,
	 f_n_proposed_blocks BIGINT
);
`

// TODO: Add missing
// MissedAttestationsKeys []string
// LostBalanceKeys        []string
var insertValidatorPerformance = `
INSERT INTO t_POOLNAMEPLACEHOLDER(
	f_epoch,
	f_n_total_votes,
	f_n_incorrect_source,
	f_n_incorrect_target,
	f_n_incorrect_head,
	f_n_validating_keys,
	f_n_valitadors_with_less_balace,
	f_epoch_earned_balance,
	f_epoch_lost_balace)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (f_epoch)
DO UPDATE SET
   f_n_total_votes=EXCLUDED.f_n_total_votes,
	 f_n_incorrect_source=EXCLUDED.f_n_incorrect_source,
	 f_n_incorrect_target=EXCLUDED.f_n_incorrect_target,
	 f_n_incorrect_head=EXCLUDED.f_n_incorrect_head,
	 f_n_validating_keys=EXCLUDED.f_n_validating_keys,
	 f_n_valitadors_with_less_balace=EXCLUDED.f_n_valitadors_with_less_balace,
	 f_epoch_earned_balance=EXCLUDED.f_epoch_earned_balance,
	 f_epoch_lost_balace=EXCLUDED.f_epoch_lost_balace
`


var insertProposalDuties = `
INSERT INTO t_POOLNAMEPLACEHOLDER(
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
		a.GetQuery(createPool)); err != nil {
		return err
	}
	return nil
}

func (a *Postgresql) GetQuery(query string) string {
	return strings.Replace(query, poolNamePlaceholder, a.PoolName, -1)
}

func (a *Postgresql) StoreProposalDuties(epoch uint64, scheduledBlocks uint64, proposedBlocks uint64) error {
	_, err := a.postgresql.Exec(
		context.Background(),
		a.GetQuery(insertProposalDuties),
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
		a.GetQuery(insertValidatorPerformance),
		validatorPerformance.Epoch,
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
