package metrics

import (
	//"bytes"
	"context"
	"encoding/hex"
	"eth-pools-metrics/prometheus"       // TODO: Set github prefix when released
	"eth-pools-metrics/prysm-concurrent" // TODO: Use Github prefix when released
	"eth-pools-metrics/thegraph"         // TODO: Use Github prefix when released
	"fmt"
	"github.com/pkg/errors"
	ethTypes "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/v2/config/params"
	ethpb "github.com/prysmaticlabs/prysm/v2/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v2/time/slots"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/big"
	"strconv"
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
	duties          *ethpb.DutiesResponse
	blocks          *ethpb.ListBeaconBlocksResponse

	// Calculated metrics
	ValsMetrics ValidatorsMetrics
}

type ValidatorsMetrics struct {
	NOfValidators               uint64
	NOfTotalVotes               uint64
	NOfIncorrectSource          uint64
	NOfIncorrectTarget          uint64
	NOfIncorrectHead            uint64
	NOfValsWithDecreasedBalance uint64

	BalanceDecreasedPercent float64

	CumulativeRewards *big.Int
	DepositedAmount   *big.Int
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

	slotsInEpoch, err := getSlotsInEpoch(ctx, beaconClient)
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
	a.StreamDuties()

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
			activeKeys := filterActiveValidators(valsStatus)
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

	go func() {
		for {
			if a.activeKeys == nil {
				log.Warn("No active keys to calculate the rewards")
				time.Sleep(30 * time.Second)
				continue
			}
			head, err := GetChainHead(context.Background(), a.beaconChainClient)
			if err != nil {
				log.Error("error getting chain head", err)
			}
			cumulativeRewards, depositedAmount, err := a.GetRewards(context.Background(), uint64(head.FinalizedEpoch))
			if err != nil {
				log.Error("could not get rewards and balances", err)
				time.Sleep(30 * time.Second)
				continue
			}

			a.ValsMetrics.CumulativeRewards = cumulativeRewards
			a.ValsMetrics.DepositedAmount = depositedAmount
			prometheus.DepositedAmount.Set(float64(a.ValsMetrics.DepositedAmount.Uint64()))
			prometheus.CumulativeRewards.Set(float64(a.ValsMetrics.CumulativeRewards.Uint64()))

			log.WithFields(log.Fields{
				"Epoch":             uint64(head.FinalizedEpoch),
				"DepositedAmount":   a.ValsMetrics.DepositedAmount.Uint64(),
				"CumulativeRewards": a.ValsMetrics.CumulativeRewards.Uint64(),
			}).Info("Rewards/Balances:")

			// Do not fetch every epoch. For a large number of validators it would be too much
			// TODO: Set as config parameter
			time.Sleep(30 * 60 * time.Second)
		}
	}()

	go func() {
		for {

			time.Sleep(2 * time.Second)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			// Fetch needed data to run the metrics
			newData, err := a.FetchData(ctx)
			if err != nil {
				log.WithError(err).Warn("Failed to fetch metrics data")
				//time.Sleep(time.Minute)
				continue
			}

			if !newData {
				continue
			}

			// Calculate metrics on the fetched data
			a.CalculateMetrics()

			// Log the information
			a.Monitor()

			// Update prometheus metrics
			// TODO Set epoch also??
			// TODO: Pass the whole struct to prometheus module
			prometheus.NOfTotalVotes.Set(float64(a.ValsMetrics.NOfTotalVotes))
			prometheus.NOfIncorrectSource.Set(float64(a.ValsMetrics.NOfIncorrectSource))
			prometheus.NOfIncorrectTarget.Set(float64(a.ValsMetrics.NOfIncorrectTarget))
			prometheus.NOfIncorrectHead.Set(float64(a.ValsMetrics.NOfIncorrectHead))
			prometheus.BalanceDecreasedPercent.Set(a.ValsMetrics.BalanceDecreasedPercent)
		}
	}()
}

func (a *Metrics) GetRewards(ctx context.Context, epoch uint64) (*big.Int, *big.Int, error) {
	// Nov 2021: Balances - Deposits matches the rewards
	// but this may change once withdrawals are enabled
	balances, err := a.GetBalances(ctx, epoch)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get balances from beacon chain")
	}

	totalRewards := big.NewInt(0)
	totalDeposits := big.NewInt(0)
	for _, b := range balances {
		status := ethpb.ValidatorStatus(ethpb.ValidatorStatus_value[b.Status])
		if status == ethpb.ValidatorStatus_ACTIVE ||
			status == ethpb.ValidatorStatus_EXITING ||
			status == ethpb.ValidatorStatus_SLASHING {
			deposit := big.NewInt(0).SetUint64(depositInGigaWei)
			balance := big.NewInt(0).SetUint64(b.Balance)
			reward := big.NewInt(0).Sub(balance, deposit)
			totalRewards.Add(totalRewards, reward)
			totalDeposits.Add(totalDeposits, deposit)
		}
	}
	return totalRewards, totalDeposits, nil
}

func (a *Metrics) GetBalances(ctx context.Context, epoch uint64) ([]*ethpb.ValidatorBalances_Balance, error) {
	request := ethpb.ListValidatorBalancesRequest{
		QueryFilter: &ethpb.ListValidatorBalancesRequest_Epoch{ethTypes.Epoch(epoch)},
		PublicKeys:  a.activeKeys,
	}

	balances := make([]*ethpb.ValidatorBalances_Balance, 0)
	for {
		resp, err := a.beaconChainClient.ListValidatorBalances(ctx, &request)
		if err != nil {
			return nil, errors.Wrap(err, "error getting validator balances")
		}

		balances = append(balances, resp.Balances...)
		log.Info("NextPageToken: ", resp.NextPageToken, " length: ", len(resp.Balances))

		if resp.NextPageToken == "" {
			break
		} else {
			request.PageToken = resp.NextPageToken
		}
	}
	return balances, nil
}

func GetChainHead(ctx context.Context, beaconChainClient ethpb.BeaconChainClient) (*ethpb.ChainHead, error) {
	chainHead, err := beaconChainClient.GetChainHead(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "error getting chain head")
	}
	return chainHead, nil
}

func filterActiveValidators(vals *ethpb.MultipleValidatorStatusResponse) [][]byte {
	activeKeys := make([][]byte, 0)
	for i := range vals.PublicKeys {
		if IsKeyActive(vals.Statuses[i].Status) {
			activeKeys = append(activeKeys, vals.PublicKeys[i])
		}
	}
	return activeKeys
}

//Fetches data from the beacon chain for a given set of validators. Note
//that not all request accepts the epoch as input, so this function takes
//care of synching with the beacon so that all fetched data refers to the same
//epoch
func (a *Metrics) FetchData(ctx context.Context) (bool, error) {
	head, err := GetChainHead(ctx, a.beaconChainClient)
	if err != nil {
		return false, errors.Wrap(err, "error getting chain head")
	}

	// Run metrics in already completed epochs
	metricsEpoch := uint64(head.HeadEpoch) - 1
	metricsSlot := uint64(head.HeadSlot)

	log.Info("Slot: ", ethTypes.Slot(metricsSlot)%params.BeaconConfig().SlotsPerEpoch)

	if a.depositedKeys == nil {
		log.Warn("No active keys to get vals performance")
		time.Sleep(30 * time.Second)
		return false, nil
	}

	// Wait until the last slot to ensure all attestations are included
	if a.Epoch >= metricsEpoch || !slots.IsEpochEnd(head.HeadSlot) {
		return false, nil
	}

	slotTime, err := slots.ToTime(uint64(a.genesisSeconds), ethTypes.Slot(head.HeadSlot+1))

	if err != nil {
		return false, errors.Wrap(err, "could not get next slot time")
	}

	// Set as deadline the begining of the first slot of the next epoch
	ctx, cancel := context.WithDeadline(ctx, slotTime)
	defer cancel()

	a.Epoch = metricsEpoch
	a.Slot = metricsSlot

	log.WithFields(log.Fields{
		"Epoch": metricsEpoch,
		"Slot":  metricsSlot,
		// zero-indexed
		"SlotInEpoch": ethTypes.Slot(metricsSlot) % params.BeaconConfig().SlotsPerEpoch,
	}).Info("Fetching new validators info")

	req := &ethpb.ValidatorPerformanceRequest{
		PublicKeys: a.activeKeys,
	}

	valsPerformance, err := a.beaconChainClient.GetValidatorPerformance(ctx, req)
	if err != nil {
		return false, errors.Wrap(err, "could not get validator performance from beacon client")
	}

	a.valsPerformance = valsPerformance

	for i := range valsPerformance.MissingValidators {
		log.WithFields(log.Fields{
			"Epoch":   a.Epoch,
			"Address": hex.EncodeToString(a.valsPerformance.MissingValidators[i]),
		}).Warn("Validator performance not found in beacon chain")
	}

	log.Info("Remaining time for next slot: ", ctx)

	return true, nil
}

func (a *Metrics) CalculateMetrics() {
	// Calculate the metrics
	nOfTotalVotes, nOfIncorrectSource, nOfIncorrectTarget, nOfIncorrectHead := a.getIncorrectAttestations()
	nOfValsWithDecreasedBalance, nOfValidators := a.getNumOfBalanceDecreasedVals()

	// And store them
	a.ValsMetrics.NOfValidators = nOfValidators
	a.ValsMetrics.NOfTotalVotes = nOfTotalVotes
	a.ValsMetrics.NOfIncorrectSource = nOfIncorrectSource
	a.ValsMetrics.NOfIncorrectTarget = nOfIncorrectTarget
	a.ValsMetrics.NOfIncorrectHead = nOfIncorrectHead
	a.ValsMetrics.NOfValsWithDecreasedBalance = nOfValsWithDecreasedBalance
	a.ValsMetrics.BalanceDecreasedPercent = (float64(nOfValsWithDecreasedBalance) / float64(nOfValidators)) * 100
}

// Gets the total number of votes and the incorrect ones
func (a *Metrics) getIncorrectAttestations() (uint64, uint64, uint64, uint64) {
	// The source is the attestation itself
	// https://pintail.xyz/posts/validator-rewards-in-practice/?s=03#attestation-efficiency
	nOfIncorrectSource := uint64(0)
	nOfIncorrectTarget := uint64(0)
	nOfIncorrectHead := uint64(0)
	for i := range a.valsPerformance.PublicKeys {
		nOfIncorrectSource += boolToUint64(!a.valsPerformance.CorrectlyVotedSource[i])
		nOfIncorrectTarget += boolToUint64(!a.valsPerformance.CorrectlyVotedTarget[i])
		nOfIncorrectHead += boolToUint64(!a.valsPerformance.CorrectlyVotedHead[i])
		// since missing source is the most severe, log it
		if !a.valsPerformance.CorrectlyVotedSource[i] {
			log.Info("Key that missed the attestation: ", hex.EncodeToString(a.valsPerformance.PublicKeys[i]), "--", a.valsPerformance.CorrectlyVotedSource[i], "--", a.valsPerformance.BalancesAfterEpochTransition[i])
		}
	}

	// Each validator contains three votes: source, target and head
	nOfTotalVotes := uint64(len(a.valsPerformance.PublicKeys)) * 3

	return nOfTotalVotes, nOfIncorrectSource, nOfIncorrectTarget, nOfIncorrectHead
}

// Gets the total number of validators and the ones that decreased in value
func (a *Metrics) getNumOfBalanceDecreasedVals() (uint64, uint64) {
	nOfValsWithDecreasedBalance := uint64(0)
	for i := range a.valsPerformance.PublicKeys {
		if a.valsPerformance.BalancesAfterEpochTransition[i] < a.valsPerformance.BalancesBeforeEpochTransition[i] {
			log.Info("Key with decr balance: ", hex.EncodeToString(a.valsPerformance.PublicKeys[i]), "--", a.valsPerformance.BalancesBeforeEpochTransition[i], "--", a.valsPerformance.BalancesAfterEpochTransition[i])
			nOfValsWithDecreasedBalance++
		}
	}
	nOfValidators := uint64(len(a.valsPerformance.PublicKeys))

	return nOfValsWithDecreasedBalance, nOfValidators
}

func (a *Metrics) Monitor() {
	logEpochSlot := log.WithFields(log.Fields{
		"Epoch": a.Epoch,
		"Slot":  a.Slot,
	})

	logEpochSlot.WithFields(log.Fields{
		"nOfTotalVotes":      a.ValsMetrics.NOfTotalVotes,
		"nOfIncorrectSource": a.ValsMetrics.NOfIncorrectSource,
		"nOfIncorrectTarget": a.ValsMetrics.NOfIncorrectTarget,
		"nOfIncorrectHead":   a.ValsMetrics.NOfIncorrectHead,
	}).Info("Incorrect voting:")

	logEpochSlot.WithFields(log.Fields{
		"ActiveValidators":    len(a.activeKeys),
		"DepositedValidators": len(a.depositedKeys),
		"SlashedValidators":   "TODO",
		"ExitingValidators":   "TODO",
		"OtherStates":         "TODO",
	}).Info("Validator Status:")

	logEpochSlot.WithFields(log.Fields{
		"PercentIncorrectSource": (float64(a.ValsMetrics.NOfIncorrectSource) / float64(a.ValsMetrics.NOfTotalVotes)) * 100,
		"PercentIncorrectTarget": (float64(a.ValsMetrics.NOfIncorrectTarget) / float64(a.ValsMetrics.NOfTotalVotes)) * 100,
		"PercentIncorrectHead":   (float64(a.ValsMetrics.NOfIncorrectHead) / float64(a.ValsMetrics.NOfTotalVotes)) * 100,
	}).Info("Incorrect voting percents:")

	logEpochSlot.WithFields(log.Fields{
		"nOfValidators":               a.ValsMetrics.NOfValidators,
		"nOfValsWithDecreasedBalance": a.ValsMetrics.NOfValsWithDecreasedBalance,
		"balanceDecreasedPercent":     a.ValsMetrics.BalanceDecreasedPercent,
	}).Info("Balance decreased:")
}

// TODO: Move
func boolToUint64(in bool) uint64 {
	if in {
		return uint64(1)
	}
	return uint64(0)
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

// TODO: Move
func getSlotsInEpoch(ctx context.Context, beaconChainClient ethpb.BeaconChainClient) (uint64, error) {
	beaconConfig, err := GetBeaconConfig(ctx, beaconChainClient)
	if err != nil {
		return 0, errors.Wrap(err, "error getting beacon config")
	}

	/* Will most likely be 32 but it may change in other networks */
	slotsInEpochStr := beaconConfig.Config["SlotsPerEpoch"]
	slotsInEpoch, err := strconv.ParseUint(slotsInEpochStr, 10, 64)

	if err != nil {
		return 0, errors.Wrap(err, "error parsing slotsInEpoch string to uint64")
	}

	return slotsInEpoch, nil
}

// TODO: Move
func GetBeaconConfig(ctx context.Context, beaconChainClient ethpb.BeaconChainClient) (*ethpb.BeaconConfig, error) {
	beaconConfig, err := beaconChainClient.GetBeaconConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "error getting beacon config")
	}
	return beaconConfig, nil
}

// TODO: Move
func IsKeyActive(status ethpb.ValidatorStatus) bool {
	if status == ethpb.ValidatorStatus_ACTIVE ||
		status == ethpb.ValidatorStatus_EXITING ||
		status == ethpb.ValidatorStatus_SLASHING {
		return true
	}
	return false
}
