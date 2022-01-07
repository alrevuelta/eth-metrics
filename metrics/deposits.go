package metrics

import (
	//"github.com/alrevuelta/eth-pools-metrics/prometheus"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// TODO: Temporal solution:
// - TheGraph API calls has some limits, so we can't query in every epoch
// - Race condition with the depositedKeys
// - Fetches the deposits every hour
func (a *Metrics) StreamDeposits() {
	for {
		/*
		pubKeysDeposited, err := a.theGraph.GetAllDepositedKeys()
		if err != nil {
			log.Error(err)
			time.Sleep(10 * 60 * time.Second)
			continue
		}
		*/

		pubKeysDeposited, err := getHardcodedKeys()
		if err != nil {
			log.Fatal(err)
		}

		a.depositedKeys = pubKeysDeposited

		log.WithFields(log.Fields{
			"DepositedValidators": len(pubKeysDeposited),
			// TODO: Print epoch
			//"Slot":     slot,
			//"Epoch":    uint64(slot) % a.slotsInEpoch,
		}).Info("Deposits:")

		// Temporal fix to memory leak. Perhaps having an infinite loop
		// inside a routinne is not a good idea. TODO
		runtime.GC()

		time.Sleep(60 * 60 * time.Second)
	}
}

func getHardcodedKeys() ([][]byte, error) {
	// Kintsugi validators
	var keysStr = []string{
		"0x820575d85e0368bc5f2aa55a9b3a41ea804afe28b031a84f4b8a66b9a3c6b6ab39d6f46e020b860d6cbfa987248e92d6",
		"0x935daecc77617f127226edff6af7100fdbea02477bb9376608492a0b5c9706e30f8911f5075563e140afa182df22abfb",
	}

	keys := make([][]byte, 0)

	for _, keyStr := range keysStr {
		key, err := hexutil.Decode(keyStr)
		keys = append(keys, key)
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}
