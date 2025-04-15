package factory

import (
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

// SysSCFactoryMock -
type SysSCFactoryMock struct {
	CreateSystemSCProcessorCalled func(args metachain.ArgsNewEpochStartSystemSCProcessing) (process.EpochStartSystemSCProcessor, error)
}

// CreateSystemSCProcessor -
func (mock *SysSCFactoryMock) CreateSystemSCProcessor(args metachain.ArgsNewEpochStartSystemSCProcessing) (process.EpochStartSystemSCProcessor, error) {
	if mock.CreateSystemSCProcessorCalled != nil {
		return mock.CreateSystemSCProcessorCalled(args)
	}

	return &testscommon.EpochStartSystemSCStub{}, nil
}

// IsInterfaceNil -
func (mock *SysSCFactoryMock) IsInterfaceNil() bool {
	return mock == nil
}
