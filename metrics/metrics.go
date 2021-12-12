package metrics

import (
	"context"
	"eth-pools-metrics/prometheus"       // TODO: Set github prefix when released
	"eth-pools-metrics/prysm-concurrent" // TODO: Use Github prefix when released
	"eth-pools-metrics/thegraph"         // TODO: Use Github prefix when released
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
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
	go a.StreamValidatorPerformance()

	go func() {
		for {
			// TODO: Race condition with the depositedKeys

			// TODO: Take theGraph out of metrics
			pubKeysDeposited, err := a.theGraph.GetAllDepositedKeys()
			if err != nil {
				log.Error(err)
				time.Sleep(60 * 10 * time.Second)
				continue
			}
			log.Info("Number of deposited keys: ", len(pubKeysDeposited))
			a.depositedKeys = pubKeysDeposited

			// Get the status of all the validators
			valsStatus, err := a.validatorClient.MultipleValidatorStatus(context.Background(), &ethpb.MultipleValidatorStatusRequest{
				PublicKeys: a.depositedKeys,
			})
			if err != nil {
				log.Error(err)
				time.Sleep(60 * 10 * time.Second)
				continue
			}
			a.valsStatus = valsStatus

			// TODO: Get other status

			// Get the performance of the filtered validators
			activeKeys := FilterActiveValidators(valsStatus)
			a.activeKeys = activeKeys
			log.Info("Active validators: ", len(activeKeys))
			if len(a.activeKeys) == 0 {
				log.Error(err)
				time.Sleep(60 * 10 * time.Second)
				continue
			}

			prometheus.NOfValidators.Set(float64(len(a.activeKeys)))
			prometheus.NOfDepositedValidators.Set(float64(len(a.depositedKeys)))
			// TODO: Other status (slashed, etc)

			time.Sleep(60 * 10 * time.Second)
		}
	}()
}
