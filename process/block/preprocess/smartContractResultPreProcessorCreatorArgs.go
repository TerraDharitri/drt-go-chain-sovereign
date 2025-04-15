package preprocess

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// SmartContractResultPreProcessorCreatorArgs is the struct containing the data needed to create a SmartContractResultPreProcessor
type SmartContractResultPreProcessorCreatorArgs struct {
	ScrDataPool                  dataRetriever.ShardedDataCacherNotifier
	Store                        dataRetriever.StorageService
	Hasher                       hashing.Hasher
	Marshalizer                  marshal.Marshalizer
	ScrProcessor                 process.SmartContractResultProcessor
	ShardCoordinator             sharding.Coordinator
	Accounts                     state.AccountsAdapter
	OnRequestSmartContractResult func(shardID uint32, txHashes [][]byte)
	GasHandler                   process.GasHandler
	EconomicsFee                 process.FeeHandler
	PubkeyConverter              core.PubkeyConverter
	BlockSizeComputation         BlockSizeComputationHandler
	BalanceComputation           BalanceComputationHandler
	EnableEpochsHandler          common.EnableEpochsHandler
	ProcessedMiniBlocksTracker   process.ProcessedMiniBlocksTracker
	TxExecutionOrderHandler      common.TxExecutionOrderHandler
}
