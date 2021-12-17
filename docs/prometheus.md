# Prometheus Metrics

The following metrics are exposed by prometheus. Note that in order to access them you must run prometheus connected to your `eth-pools-metrics` instance. Once done, you will be able to access the metrics in a dashboard like for example grafana. All metrics belong to a set of validators, identified by their withdrawal credentials or the address used in the deposit.

* `number_unknown_validators`: Number of validators with unknown status
* `number_deposited_validators`: Number of validators whose deposit was recognized
* `number_pending_validators`: Number of validator pending to be activated
* `number_active_validators`: Number of active validators supposed to fulfil proposal and attestation dutties
* `number_exiting_validators`: Number of validator exiting, about to stop fulfilling duties
* `number_slashing_validators`: Number of validators that were slashed
* `number_exited_validators`: Number of validators that existed the network, no longe fulfilling duties
* `number_invalid_validators`: Number of validators with invalid deposits
* `number_partiallydeposited_validators`: Number of validators that deposited <32 Eth
* `number_validating_validators`: Number of validator that are fulfilling duties. Note they might differ from the active ones.
* `number_total_votes`: Number of votes that the set of validators were supposed to cast (head+source+target) in a given epoch.
* `number_incorrect_source`: Number of incorrect source votes for a given epoch. This vote is the attestation itself.
* `number_incorrect_target`: Number of incorrect target votes for a given epoch
* `number_incorrect_head`: Number of incorrect head votes fo a given epoch
* `number_scheduled_blocks`: Number of scheduled blocks that the set of validators should propose in a given epoch
* `number_proposed_blocks`: Number of proposed blocks that the set of validators actually proposed
* `avg_inc_distance` // TODO: Deprecate. Removed with Altair fork.
* `balance_decreased_percent`: Percent of validators that decreased in balance in a given epoch transition
* `recognized_deposited_amount`: Amount of gwei that all validators deposited. It should be `32Eth*number_deposited_validators`
* `cumulative_rewards`: Cumulative rewards for all validators
