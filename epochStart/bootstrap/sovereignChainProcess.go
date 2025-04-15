package bootstrap

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/errors"
)

type sovereignChainEpochStartBootstrap struct {
	*epochStartBootstrap
}

// NewSovereignChainEpochStartBootstrap creates a new instance of sovereignChainEpochStartBootstrap
func NewSovereignChainEpochStartBootstrap(epochStartBootstrap *epochStartBootstrap) (*sovereignChainEpochStartBootstrap, error) {
	if epochStartBootstrap == nil {
		return nil, errors.ErrNilEpochStartBootstrapper
	}

	scesb := &sovereignChainEpochStartBootstrap{
		epochStartBootstrap,
	}

	scesb.bootStrapShardProcessor = &sovereignBootStrapShardProcessor{
		scesb,
	}

	scesb.shardForLatestEpochComputer = scesb
	return scesb, nil
}

// GetShardIDForLatestEpoch returns the shard ID for the latest epoch
func (scesb *sovereignChainEpochStartBootstrap) GetShardIDForLatestEpoch() (uint32, bool, error) {
	return core.SovereignChainShardId, false, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (scesb *sovereignChainEpochStartBootstrap) IsInterfaceNil() bool {
	return scesb == nil
}
