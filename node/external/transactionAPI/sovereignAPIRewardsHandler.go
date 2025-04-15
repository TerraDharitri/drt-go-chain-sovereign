package transactionAPI

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	rewardTxData "github.com/TerraDharitri/drt-go-chain-core/data/rewardTx"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"

	"github.com/TerraDharitri/drt-go-chain/errors"
)

type sovereignAPIRewards struct {
	addressPubKeyConverter core.PubkeyConverter
}

// NewSovereignAPIRewardsHandler creates an api rewards handler for sovereign chain
func NewSovereignAPIRewardsHandler(addressPubKeyConverter core.PubkeyConverter) (*sovereignAPIRewards, error) {
	if check.IfNil(addressPubKeyConverter) {
		return nil, errors.ErrNilAddressPublicKeyConverter
	}

	return &sovereignAPIRewards{
		addressPubKeyConverter: addressPubKeyConverter,
	}, nil
}

// PrepareRewardTx prepares an api reward tx with sender and receiver shard as sovereign chain
func (sar *sovereignAPIRewards) PrepareRewardTx(tx *rewardTxData.RewardTx) *transaction.ApiTransactionResult {
	if check.IfNil(tx) {
		return &transaction.ApiTransactionResult{}
	}

	return &transaction.ApiTransactionResult{
		Tx:          tx,
		Type:        string(transaction.TxTypeReward),
		Round:       tx.GetRound(),
		Epoch:       tx.GetEpoch(),
		Value:       tx.GetValue().String(),
		Sender:      "sovereign",
		Receiver:    sar.addressPubKeyConverter.SilentEncode(tx.GetRcvAddr(), log),
		SourceShard: core.SovereignChainShardId,
	}
}

// IsInterfaceNil checks if the underlying pointer is nil
func (sar *sovereignAPIRewards) IsInterfaceNil() bool {
	return sar == nil
}
