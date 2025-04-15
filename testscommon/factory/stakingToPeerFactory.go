package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	processMock "github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/process/scToProtocol"
)

// StakingToPeerFactoryMock -
type StakingToPeerFactoryMock struct {
	CreateStakingToPeerCalled func(args scToProtocol.ArgStakingToPeer) (process.SmartContractToProtocolHandler, error)
}

// CreateStakingToPeer -
func (mock *StakingToPeerFactoryMock) CreateStakingToPeer(args scToProtocol.ArgStakingToPeer) (process.SmartContractToProtocolHandler, error) {
	if mock.CreateStakingToPeerCalled != nil {
		return mock.CreateStakingToPeerCalled(args)
	}

	return &processMock.SCToProtocolStub{}, nil
}

// IsInterfaceNil -
func (mock *StakingToPeerFactoryMock) IsInterfaceNil() bool {
	return mock == nil
}
