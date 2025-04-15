package preprocess

import (
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

// ScheduledTxsExecutionFactoryArgs holds all dependencies required by the process data factory to create components
type ScheduledTxsExecutionFactoryArgs struct {
	TxProcessor             process.TransactionProcessor
	TxCoordinator           process.TransactionCoordinator
	Storer                  storage.Storer
	Marshalizer             marshal.Marshalizer
	Hasher                  hashing.Hasher
	ShardCoordinator        sharding.Coordinator
	TxExecutionOrderHandler common.TxExecutionOrderHandler
}

type shardScheduledTxsExecutionFactory struct {
}

// NewShardScheduledTxsExecutionFactory creates a new shard scheduled txs execution factory
func NewShardScheduledTxsExecutionFactory() *shardScheduledTxsExecutionFactory {
	return &shardScheduledTxsExecutionFactory{}
}

// CreateScheduledTxsExecutionHandler creates a new scheduled txs execution handler for shard chain
func (stxef *shardScheduledTxsExecutionFactory) CreateScheduledTxsExecutionHandler(args ScheduledTxsExecutionFactoryArgs) (process.ScheduledTxsExecutionHandler, error) {
	return NewScheduledTxsExecution(
		args.TxProcessor,
		args.TxCoordinator,
		args.Storer,
		args.Marshalizer,
		args.Hasher,
		args.ShardCoordinator,
		args.TxExecutionOrderHandler,
	)
}

// IsInterfaceNil returns true if there is no value under the interface
func (stxef *shardScheduledTxsExecutionFactory) IsInterfaceNil() bool {
	return stxef == nil
}
