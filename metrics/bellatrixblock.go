package metrics

import (
	"math/big"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

func MevRewardInWei(bellatrixBlock bellatrix.SignedBeaconBlock) (*big.Int, string, error) {
	totalMevReward := big.NewInt(0)
	// this should be just 1, but just in case
	numTxs := 0
	mevFeeRecipient := ""
	for _, rawTx := range bellatrixBlock.Message.Body.ExecutionPayload.Transactions {
		tx, msg, err := DecodeTx(rawTx)
		_ = tx
		if err != nil {
			return nil, mevFeeRecipient, err
		}
		// This seems to happen in smart contrat deployments
		if msg.To() == nil {
			continue
		}
		// Note that its usually the last tx but we check all just in case
		feeRecipient := bellatrixBlock.Message.Body.ExecutionPayload.FeeRecipient.String()
		txTo := msg.To().String()
		txFrom := msg.From().String()

		// If the blocks contains a tx sending from vanilaFeeRec to an address
		// We consider that address the MEV fee recipient
		// Unsure if just checking the last tx is enough
		if strings.ToLower(feeRecipient) == strings.ToLower(txFrom) {
			totalMevReward.Add(totalMevReward, msg.Value())
			mevFeeRecipient = txTo
			/*
				log.WithFields(log.Fields{
					"Slot":         bellatrixBlock.Message.Slot,
					"Block":        bellatrixBlock.Message.Body.ExecutionPayload.BlockNumber,
					"ValIndex":     bellatrixBlock.Message.ProposerIndex,
					"FeeRecipient": bellatrixBlock.Message.Body.ExecutionPayload.FeeRecipient.String(),
					"To":           msg.To().String(),
					"MevReward":    msg.Value(),
					"TxHash":       tx.Hash().String(),
				}).Info("MEV transaction detected to pool")
			*/
			numTxs++
		}
		if numTxs > 1 {
			log.Warn("More than 1 MEV transaction detected in block. check if this is correct")
		}
	}
	return totalMevReward, mevFeeRecipient, nil

}

func ToBytes20(x []byte) [20]byte {
	var y [20]byte
	copy(y[:], x)
	return y
}

func DecodeTx(rawTx []byte) (*types.Transaction, *types.Message, error) {
	var tx types.Transaction
	err := tx.UnmarshalBinary(rawTx)
	if err != nil {
		return nil, nil, err
	}

	// Supports EIP-2930 and EIP-2718 and EIP-1559 and EIP-155 and legacy transactions.
	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return nil, nil, err
	}
	return &tx, &msg, err
}
