package factory

import (
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/resolverscontainer"
	dataRetrieverMocks "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
)

// ResolversContainerFactoryMock -
type ResolversContainerFactoryMock struct {
	CreateShardResolversContainerFactoryCalled func(args resolverscontainer.FactoryArgs) (dataRetriever.ResolversContainerFactory, error)
}

// CreateShardResolversContainerFactory -
func (r *ResolversContainerFactoryMock) CreateShardResolversContainerFactory(args resolverscontainer.FactoryArgs) (dataRetriever.ResolversContainerFactory, error) {
	if r.CreateShardResolversContainerFactoryCalled != nil {
		return r.CreateShardResolversContainerFactory(args)
	}
	return &dataRetrieverMocks.ShardResolversContainerFactoryMock{}, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (r *ResolversContainerFactoryMock) IsInterfaceNil() bool {
	return r == nil
}
