package process

import (
	"github.com/TerraDharitri/drt-go-chain/api/shared"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/dtos"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"

	chainData "github.com/TerraDharitri/drt-go-chain-core/data"
)

// NodeHandler defines what a node handler should be able to do
type NodeHandler interface {
	GetProcessComponents() factory.ProcessComponentsHolder
	GetChainHandler() chainData.ChainHandler
	GetBroadcastMessenger() consensus.BroadcastMessenger
	GetShardCoordinator() sharding.Coordinator
	GetCryptoComponents() factory.CryptoComponentsHolder
	GetCoreComponents() factory.CoreComponentsHolder
	GetDataComponents() factory.DataComponentsHolder
	GetStateComponents() factory.StateComponentsHolder
	GetFacadeHandler() shared.FacadeHandler
	GetStatusCoreComponents() factory.StatusCoreComponentsHolder
	GetRunTypeComponents() factory.RunTypeComponentsHolder
	GetIncomingHeaderSubscriber() process.IncomingHeaderSubscriber
	SetKeyValueForAddress(addressBytes []byte, state map[string]string) error
	SetStateForAddress(address []byte, state *dtos.AddressState) error
	RemoveAccount(address []byte) error
	ForceChangeOfEpoch() error
	Close() error
	IsInterfaceNil() bool
}

// BlocksProcessorFactory defines what the block processor factory should be able to do
type BlocksProcessorFactory interface {
	ProcessBlock(processor process.BlockProcessor, header chainData.HeaderHandler) (chainData.HeaderHandler, chainData.BodyHandler, error)
	IsInterfaceNil() bool
}
