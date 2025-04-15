package factory

import (
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	storagerequesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/storageRequestersContainer"
)

// ShardRequestersContainerCreatorHandler defines a creator of shard requesters container creator
type ShardRequestersContainerCreatorHandler interface {
	CreateShardRequestersContainerFactory(args storagerequesterscontainer.FactoryArgs) (dataRetriever.RequestersContainerFactory, error)
	IsInterfaceNil() bool
}
