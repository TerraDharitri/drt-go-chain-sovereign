package factory

import (
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	requesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/requestersContainer"
	dataRetrieverMocks "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
)

// RequestersContainerFactoryMock -
type RequestersContainerFactoryMock struct {
	CreateRequesterContainerFactoryCalled func(args requesterscontainer.FactoryArgs) (dataRetriever.RequestersContainerFactory, error)
}

// CreateRequesterContainerFactory -
func (r *RequestersContainerFactoryMock) CreateRequesterContainerFactory(args requesterscontainer.FactoryArgs) (dataRetriever.RequestersContainerFactory, error) {
	if r.CreateRequesterContainerFactoryCalled != nil {
		return r.CreateRequesterContainerFactory(args)
	}
	return &dataRetrieverMocks.ShardRequestersContainerFactoryMock{}, nil
}

// IsInterfaceNil checks if underlying pointer is nil
func (r *RequestersContainerFactoryMock) IsInterfaceNil() bool {
	return r == nil
}
