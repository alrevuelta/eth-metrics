# eth-pools-metrics


## Introduction

Monitor the performance of your Ethereum 2.0 staking pool. Just input the withdrawal credentials that were used in the deposit contract and the network you want to run. Note that a prysm gRPC beacon-chain is required at `localhost:4000`. Tested with up to 25000 validators. Some features:
* Deposited Eth and rewards monitoring
* Rates of faulty head/source/target voting
* Monitor the percent of validators which balance decreased
* Proposed and missed blocks monitoring
* All metrics are exposed with prometheus, see `/prometheus`
* Only the latest epoch is analyzed, this tool can't go back in time
* No need to run an archival node, default config should be enough

## Build

### Docker

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

```console
$ ./eth-pools-metrics --help
Usage of ./eth-pools-metrics:
  -beacon-rpc-endpoint string
    	Address:Port of a eth2 beacon node endpoint (default "localhost:4000")
  -from-address value
    	Wallet addresses used to deposit. Can be used multiple times
  -network string
    	Ethereum 2.0 network mainnet|prater|pyrmont (default "mainnet")
  -prometheus-port int
    	Prometheus port to listen to (default 9500)
  -withdrawal-credentials value
    	Withdrawal credentials used in the deposits. Can be used multiple times
```

## Example

Note that a prysm beacon-node must be running in `localhost:4000`. Set `from-address` to the address used to make the deposits to the eth2.0 contract.
```console
$ ./eth-pools-metrics \
--network=mainnet \
--from-address=0x631c2d8d0d7a80824e602a79800a98d93e909a9e
```

## TODO:
* Experiment with sync committee duties
* Flag worst performing validators
