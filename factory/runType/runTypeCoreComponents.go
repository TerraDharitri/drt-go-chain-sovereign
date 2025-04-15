package runType

import (
	"github.com/TerraDharitri/drt-go-chain/common/enablers"
	"github.com/TerraDharitri/drt-go-chain/process/rating"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

type runTypeCoreComponents struct {
	genesisNodesSetupFactory sharding.GenesisNodesSetupFactory
	ratingsDataFactory       rating.RatingsDataFactory
	enableEpochsFactory      enablers.EnableEpochsFactory
}

// Close does nothing
func (rcc *runTypeCoreComponents) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (rcc *runTypeCoreComponents) IsInterfaceNil() bool {
	return rcc == nil
}
