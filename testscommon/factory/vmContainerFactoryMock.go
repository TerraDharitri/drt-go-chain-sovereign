package factory

import (
	factoryVm "github.com/TerraDharitri/drt-go-chain/factory/vm"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/mock"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// VMContainerFactoryMock -
type VMContainerFactoryMock struct {
	CreateVmContainerFactoryCalled func(args factoryVm.ArgsVmContainerFactory) (process.VirtualMachinesContainer, process.VirtualMachinesContainerFactory, error)
}

// CreateVmContainerFactory -
func (v *VMContainerFactoryMock) CreateVmContainerFactory(args factoryVm.ArgsVmContainerFactory) (process.VirtualMachinesContainer, process.VirtualMachinesContainerFactory, error) {
	if v.CreateVmContainerFactoryCalled != nil {
		return v.CreateVmContainerFactoryCalled(args)
	}
	return &mock.VMContainerMock{}, &mock.VmMachinesContainerFactoryMock{}, nil
}

// IsInterfaceNil -
func (v *VMContainerFactoryMock) IsInterfaceNil() bool {
	return v == nil
}
