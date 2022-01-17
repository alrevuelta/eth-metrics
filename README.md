# eth-pools-metrics

[![Tag](https://img.shields.io/github/tag/alrevuelta/eth-pools-metrics.svg)](https://github.com/alrevuelta/eth-pools-metrics/releases/)
[![Release](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/release.yml/badge.svg)](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alrevuelta/eth-pools-metrics)](https://goreportcard.com/report/github.com/alrevuelta/eth-pools-metrics)
[![Tests](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/tests.yml/badge.svg)](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/tests.yml)

## Introduction

Monitor the performance of your Ethereum 2.0 staking pool. Just input the withdrawal credential(s) or wallet address(es) that was used in the deposit contract and the network you want to run. Note that a prysm gRPC beacon-chain is required at `localhost:4000`. Tested with up to 30000 validators. Some features:
* Deposited Eth and rewards monitoring
* Rates of faulty head/source/target voting
* Monitor the percent of validators which balance decreased
* Monitor the amount of eth that was earned/lost in an epoch
* Proposed and missed blocks monitoring
* All metrics are exposed with prometheus, see `/prometheus`
* Only the latest epoch is analyzed, this tool can't go back in time from deployment
* No need to run an archival node, default config should be enough

See [this](https://github.com/alrevuelta/eth-pools-metrics/blob/master/prometheus/prometheus.go) for more information about the metrics and [this](https://github.com/alrevuelta/eth-pools-metrics/blob/master/docs/pools.md) if you want to get your pool monitored.

## Build

### Docker

Note that the docker image is publicly available and can be fetched as follows:

```console
docker pull alrevuelta/eth-pools-metrics:latest
```

Build with docker:

```console
git clone https://github.com/alrevuelta/eth-pools-metrics.git
docker build -t eth-pools-metrics .
```

### Source

```console
git clone https://github.com/alrevuelta/eth-pools-metrics.git
go build
```

## Requirements

Depending on how you want to use `eth-metrics`, you may need to run some extra software:
* Ethereum 2.0 deposits can be fetched in two ways: from [thegraph](https://thegraph.com/) or from a local database. Depending on your use case, you may want to use one or the other. For large amounts of validators, perhaps using a local database with pre-indexed deposits is the best approach, since thegraph has some api call limits.
* If you want to access the metrics, you may want to deploy `prometheus`.

### Custom deposits database

If you opt for running your own deposits indexer instead of just relying on thegraph, we recommend `chaind` project. Assuming you already have a postgres database running, you can run chaind as follows. This will create a `t_eth1_deposits` table that will be populated with all deposits to the deposits smart contract. Note that this table can take few hours to sync.

```console
./chaind \
--blocks.enable=false \
--finalizer.enable=false \
--summarizer.enable=false \
--summarizer.epochs.enable=false \
--summarizer.blocks.enable=false \
--summarizer.validators.enable=false \
--validators.enable=false \
--validators.balances.enable=false \
--beacon-committees.enable=false \
--proposer-duties.enable=false \
--sync-committees.enable=false \
--eth1deposits.enable=true \
--eth1deposits.start-block=11185311 \
--eth1client.address=https://your-eth1-endpoint \
--chaindb.url=postgresql://postgres:password@url:5432/user \
--eth2client.address="http://your-eth2-endpoint" \
--log-level=trace
```

### Prometheus

You can expose the metrics to access them in a dashboard like grafana by running `prometheus` with [this configuration](https://github.com/alrevuelta/eth-pools-metrics/blob/master/docs/prometheus.yml).

Some notes:
* Block 11185311 is used as starting block for mainnet, that was when the first deposit was registered.
* `chaind` can index other data, but only eth1 deposits is enabled.

## Usage

The following flags are available:

```console
$ ./eth-pools-metrics --help
Usage of ./eth-pools-metrics:
  -beacon-rpc-endpoint string
    	Address:Port of a eth2 beacon node endpoint (default "localhost:4000")
  -from-address value
    	Wallet addresses used to deposit. Can be used multiple times
  -network string
    	Ethereum 2.0 network mainnet|prater|pyrmont (default "mainnet")
  -pool-name string
    	Name of the pool being monitored. If known, addresses are loaded by default (see known pools) (default "required")
  -prometheus-port int
    	Prometheus port to listen to (default 9500)
  -version
    	Prints the release version and exits
  -withdrawal-credentials value
    	Withdrawal credentials used in the deposits. Can be used multiple times
```

## Example

Note that a prysm beacon-node must be running in `beacon-rpc-endpoint`. Set `from-address` to the address used to make the deposits to the eth2.0 contract. This addess(es) will be used to identify the validators to monitor. Metrics are shown as logs, but they are also exposed to prometheus.

```console
$ ./eth-pools-metrics \
--network=mainnet \
--from-address=0x631c2d8d0d7a80824e602a79800a98d93e909a9e
--beacon-rpc-endpoint=localhost:4000
--pool-name=my-pool-name
```

If the `pool-name` is known (see `pools/pools.go`) there is no need to provide `from-address`. Some of the major exchanges and pools addresses are already configured, so that you don't have to input all of them manually. For example, one can monitor all `poloniex` validators with the following command.

```console
$ ./eth-pools-metrics \
--network=mainnet \
--beacon-rpc-endpoint=localhost:4000
--pool-name=poloniex
```

## TODO:
* Experiment with sync committee duties
* Flag worst performing validators
