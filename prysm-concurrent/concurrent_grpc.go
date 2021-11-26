package prysmconcurrent

import (
	"context"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"google.golang.org/grpc"
	"sync"
	//"math/big"
	"golang.org/x/sync/errgroup"

	log "github.com/sirupsen/logrus"
)

type PrysmConcurrent struct {
	beaconRpcEndpoint string
	beaconChainClient ethpb.BeaconChainClient
	validatorClient   ethpb.BeaconNodeValidatorClient
	nodeClient        ethpb.NodeClient
}

func NewPrysmConcurrent(
	ctx context.Context,
	beaconRpcEndpoint string) (*PrysmConcurrent, error) {

	dialContext, err := grpc.DialContext(ctx, beaconRpcEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "could not create dial context")
	}

	beaconClient := ethpb.NewBeaconChainClient(dialContext)
	validatorClient := ethpb.NewBeaconNodeValidatorClient(dialContext)
	nodeClient := ethpb.NewNodeClient(dialContext)

	return &PrysmConcurrent{
		beaconRpcEndpoint: beaconRpcEndpoint,
		beaconChainClient: beaconClient,
		validatorClient:   validatorClient,
		nodeClient:        nodeClient,
	}, nil
}

func (a *PrysmConcurrent) ParalelGetDuties(
	ctx context.Context,
	req *ethpb.DutiesRequest,
	chunkSize int) (*ethpb.DutiesResponse, error) {

	epoch := req.Epoch
	activeKeys := req.PublicKeys

	var wg sync.WaitGroup
	var g errgroup.Group
	var lock sync.Mutex

	res := &ethpb.DutiesResponse{}

	for i := 0; i < len(activeKeys); i += chunkSize {
		wg.Add(1)

		i := i
		end := i + chunkSize

		if end > len(activeKeys) {
			end = len(activeKeys)
		}

		keyChunk := activeKeys[i:end]

		g.Go(func() error {
			lock.Lock()
			defer lock.Unlock()
			defer wg.Done()

			chunkReq := &ethpb.DutiesRequest{
				Epoch:      epoch,
				PublicKeys: keyChunk,
			}

			chunkDuties, err := a.validatorClient.GetDuties(ctx, chunkReq)
			if err != nil {
				return errors.Wrap(err, "could not get duties for validators")
			}
			res.Duties = append(res.Duties, chunkDuties.Duties...)
			res.CurrentEpochDuties = append(res.CurrentEpochDuties, chunkDuties.CurrentEpochDuties...)
			res.NextEpochDuties = append(res.NextEpochDuties, chunkDuties.NextEpochDuties...)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return res, nil
}

// TODO: Not implemented
func (a *PrysmConcurrent) ParalelGetMultipleValidatorStatus(ctx context.Context, req *ethpb.MultipleValidatorStatusRequest) (*ethpb.MultipleValidatorStatusResponse, error) {
	log.Warn("ParalelGetMultipleValidatorStatus is not implemented, using a naive single thread")

	valsStatus, err := a.validatorClient.MultipleValidatorStatus(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "could not get multiple validator status")
	}
	return valsStatus, nil
}

// TODO: Not implemented
func (a *PrysmConcurrent) ParalelGetValidatorPerformance(ctx context.Context, req *ethpb.ValidatorPerformanceRequest) (*ethpb.ValidatorPerformanceResponse, error) {
	log.Warn("ParalelGetValidatorPerformance is not implemented, using a naive single thread")

	valsPerformance, err := a.beaconChainClient.GetValidatorPerformance(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "could not get validator performance from beacon client")
	}
	return valsPerformance, nil
}
