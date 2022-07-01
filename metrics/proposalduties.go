package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/alrevuelta/eth-pools-metrics/prometheus"

	"github.com/alrevuelta/eth-pools-metrics/schemas"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	log "github.com/sirupsen/logrus"
)

type ProposalDuties struct {
	httpClient    *http.Service
	eth1Endpoint  string
	eth2Endpoint  string
	fromAddresses []string
	poolNames     []string
}

func NewProposalDuties(
	eth1Endpoint string,
	eth2Endpoint string,
	fromAddresses []string,
	poolNames []string) (*ProposalDuties, error) {

	client, err := http.New(context.Background(),
		http.WithTimeout(60*time.Second),
		http.WithAddress(eth2Endpoint),
		http.WithLogLevel(zerolog.WarnLevel),
	)
	if err != nil {
		return nil, err
	}

	httpClient := client.(*http.Service)

	return &ProposalDuties{
		httpClient:    httpClient,
		eth2Endpoint:  eth2Endpoint,
		fromAddresses: fromAddresses,
		poolNames:     poolNames,
		eth1Endpoint:  eth1Endpoint,
	}, nil
}

func (p *ProposalDuties) RunProposalMetrics(
	activeKeys []uint64,
	poolName string,
	metrics *schemas.ProposalDutiesMetrics) error {

	poolProposals := getPoolProposalDuties(
		metrics,
		poolName,
		activeKeys)

	logProposalDuties(poolProposals, poolName)
	setPrometheusProposalDuties(poolProposals, poolName)
	return nil

}

func (p *ProposalDuties) GetProposalDuties(epoch uint64) ([]*api.ProposerDuty, error) {
	// Empty indexes to force fetching all duties
	indexes := make([]phase0.ValidatorIndex, 0)

	duties, err := p.httpClient.ProposerDuties(
		context.Background(),
		phase0.Epoch(epoch),
		indexes)

	if err != nil {
		return make([]*api.ProposerDuty, 0), err
	}

	return duties, nil
}

func (p *ProposalDuties) GetProposedBlocks(epoch uint64) ([]*api.BeaconBlockHeader, error) {

	epochBlockHeaders := make([]*api.BeaconBlockHeader, 0)
	slotsInEpoch := uint64(32)

	slotWithinEpoch := uint64(0)
	for slotWithinEpoch < slotsInEpoch {
		epochStr := strconv.FormatUint(epoch*slotsInEpoch+slotWithinEpoch, 10)

		blockHeader, err := p.httpClient.BeaconBlockHeader(context.Background(), epochStr)
		if err != nil {
			return epochBlockHeaders, err
		}
		epochBlockHeaders = append(epochBlockHeaders, blockHeader)
		slotWithinEpoch++
	}

	return epochBlockHeaders, nil
}

func (p *ProposalDuties) GetProposalMetrics(
	proposalDuties []*api.ProposerDuty,
	proposedBlocks []*api.BeaconBlockHeader) (schemas.ProposalDutiesMetrics, error) {

	proposalMetrics := schemas.ProposalDutiesMetrics{
		Epoch:     0,
		Scheduled: make([]schemas.Duty, 0),
		Proposed:  make([]schemas.Duty, 0),
		Missed:    make([]schemas.Duty, 0),
	}

	if proposalDuties[0].Slot != proposedBlocks[0].Header.Message.Slot {
		return proposalMetrics, errors.New("duties and proposals contains different slots")
	}
	if len(proposalDuties) != len(proposedBlocks) {
		return proposalMetrics, errors.New("duties and blocks have different sizes")
	}

	proposalMetrics.Epoch = uint64(proposalDuties[0].Slot) / 32

	for _, duty := range proposalDuties {
		proposalMetrics.Scheduled = append(
			proposalMetrics.Scheduled,
			schemas.Duty{
				ValIndex: uint64(duty.ValidatorIndex),
				Slot:     uint64(duty.Slot),
				Graffiti: "NA",
			})
	}

	for _, block := range proposedBlocks {
		proposalMetrics.Proposed = append(
			proposalMetrics.Proposed,
			schemas.Duty{
				ValIndex: uint64(block.Header.Message.ProposerIndex),
				Slot:     uint64(block.Header.Message.Slot),
				Graffiti: "TODO",
			})

	}

	return proposalMetrics, nil
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

// TODO: This is very inefficient
func getPoolProposalDuties(
	metrics *schemas.ProposalDutiesMetrics,
	poolName string,
	activeValidatorIndexes []uint64) *schemas.ProposalDutiesMetrics {

	poolDuties := schemas.ProposalDutiesMetrics{
		Epoch:     metrics.Epoch,
		Scheduled: make([]schemas.Duty, 0),
		Proposed:  make([]schemas.Duty, 0),
		Missed:    make([]schemas.Duty, 0),
	}

	// Check if this pool has any assigned proposal duties
	for i := range metrics.Scheduled {
		if IsValidatorIn(metrics.Scheduled[i].ValIndex, activeValidatorIndexes) {
			poolDuties.Scheduled = append(poolDuties.Scheduled, metrics.Scheduled[i])
		}
		if IsValidatorIn(metrics.Proposed[i].ValIndex, activeValidatorIndexes) {
			poolDuties.Proposed = append(poolDuties.Proposed, metrics.Proposed[i])
		}
	}

	poolDuties.Missed = getMissedDuties(poolDuties.Scheduled, poolDuties.Proposed)

	return &poolDuties
}

func logProposalDuties(
	poolDuties *schemas.ProposalDutiesMetrics,
	poolName string) {

	for _, d := range poolDuties.Scheduled {
		log.WithFields(log.Fields{
			"PoolName":       poolName,
			"ValIndex":       d.ValIndex,
			"Slot":           d.Slot,
			"Epoch":          poolDuties.Epoch,
			"TotalScheduled": len(poolDuties.Scheduled),
		}).Info("Scheduled Duty")
	}

	for _, d := range poolDuties.Proposed {
		log.WithFields(log.Fields{
			"PoolName":      poolName,
			"ValIndex":      d.ValIndex,
			"Slot":          d.Slot,
			"Epoch":         poolDuties.Epoch,
			"Graffiti":      d.Graffiti,
			"TotalProposed": len(poolDuties.Proposed),
		}).Info("Proposed Duty")
	}

	for _, d := range poolDuties.Missed {
		log.WithFields(log.Fields{
			"PoolName":    poolName,
			"ValIndex":    d.ValIndex,
			"Slot":        d.Slot,
			"Epoch":       poolDuties.Epoch,
			"TotalMissed": len(poolDuties.Missed),
		}).Info("Missed Duty")
	}
}

func setPrometheusProposalDuties(
	metrics *schemas.ProposalDutiesMetrics,
	poolName string) {

	prometheus.NOfProposedBlocks.WithLabelValues(
		poolName).Set(float64(len(metrics.Proposed)))

	prometheus.NOfMissedBlocks.WithLabelValues(
		poolName).Set(float64(len(metrics.Missed)))

	for _, d := range metrics.Proposed {
		_ = d
		/* TODO: Not sure, add pool label
		prometheus.ProposedBlocks.WithLabelValues(
			UToStr(metrics.Epoch),
			UToStr(d.ValIndex)).Inc()
		*/
	}

	for _, d := range metrics.Missed {
		_ = d
		/* TODO: Not sure, add pool label
		prometheus.MissedBlocks.WithLabelValues(
			UToStr(metrics.Epoch),
			UToStr(d.ValIndex)).Inc()
		*/
	}
}

/*
func (a *Metrics) RunProposalMetrics() {
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
*/

/*


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

*/
