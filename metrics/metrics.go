package metrics

import (
	"context"
	"eth-pools-metrics/prysm-concurrent" // TODO: Use Github prefix when released
	"eth-pools-metrics/thegraph"         // TODO: Use Github prefix when released
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	//log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	gigaWei          = uint64(1_000_000_000)
	depositInGigaWei = uint64(32) * gigaWei
)

// Max uint64 value
const maxUint64Value = uint64(^uint(0))

type Metrics struct {
	beaconChainClient ethpb.BeaconChainClient
	validatorClient   ethpb.BeaconNodeValidatorClient
	nodeClient        ethpb.NodeClient
	prysmConcurrent   *prysmconcurrent.PrysmConcurrent
	genesisSeconds    uint64
	slotsInEpoch      uint64

	depositedKeys [][]byte
	activeKeys    [][]byte
	withCredList  []string
	theGraph      *thegraph.Thegraph

	// Slot and epoch and its raw data
	Epoch           uint64
	Slot            uint64
	valsPerformance *ethpb.ValidatorPerformanceResponse
	valsStatus      *ethpb.MultipleValidatorStatusResponse
}

func NewMetrics(
	ctx context.Context,
	beaconRpcEndpoint string,
	network string,
	withCredList []string,
	fromAddresses []string) (*Metrics, error) {

	theGraph, err := thegraph.NewThegraph(network, withCredList, fromAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "error creating thegraph")
	}

	dialContext, err := grpc.DialContext(ctx, beaconRpcEndpoint, grpc.WithInsecure())
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

	prysmConcurrent, err := prysmconcurrent.NewPrysmConcurrent(ctx, beaconRpcEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "error creating prysm concurrent")
	}

	return &Metrics{
		prysmConcurrent:   prysmConcurrent,
		theGraph:          theGraph,
		beaconChainClient: beaconClient,
		validatorClient:   validatorClient,
		nodeClient:        nodeClient,
		withCredList:      withCredList,
		genesisSeconds:    uint64(genesis.GenesisTime.Seconds),
		slotsInEpoch:      uint64(slotsInEpoch),
	}, nil
}

func (a *Metrics) Run() {
	go a.StreamDuties()
	go a.StreamRewards()
	go a.StreamDeposits()
	go a.StreamValidatorPerformance()
	go a.StreamValidatorStatus()
}
