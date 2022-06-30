package metrics

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/rs/zerolog"

	"github.com/alrevuelta/eth-pools-metrics/config"
	"github.com/alrevuelta/eth-pools-metrics/pools"
	"github.com/alrevuelta/eth-pools-metrics/postgresql"
	prysmconcurrent "github.com/alrevuelta/eth-pools-metrics/prysm-concurrent"
	"github.com/alrevuelta/eth-pools-metrics/thegraph"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v2/time/slots"
	log "github.com/sirupsen/logrus"

	//log "github.com/sirupsen/logrus"
	ethTypes "github.com/prysmaticlabs/eth2-types"
)

type Metrics struct {
	beaconChainClient ethpb.BeaconChainClient
	validatorClient   ethpb.BeaconNodeValidatorClient
	nodeClient        ethpb.NodeClient
	prysmConcurrent   *prysmconcurrent.PrysmConcurrent
	genesisSeconds    uint64
	slotsInEpoch      uint64

	depositedKeys  [][]byte
	validatingKeys [][]byte
	withCredList   []string
	fromAddrList   []string
	eth1Address    string
	eth2Address    string
	theGraph       *thegraph.Thegraph
	postgresql     *postgresql.Postgresql

	httpClient *http.Service

	beaconState *BeaconState

	// Slot and epoch and its raw data
	// TODO: Remove, each metric task has its pace
	Epoch uint64
	Slot  uint64

	PoolNames []string
}

func NewMetrics(
	ctx context.Context,
	config *config.Config) (*Metrics, error) {

	/* TODO: Get from a http endpoint instead of prysm gRPC
	genesis, err := nodeClient.GetGenesis(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "error getting genesis info")
	}*/

	/* TODO: Get from a http endpoint instead of prysm gRPC
	slotsInEpoch, err := GetSlotsInEpoch(ctx, beaconClient)
	if err != nil {
		return nil, errors.Wrap(err, "error getting slots in epoch from config")
	}*/

	var pg *postgresql.Postgresql
	var err error
	if config.Postgres != "" {
		pg, err = postgresql.New(config.Postgres)
		if err != nil {
			return nil, errors.Wrap(err, "could not create postgresql")
		}
		err = pg.CreateTable()
		if err != nil {
			return nil, errors.Wrap(err, "error creating pool table to store data")
		}
	}

	for _, poolName := range config.PoolNames {
		if strings.HasSuffix(poolName, ".txt") {
			pubKeysDeposited, err := pools.ReadCustomValidatorsFile(poolName)
			if err != nil {
				log.Fatal(err)
			}
			log.Info("File: ", poolName, " contains ", len(pubKeysDeposited), " keys")

		}
	}

	client, err := http.New(context.Background(),
		http.WithTimeout(60*time.Second),
		http.WithAddress(config.Eth2Address),
		http.WithLogLevel(zerolog.WarnLevel),
	)
	if err != nil {
		return nil, err
	}

	httpClient := client.(*http.Service)

	return &Metrics{
		withCredList: config.WithdrawalCredentials,
		fromAddrList: config.FromAddress,
		//genesisSeconds:    uint64(genesis.GenesisTime.Seconds),
		//slotsInEpoch:      uint64(slotsInEpoch),
		eth1Address: config.Eth1Address,
		eth2Address: config.Eth2Address,
		postgresql:  pg,
		PoolNames:   config.PoolNames,
		httpClient:  httpClient,
	}, nil
}

func (a *Metrics) Run() {
	bc, err := NewBeaconState(
		a.eth1Address,
		a.eth2Address,
		a.postgresql,
		a.fromAddrList,
		a.PoolNames,
	)
	if err != nil {
		log.Fatal(err)
		// TODO: Add return here.
	}
	a.beaconState = bc
	for _, poolName := range a.PoolNames {
		if poolName == "rocketpool" {
			go pools.RocketPoolFetcher(a.eth1Address)
			break
		}
	}
	go a.Loop()
}

func (a *Metrics) Loop() {
	var prevEpoch uint64 = uint64(0)
	var prevBeaconState *spec.VersionedBeaconState = nil
	// TODO: Refactor and hoist some stuff out to a function
	for {
		// Before doing anything, check if we are in the next epoch
		headSlot, err := a.httpClient.NodeSyncing(context.Background())
		if err != nil {
			log.Error("Could not get node sync status:", err)
			continue
		}

		if headSlot.IsSyncing {
			log.Error("Node is not in sync")
			continue
		}

		currentEpoch := uint64(headSlot.HeadSlot)/uint64(32) - 3

		if prevEpoch >= currentEpoch {
			// do nothing
			time.Sleep(5 * time.Second)
			continue
		}

		currentBeaconState, err := a.beaconState.GetBeaconState(currentEpoch)
		if err != nil {
			prevBeaconState = nil
			log.Error("Error fetching beacon state:", err)
			continue
		}

		// if no prev beacon state is known, fetch it
		if prevBeaconState == nil {
			prevBeaconState, err = a.beaconState.GetBeaconState(currentEpoch - 1)
			if err != nil {
				log.Error(err)
				continue
			}
		}

		// Iterate all pools and calculate metrics using the fetched data
		for _, poolName := range a.PoolNames {
			poolName, pubKeys, err := a.GetValidatorKeys(poolName)
			if err != nil {
				log.Error("TODO", err)
				continue
			}
			a.beaconState.Run(pubKeys, poolName, currentBeaconState, prevBeaconState)
		}

		prevBeaconState = currentBeaconState
		prevEpoch = currentEpoch
	}
}

// Get the validator keys from different sources:
// - pool.txt: Opens the file and read the keys from it
// - rocketpool: Special case, see pools
// - poolname: Gets the keys from the address used for the deposit
func (a *Metrics) GetValidatorKeys(poolName string) (string, [][]byte, error) {
	var pubKeysDeposited [][]byte
	var err error
	if strings.HasSuffix(poolName, ".txt") {
		pubKeysDeposited, err = pools.ReadCustomValidatorsFile(poolName)
		if err != nil {
			log.Fatal(err)
		}
		// trim the file path and extension
		poolName = filepath.Base(poolName)
		poolName = strings.TrimSuffix(poolName, filepath.Ext(poolName))
	} else if poolName == "rocketpool" {
		pubKeysDeposited = pools.RocketPoolKeys
	} else {
		poolAddressList := pools.PoolsAddresses[poolName]
		pubKeysDeposited, err = a.postgresql.GetKeysByFromAddresses(poolAddressList)
		if err != nil {
			return "", nil, err
		}
	}
	return poolName, pubKeysDeposited, nil
}

func (a *Metrics) EpochToTime(epoch uint64) (time.Time, error) {
	epochTime, err := slots.ToTime(uint64(a.genesisSeconds), ethTypes.Slot(epoch*a.slotsInEpoch))
	if err != nil {
		return time.Time{}, err
	}
	return epochTime, nil
}
