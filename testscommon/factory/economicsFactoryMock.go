package factory

import (
	metachainEpochStart "github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// EconomicsFactoryMock -
type EconomicsFactoryMock struct {
	CreateEndOfEpochEconomicsCalled func(args metachainEpochStart.ArgsNewEpochEconomics) (process.EndOfEpochEconomics, error)
}

// CreateEndOfEpochEconomics -
func (f *EconomicsFactoryMock) CreateEndOfEpochEconomics(args metachainEpochStart.ArgsNewEpochEconomics) (process.EndOfEpochEconomics, error) {
	if f.CreateEndOfEpochEconomicsCalled != nil {
		return f.CreateEndOfEpochEconomicsCalled(args)
	}

	return &mock.EpochEconomicsStub{}, nil
}

// IsInterfaceNil -
func (f *EconomicsFactoryMock) IsInterfaceNil() bool {
	return f == nil
}
