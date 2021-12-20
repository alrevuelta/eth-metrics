# eth-pools-metrics

[![Tag](https://img.shields.io/github/tag/alrevuelta/eth-pools-metrics.svg)](https://github.com/alrevuelta/eth-pools-metrics/releases/)
[![Release](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/release.yml/badge.svg)](https://github.com/alrevuelta/eth-pools-metrics/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alrevuelta/eth-pools-metrics)](https://goreportcard.com/report/github.com/alrevuelta/eth-pools-metrics)

## Introduction

Monitor the performance of your Ethereum 2.0 staking pool. Just input the withdrawal credential(s) or wallet address(es) that was used in the deposit contract and the network you want to run. Note that a prysm gRPC beacon-chain is required at `localhost:4000`. Tested with up to 25000 validators. Some features:
* Deposited Eth and rewards monitoring
* Rates of faulty head/source/target voting
* Monitor the percent of validators which balance decreased
* Proposed and missed blocks monitoring
* All metrics are exposed with prometheus, see `/prometheus`
* Only the latest epoch is analyzed, this tool can't go back in time from deployment
* No need to run an archival node, default config should be enough

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
