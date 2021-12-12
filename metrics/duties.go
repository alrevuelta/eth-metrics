package metrics

import (
	"context"
	"eth-pools-metrics/prometheus" // TODO: Set github prefix when released
	"fmt"
	"github.com/pkg/errors"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"time"
)

type Duty struct {
	valIndex uint64
	slot     ethTypes.Slot
}

// Continuously reports scheduled and fulfilled duties for the validators for
// the latest finalized epoch
func (a *Metrics) StreamDuties() {
	lastEpoch := uint64(0)
	for {
		if a.activeKeys == nil {
			log.Warn("No active keys to get duties")
			time.Sleep(10 * time.Second)
			continue
		}

		head, err := GetChainHead(context.Background(), a.beaconChainClient)
		if err != nil {
			log.Error("error getting chain head: ", err)
		}
		if uint64(head.FinalizedEpoch) <= lastEpoch {
			time.Sleep(5 * time.Second)
			continue
		}
		log.Info("Fetching duties for epoch: ", head.FinalizedEpoch)
		duties, blocks, err := a.FetchDuties(context.Background(), uint64(head.FinalizedEpoch))
		if err != nil {
			log.Error("could not get duties: ", err)
			continue
		}
		nOfScheduledBlocks, nOfProposedBlocks := a.getProposalDuties(duties, blocks)

		prometheus.NOfScheduledBlocks.Set(float64(nOfScheduledBlocks))
		prometheus.NOfProposedBlocks.Set(float64(nOfProposedBlocks))

		log.WithFields(log.Fields{
			"Epoch":           head.FinalizedEpoch,
			"RequestedDuties": nOfScheduledBlocks,
			"PerformedDuties": nOfProposedBlocks,
		}).Info("Block proposals duties:")
		lastEpoch = uint64(head.FinalizedEpoch)
	}
}

func (a *Metrics) FetchDuties(
	ctx context.Context,
	epoch uint64) (
	*ethpb.DutiesResponse,
	*ethpb.ListBeaconBlocksResponse,
	error) {

	dutReq := &ethpb.DutiesRequest{
		Epoch:      ethTypes.Epoch(epoch),
		PublicKeys: a.activeKeys,
	}

	// TODO: Move this
	chunkSize := 2000
	duties, err := a.prysmConcurrent.ParalelGetDuties(ctx, dutReq, chunkSize)

	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get duties")
	}

	// Get the blocks in the current epoch
	blocks, err := a.beaconChainClient.ListBeaconBlocks(ctx, &ethpb.ListBlocksRequest{
		QueryFilter: &ethpb.ListBlocksRequest_Epoch{
			Epoch: ethTypes.Epoch(epoch),
		},
	})

	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get blocks")
	}

	return duties, blocks, nil
}

// Returns the number of duties in an epoch for all our validators and the number
// of performed proposals
func (a *Metrics) getProposalDuties(
	duties *ethpb.DutiesResponse,
	blocks *ethpb.ListBeaconBlocksResponse) (uint64, uint64) {

	if duties == nil {
		log.Warn("No data is available to calculate the duties")
		return 0, 0
	}

	// Store the proposing duties that belongs to our validators
	proposalDuties := make([]Duty, 0)

	// Scan all duties in the given epoch
	for i := range duties.CurrentEpochDuties {
		// If there are any proposal duties append them
		if len(duties.CurrentEpochDuties[i].ProposerSlots) > 0 {
			// Pub Key is also available res.CurrentEpochDuties[i].PublicKey
			valIndex := uint64(duties.CurrentEpochDuties[i].ValidatorIndex)
			// Most likely there will be only a single proposal per epoch
			for _, propSlot := range duties.CurrentEpochDuties[i].ProposerSlots {
				proposalDuties = append(proposalDuties, Duty{valIndex, propSlot})
				log.WithFields(log.Fields{
					"PublicKey": fmt.Sprintf("%x", duties.CurrentEpochDuties[i].PublicKey),
					"ValIndex":  valIndex,
					"Slot":      propSlot,
					"Epoch":     uint64(propSlot) % a.slotsInEpoch,
				}).Info("Proposal Duty Found:")
			}
		}
	}

	// Just return if no proposal duties were found for us
	if len(proposalDuties) == 0 {
		return 0, 0
	}

	proposalsPerformed := uint64(0)

	// Iterate our validator proposal duties
	for _, duty := range proposalDuties {
		// Iterate all blocks and check if we proposed the ones we should
		for _, block := range blocks.BlockContainers {
			propIndex, slot, graffiti := getBlockParams(block)
			// If the block at the slot was proposed by us (valIndex)
			if duty.valIndex == propIndex && duty.slot == slot {
				log.WithFields(log.Fields{
					"ValIndex": propIndex,
					"Slot":     slot,
					"Epoch":    uint64(slot) % a.slotsInEpoch,
					"Graffiti": graffiti,
				}).Info("Proposal Duty Completion Verified:")
				proposalsPerformed++
				break
			}
		}
	}

	totalProposalDuties := uint64(len(proposalDuties))

	return totalProposalDuties, proposalsPerformed
}

func getBlockParams(block *ethpb.BeaconBlockContainer) (uint64, ethTypes.Slot, string) {
	var propIndex uint64
	var slot ethTypes.Slot
	var graffiti string

	if block.GetAltairBlock() == nil {
		propIndex = uint64(block.GetPhase0Block().Block.ProposerIndex)
		slot = block.GetPhase0Block().Block.Slot
		graffiti = fmt.Sprintf("%s", block.GetPhase0Block().Block.Body.Graffiti)
	} else {
		propIndex = uint64(block.GetAltairBlock().Block.ProposerIndex)
		slot = block.GetAltairBlock().Block.Slot
		graffiti = fmt.Sprintf("%s", block.GetAltairBlock().Block.Body.Graffiti)
	}
	return propIndex, slot, graffiti
}
