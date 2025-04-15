package factory

import (
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/broadcast"
	consensusMock "github.com/TerraDharitri/drt-go-chain/consensus/mock"
)

// ShardChainMessengerFactoryMock -
type ShardChainMessengerFactoryMock struct {
	CreateShardChainMessengerCalled func(args broadcast.ShardChainMessengerArgs) (consensus.BroadcastMessenger, error)
}

// CreateShardChainMessenger -
func (mock *ShardChainMessengerFactoryMock) CreateShardChainMessenger(args broadcast.ShardChainMessengerArgs) (consensus.BroadcastMessenger, error) {
	if mock.CreateShardChainMessengerCalled != nil {
		return mock.CreateShardChainMessengerCalled(args)
	}

	return &consensusMock.BroadcastMessengerMock{}, nil
}

// IsInterfaceNil -
func (mock *ShardChainMessengerFactoryMock) IsInterfaceNil() bool {
	return mock == nil
}
