package scToProtocol

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/process"
)

type modifiedMBShardIDCheckerHandler interface {
	isModifiedStateMBValid(miniBlock *block.MiniBlock) bool
}

// StakingToPeerFactoryHandler defines the factory interface to create sc to protocol handler
type StakingToPeerFactoryHandler interface {
	CreateStakingToPeer(args ArgStakingToPeer) (process.SmartContractToProtocolHandler, error)
	IsInterfaceNil() bool
}
