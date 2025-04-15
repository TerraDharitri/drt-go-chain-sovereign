package sharding

import "github.com/TerraDharitri/drt-go-chain-core/core"

// NodesSetupArgs defines arguments needed to create a genesis nodes setup handler
type NodesSetupArgs struct {
	NodesFilePath            string
	AddressPubKeyConverter   core.PubkeyConverter
	ValidatorPubKeyConverter core.PubkeyConverter
	GenesisMaxNumShards      uint32
}
