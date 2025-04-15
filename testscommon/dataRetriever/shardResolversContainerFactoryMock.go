package dataRetriever

import (
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
)

// ShardResolversContainerFactoryMock -
type ShardResolversContainerFactoryMock struct {
	CreateCalled func() (dataRetriever.ResolversContainer, error)
}

// Create -
func (s *ShardResolversContainerFactoryMock) Create() (dataRetriever.ResolversContainer, error) {
	if s.CreateCalled != nil {
		return s.Create()
	}
	return &ResolversContainerStub{}, nil
}

// IsInterfaceNil -
func (s *ShardResolversContainerFactoryMock) IsInterfaceNil() bool {
	return s == nil
}
