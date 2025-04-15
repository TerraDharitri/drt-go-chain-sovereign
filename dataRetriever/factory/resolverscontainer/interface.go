package resolverscontainer

import "github.com/TerraDharitri/drt-go-chain/dataRetriever"

// ShardResolversContainerFactoryCreator defines a shard resolvers container factory creator
type ShardResolversContainerFactoryCreator interface {
	CreateShardResolversContainerFactory(args FactoryArgs) (dataRetriever.ResolversContainerFactory, error)
	IsInterfaceNil() bool
}
