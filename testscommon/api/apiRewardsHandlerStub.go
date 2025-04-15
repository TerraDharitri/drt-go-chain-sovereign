package api

import (
	rewardTxData "github.com/TerraDharitri/drt-go-chain-core/data/rewardTx"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
)

// APIRewardsHandlerStub -
type APIRewardsHandlerStub struct {
	PrepareRewardTxCalled func(tx *rewardTxData.RewardTx) *transaction.ApiTransactionResult
}

// PrepareRewardTx -
func (stub *APIRewardsHandlerStub) PrepareRewardTx(tx *rewardTxData.RewardTx) *transaction.ApiTransactionResult {
	if stub.PrepareRewardTxCalled != nil {
		return stub.PrepareRewardTxCalled(tx)
	}

	return &transaction.ApiTransactionResult{}
}

// IsInterfaceNil -
func (stub *APIRewardsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
