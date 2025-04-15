package factory

import (
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	storagerequesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/storageRequestersContainer"
	dataRetrieverMock "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
)

// ShardRequestersContainerCreatorMock -
type ShardRequestersContainerCreatorMock struct {
	CreateShardRequestersContainerFactoryCalled func(args storagerequesterscontainer.FactoryArgs) (dataRetriever.RequestersContainerFactory, error)
}

// CreateShardRequestersContainerFactory -
func (mock *ShardRequestersContainerCreatorMock) CreateShardRequestersContainerFactory(args storagerequesterscontainer.FactoryArgs) (dataRetriever.RequestersContainerFactory, error) {
	if mock.CreateShardRequestersContainerFactoryCalled != nil {
		return mock.CreateShardRequestersContainerFactoryCalled(args)
	}

	return &dataRetrieverMock.ShardRequestersContainerFactoryMock{}, nil
}

// IsInterfaceNil -
func (mock *ShardRequestersContainerCreatorMock) IsInterfaceNil() bool {
	return mock == nil
}
