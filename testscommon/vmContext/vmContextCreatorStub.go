package vmContext

import (
	"github.com/TerraDharitri/drt-go-chain/vm"
	"github.com/TerraDharitri/drt-go-chain/vm/mock"
	"github.com/TerraDharitri/drt-go-chain/vm/systemSmartContracts"
)

// VMContextCreatorStub -
type VMContextCreatorStub struct {
	CreateVmContextCalled func(args systemSmartContracts.VMContextArgs) (vm.ContextHandler, error)
}

// CreateVmContext -
func (stub *VMContextCreatorStub) CreateVmContext(args systemSmartContracts.VMContextArgs) (vm.ContextHandler, error) {
	if stub.CreateVmContextCalled != nil {
		return stub.CreateVmContextCalled(args)
	}

	return &mock.SystemEIStub{}, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (stub *VMContextCreatorStub) IsInterfaceNil() bool {
	return stub == nil
}
