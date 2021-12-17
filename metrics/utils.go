package metrics

import (
	"context"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

func GetChainHead(ctx context.Context, beaconChainClient ethpb.BeaconChainClient) (*ethpb.ChainHead, error) {
	chainHead, err := beaconChainClient.GetChainHead(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "error getting chain head")
	}
	return chainHead, nil
}

func BoolToUint64(in bool) uint64 {
	if in {
		return uint64(1)
	}
	return uint64(0)
}

func GetSlotsInEpoch(ctx context.Context, beaconChainClient ethpb.BeaconChainClient) (uint64, error) {
	beaconConfig, err := GetBeaconConfig(ctx, beaconChainClient)
	if err != nil {
		return 0, errors.Wrap(err, "error getting beacon config")
	}

	slotsInEpochStr := beaconConfig.Config["SlotsPerEpoch"]
	slotsInEpoch, err := strconv.ParseUint(slotsInEpochStr, 10, 64)

	if err != nil {
		return 0, errors.Wrap(err, "error parsing slotsInEpoch string to uint64")
	}

	return slotsInEpoch, nil
}

func GetBeaconConfig(ctx context.Context, beaconChainClient ethpb.BeaconChainClient) (*ethpb.BeaconConfig, error) {
	beaconConfig, err := beaconChainClient.GetBeaconConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "error getting beacon config")
	}
	return beaconConfig, nil
}
