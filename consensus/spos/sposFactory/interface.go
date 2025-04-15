package sposFactory

import (
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/broadcast"
)

// BroadCastShardMessengerFactoryHandler defines a shard messenger factory handler
type BroadCastShardMessengerFactoryHandler interface {
	CreateShardChainMessenger(args broadcast.ShardChainMessengerArgs) (consensus.BroadcastMessenger, error)
	IsInterfaceNil() bool
}
