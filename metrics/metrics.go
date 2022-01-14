package metrics

import (
	"context"
	"time"

	"github.com/alrevuelta/eth-pools-metrics/config"
	"github.com/alrevuelta/eth-pools-metrics/postgresql"
	prysmconcurrent "github.com/alrevuelta/eth-pools-metrics/prysm-concurrent"
	"github.com/alrevuelta/eth-pools-metrics/thegraph"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v2/time/slots"

	//log "github.com/sirupsen/logrus"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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
	theGraph       *thegraph.Thegraph
	postgresql     *postgresql.Postgresql

	// Slot and epoch and its raw data
	// TODO: Remove, each metric task has its pace
	Epoch uint64
	Slot  uint64

	PoolName string
}

func NewMetrics(
	ctx context.Context,
	config *config.Config) (*Metrics, error) {

	theGraph, err := thegraph.NewThegraph(
		config.Network,
		config.WithdrawalCredentials,
		config.FromAddress)

	if err != nil {
		return nil, errors.Wrap(err, "error creating thegraph")
	}

	dialContext, err := grpc.DialContext(ctx, config.BeaconRpcEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "could not create dial context")
	}

	beaconClient := ethpb.NewBeaconChainClient(dialContext)
	validatorClient := ethpb.NewBeaconNodeValidatorClient(dialContext)
	nodeClient := ethpb.NewNodeClient(dialContext)

	genesis, err := nodeClient.GetGenesis(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "error getting genesis info")
	}

	slotsInEpoch, err := GetSlotsInEpoch(ctx, beaconClient)
	if err != nil {
		return nil, errors.Wrap(err, "error getting slots in epoch from config")
	}

	prysmConcurrent, err := prysmconcurrent.NewPrysmConcurrent(ctx, config.BeaconRpcEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "error creating prysm concurrent")
	}

	var pg *postgresql.Postgresql
	if config.Postgres != "" {
		pg, err = postgresql.New(config.Postgres)
		if err != nil {
			return nil, errors.Wrap(err, "could not create postgresql")
		}
		err := pg.CreateTable()
		if err != nil {
			return nil, errors.Wrap(err, "error creating pool table to store data")
		}
	}

	return &Metrics{
		prysmConcurrent:   prysmConcurrent,
		theGraph:          theGraph,
		beaconChainClient: beaconClient,
		validatorClient:   validatorClient,
		nodeClient:        nodeClient,
		withCredList:      config.WithdrawalCredentials,
		genesisSeconds:    uint64(genesis.GenesisTime.Seconds),
		slotsInEpoch:      uint64(slotsInEpoch),
		postgresql:        pg,
		PoolName:          config.PoolName,
	}, nil
}

func (a *Metrics) Run() {
	go a.StreamDuties()
	go a.StreamRewards()
	go a.StreamDeposits()
	go a.StreamValidatorPerformance()
	go a.StreamValidatorStatus()
}

func (a *Metrics) EpochToTime(epoch uint64) (time.Time, error) {
	epochTime, err := slots.ToTime(uint64(a.genesisSeconds), ethTypes.Slot(epoch*a.slotsInEpoch))
	if err != nil {
		return time.Time{}, err
	}
	return epochTime, nil
}
