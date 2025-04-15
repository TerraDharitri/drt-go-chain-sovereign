package factory

import (
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

// ValidatorInfoCreatorFactoryMock -
type ValidatorInfoCreatorFactoryMock struct {
	CreateValidatorInfoCreatorCalled func(args metachain.ArgsNewValidatorInfoCreator) (process.EpochStartValidatorInfoCreator, error)
}

// CreateValidatorInfoCreator -
func (mock *ValidatorInfoCreatorFactoryMock) CreateValidatorInfoCreator(args metachain.ArgsNewValidatorInfoCreator) (process.EpochStartValidatorInfoCreator, error) {
	if mock.CreateValidatorInfoCreatorCalled != nil {
		return mock.CreateValidatorInfoCreatorCalled(args)
	}

	return &testscommon.EpochValidatorInfoCreatorStub{}, nil
}

// IsInterfaceNil -
func (mock *ValidatorInfoCreatorFactoryMock) IsInterfaceNil() bool {
	return mock == nil
}
