# eth-pools-metrics

[![Tag](https://img.shields.io/github/tag/alrevuelta/eth-pools-metrics.svg)](https://github.com/alrevuelta/eth-pools-metrics/releases/)
[![Release](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/release.yml/badge.svg)](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alrevuelta/eth-pools-metrics)](https://goreportcard.com/report/github.com/alrevuelta/eth-pools-metrics)
[![Tests](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/tests.yml/badge.svg)](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/tests.yml)

## Introduction

Monitor the performance of your ethereum consensus staking pool. Just input the withdrawal credential(s) or wallet address(es) that was used in the deposit contract and the network you want to run. This will be used to identify your validators. Some of the parameters that are monitored:
* Deposited Eth and rewards
* Rates of faulty head/source/target votes (see GASPER algorithm)
* Delta in rewards/penalties between consecutive epochs
* Proposed and missed blocks for each epoch

Some features: 
* All metrics are exposed with prometheus, see `/prometheus`
* Calculates all metrics streaming the latest head-1 epoch
* No need to run an archival node, default config should be enough

See [this](https://github.com/alrevuelta/eth-pools-metrics/blob/master/prometheus/prometheus.go) for more information about the metrics and [this](https://github.com/alrevuelta/eth-pools-metrics/blob/master/docs/pools.md) if you want to get your pool monitored.

**A note to old users:** This project started using [prysm](https://github.com/prysmaticlabs/prysm) gRPC but has migrated to the http api to be cross compatible with all clients. If you are still interested in the gRPC implementation, see [v0.0.10](https://github.com/alrevuelta/eth-metrics/releases/tag/v0.0.10) release.


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

This project requires:
* An ethereum `consensus` client compliant with the http api
* An ethereum `execution` client compliant with the http api
* `chaind` instance indexing deposits
* `prometheus` (optional)

### consensus-client

### execution-client

### chaind

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

### prometheus (optional)

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
  -eth1address string
    	Ethereum 1 http endpoint. To be used by rocket pool
  -eth2address string
    	Ethereum 2 http endpoint
  -from-address value
    	Wallet addresses used to deposit. Can be used multiple times
  -pool-name value
    	Pool name to monitor. Can be useed multiple times
  -postgres string
    	Postgres db endpoit: postgresql://user:password@netloc:port/dbname (optional)
  -prometheus-port int
    	Prometheus port to listen to (default 9500)
  -version
    	Prints the release version and exits
  -withdrawal-credentials value
    	Withdrawal credentials used in the deposits. Can be used multiple times
```

## Example

Log metrics for kraken pool (see `pools/pools.go`)

```console
$ ./eth-pools-metrics \
--postgres=xxx://yyy:kkk@localhost:port/zzz \
--eth1address=https:your-execution-endpoint \
--eth2address=https:your-consensus-endpoint \
--pool-name=kraken
```

You can provide your own hardcoded set of validators by providing a file containing the public keys of the validators, one key per line.
```console
$ ./eth-pools-metrics \
--postgres=xxx://yyy:kkk@localhost:port/zzz \
--eth1address=https:your-execution-endpoint \
--eth2address=https:your-consensus-endpoint \
--pool-name=/validators/coinbase.txt
```
