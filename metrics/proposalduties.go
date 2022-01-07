package metrics

import (
	"context"
	"fmt"
	"github.com/alrevuelta/eth-pools-metrics/prometheus"
	"github.com/alrevuelta/eth-pools-metrics/schemas"
	"github.com/pkg/errors"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

// Continuously reports scheduled and fulfilled duties for the validators for
// the latest finalized epoch
func (a *Metrics) StreamDuties() {
	lastEpoch := uint64(0)
	for {
		if a.validatingKeys == nil {
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
		metrics := getProposalDuties(duties, blocks)
		metrics.Epoch = uint64(head.FinalizedEpoch)

		logProposalDuties(metrics)
		setPrometheusProposalDuties(metrics)

		lastEpoch = uint64(head.FinalizedEpoch)

		// Temporal fix to memory leak. Perhaps having an infinite loop
		// inside a routinne is not a good idea. TODO
		runtime.GC()
	}
}

func logProposalDuties(metrics *schemas.ProposalDutiesMetrics) {
	for _, d := range metrics.Scheduled {
		log.WithFields(log.Fields{
			"ValIndex":       d.ValIndex,
			"Slot":           d.Slot,
			"Epoch":          metrics.Epoch,
			"TotalScheduled": len(metrics.Scheduled),
		}).Info("Scheduled Duty")
	}

	for _, d := range metrics.Proposed {
		log.WithFields(log.Fields{
			"ValIndex":      d.ValIndex,
			"Slot":          d.Slot,
			"Epoch":         metrics.Epoch,
			"Graffiti":      d.Graffiti,
			"TotalProposed": len(metrics.Proposed),
		}).Info("Proposed Duty")
	}

	for _, d := range metrics.Missed {
		log.WithFields(log.Fields{
			"ValIndex":    d.ValIndex,
			"Slot":        d.Slot,
			"Epoch":       metrics.Epoch,
			"TotalMissed": len(metrics.Missed),
		}).Info("Missed Duty")
	}
}

func setPrometheusProposalDuties(metrics *schemas.ProposalDutiesMetrics) {
	prometheus.NOfScheduledBlocks.Set(float64(len(metrics.Scheduled)))
	prometheus.NOfProposedBlocks.Set(float64(len(metrics.Proposed)))

	for _, d := range metrics.Proposed {
		prometheus.ProposedBlocks.WithLabelValues(
			UToStr(metrics.Epoch),
			UToStr(d.ValIndex)).Inc()
	}

	for _, d := range metrics.Missed {
		prometheus.MissedBlocks.WithLabelValues(
			UToStr(metrics.Epoch),
			UToStr(d.ValIndex)).Inc()
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
		PublicKeys: a.validatingKeys,
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
func getProposalDuties(
	duties *ethpb.DutiesResponse,
	blocks *ethpb.ListBeaconBlocksResponse) *schemas.ProposalDutiesMetrics {

	metrics := &schemas.ProposalDutiesMetrics{
		Scheduled: make([]schemas.Duty, 0),
		Proposed:  make([]schemas.Duty, 0),
		Missed:    make([]schemas.Duty, 0),
	}

	if duties == nil {
		log.Warn("No data is available to calculate the duties")
		return metrics
	}

	// Scan all duties in the given epoch
	for i := range duties.CurrentEpochDuties {
		// If there are any proposal duties append them
		if len(duties.CurrentEpochDuties[i].ProposerSlots) > 0 {
			// Pub Key is also available res.CurrentEpochDuties[i].PublicKey
			valIndex := uint64(duties.CurrentEpochDuties[i].ValidatorIndex)
			// Most likely there will be only a single proposal per epoch
			for _, propSlot := range duties.CurrentEpochDuties[i].ProposerSlots {
				metrics.Scheduled = append(metrics.Scheduled, schemas.Duty{ValIndex: valIndex, Slot: propSlot})
			}
		}
	}

	// Just return if no proposal duties were found for us
	if len(metrics.Scheduled) == 0 {
		return metrics
	}

	// Iterate our validator proposal duties
	for _, duty := range metrics.Scheduled {
		// Iterate all blocks and check if we proposed the ones we should
		for _, block := range blocks.BlockContainers {
			propIndex, slot, graffiti := getBlockParams(block)
			// If the block at the slot was proposed by us (valIndex)
			if duty.ValIndex == propIndex && duty.Slot == slot {
				metrics.Proposed = append(metrics.Proposed, schemas.Duty{
					ValIndex: propIndex,
					Slot:     slot,
					Graffiti: graffiti})
				break
			}
		}
	}

	metrics.Missed = getMissedDuties(metrics.Scheduled, metrics.Proposed)

	return metrics
}

func getMissedDuties(scheduled []schemas.Duty, proposed []schemas.Duty) []schemas.Duty {
	missed := make([]schemas.Duty, 0)

	for _, s := range scheduled {
		found := false
		for _, p := range proposed {
			if s.Slot == p.Slot && s.ValIndex == p.ValIndex {
				found = true
				break
			}
		}
		if found == false {
			missed = append(missed, s)
		}
	}

	return missed
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
	// TODO: Add merge block when implemented
	return propIndex, slot, graffiti
}
