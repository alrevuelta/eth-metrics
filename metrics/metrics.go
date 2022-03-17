package metrics

import (
	"context"
	"time"

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

	// Slot and epoch and its raw data
	// TODO: Remove, each metric task has its pace
	Epoch uint64
	Slot  uint64

	PoolNames []string
}

func NewMetrics(
	ctx context.Context,
	config *config.Config) (*Metrics, error) {

	/*
		theGraph, err := thegraph.NewThegraph(
			//config.Network,
			"ignore",
			config.WithdrawalCredentials,
			config.FromAddress)

		if err != nil {
			return nil, errors.Wrap(err, "error creating thegraph")
		}
	*/
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

	return &Metrics{
		//theGraph:     theGraph,
		withCredList: config.WithdrawalCredentials,
		fromAddrList: config.FromAddress,
		//genesisSeconds:    uint64(genesis.GenesisTime.Seconds),
		//slotsInEpoch:      uint64(slotsInEpoch),
		eth1Address: config.Eth1Address,
		eth2Address: config.Eth2Address,
		postgresql:  pg,
		PoolNames:   config.PoolNames,
	}, nil
}

func (a *Metrics) Run() {
	//go a.StreamDuties()
	//go a.StreamRewards()
	//go a.StreamDeposits()
	//go a.StreamValidatorPerformance()
	//go a.StreamValidatorStatus()
	bc, err := NewBeaconState(
		a.eth1Address,
		a.eth2Address,
		a.postgresql,
		a.fromAddrList,
		a.PoolNames,
	)
	if err != nil {
		log.Error(err)
		// TODO: Add return here.
	}
	for _, poolName := range a.PoolNames {
		if poolName == "rocketpool" {
			go pools.RocketPoolFetcher(a.eth1Address)
			break
		}
	}
	go bc.Run()
}

func (a *Metrics) EpochToTime(epoch uint64) (time.Time, error) {
	epochTime, err := slots.ToTime(uint64(a.genesisSeconds), ethTypes.Slot(epoch*a.slotsInEpoch))
	if err != nil {
		return time.Time{}, err
	}
	return epochTime, nil
}
