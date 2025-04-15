package factory

import (
	"github.com/TerraDharitri/drt-go-chain/node/mock"
	"github.com/TerraDharitri/drt-go-chain/process"
	processBlock "github.com/TerraDharitri/drt-go-chain/process/block"
)

// HeaderValidatorFactoryMock -
type HeaderValidatorFactoryMock struct {
	CreateHeaderValidatorCalled func(args processBlock.ArgsHeaderValidator) (process.HeaderConstructionValidator, error)
}

// CreateHeaderValidator -
func (h *HeaderValidatorFactoryMock) CreateHeaderValidator(args processBlock.ArgsHeaderValidator) (process.HeaderConstructionValidator, error) {
	if h.CreateHeaderValidatorCalled != nil {
		return h.CreateHeaderValidatorCalled(args)
	}
	return &mock.HeaderValidatorStub{}, nil
}

// IsInterfaceNil -
func (h *HeaderValidatorFactoryMock) IsInterfaceNil() bool {
	return h == nil
}
