# eth-pools-metrics

Monitor the performance of your Ethereum 2.0 staking pool. Just input the withdrawal credentials that were used in the deposit contract and the network you want to run. Note that a prysm gRPC beacon-chain is required at `localhost:4000`. Tested with up to 25000 validators. Some features:
* Deposited Eth and rewards monitoring
* Rates of faulty head/source/target voting
* Monitor the percent of validators which balance decreased
* Proposed and missed blocks monitoring
* All metrics are exposed with prometheus, see `/prometheus`

## Usage:

```console
$ ./eth-pools-metrics --help
Usage of ./eth-pools-metrics:
  -beacon-rpc-endpoint string
    	Address:Port of a eth2 beacon node endpoint (default "localhost:4000")
  -network string
    	Ethereum 2.0 network mainnet|prater|pyrmont (default "mainnet")
  -prometheus-port int
    	Prometheus port to listen to (default 9500)
  -withdrawal-credentials string
    	Hex withdrawal credentials following the spec, without 0x prefix
```

```console
$ go build
$ ./eth-pools-metrics \
--network=mainnet \
--withdrawal-credentials="004f58172d06b6d54c015d688511ad5656450933aff85dac123cd09410a0825c"
```

## TODO:
* Get validators per depositor address
* Experiment with sync committee duties
* Flag worst performing validators
