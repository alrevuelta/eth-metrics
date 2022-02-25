package pools

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/client"
)

type RocketpoolMinipool struct {
	Address     []byte
	Pubkey      []byte
	NodeAddress []byte
}

// Mainnet rocket pool contract
var rocketStorage = "0x1d8f8f00cfa6758d7bE78336684788Fb0ee0Fa46"

// Memory cache of already known minipools and its data
var MinipoolsByAddress map[string]*RocketpoolMinipool = make(map[string]*RocketpoolMinipool, 0)

var RocketPoolKeys [][]byte

func RocketPoolFetcher(eth1Address string) {
	todoSetAsFlag := 60 * time.Minute
	ticker := time.NewTicker(todoSetAsFlag)
	for ; true; <-ticker.C {
		keys, err := GetRocketPoolKeys(eth1Address)
		if err != nil {
			log.Error("could not get rocketpool keys: ", err)
		}
		RocketPoolKeys = keys
	}
}

func GetRocketPoolKeys(eth1Address string) ([][]byte, error) {
	log.Info("Fetching rocket pool keys")
	t0 := time.Now()
	proxy := client.NewEth1ClientProxy(60*time.Second, eth1Address)
	rp, err := rocketpool.NewRocketPool(proxy, common.HexToAddress(rocketStorage))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("bad contract address: %s", rocketStorage))
	}

	minipools, err := minipool.GetMinipoolAddresses(rp, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error getting minipool addresses")
	}

	statsNew := 0
	statsCache := 0

	// Get the validator pubkey for each minipool
	for _, minipoolAddress := range minipools {

		// Since this should not change, avoid fetching already known mini pools
		if _, exists := MinipoolsByAddress[minipoolAddress.Hex()]; exists {
			statsCache++
		} else {
			info, err := getMiniPoolInfo(rp, minipoolAddress)

			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("could not get minipool info: %s", minipoolAddress.Hex()))
			}

			MinipoolsByAddress[minipoolAddress.Hex()] = info
			statsNew++
		}
	}

	log.WithFields(log.Fields{
		"NewDetectedKeys": statsNew,
		"CachedKeys":      statsCache,
		"Duration":        time.Since(t0),
	}).Info("RocketPool Keys:")

	return getKeys(), nil
}

func getMiniPoolInfo(
	rp *rocketpool.RocketPool,
	address common.Address) (*RocketpoolMinipool, error) {

	mp, err := minipool.NewMinipool(rp, address)
	if err != nil {
		return nil, errors.Wrap(err, "error creating minipool")
	}

	nodeAddress, err := mp.GetNodeAddress(nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not get node address of minipool")
	}

	pubkey, err := minipool.GetMinipoolPubkey(
		rp,
		common.BytesToAddress(address.Bytes()),
		nil)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not get minipool key: %s", address))
	}

	info := &RocketpoolMinipool{
		Address:     address.Bytes(),
		NodeAddress: nodeAddress.Bytes(),
		Pubkey:      pubkey.Bytes(),
	}
	return info, nil
}

func getKeys() [][]byte {
	keys := make([][]byte, 0)
	for _, element := range MinipoolsByAddress {
		keys = append(keys, element.Pubkey)
	}
	return keys
}
