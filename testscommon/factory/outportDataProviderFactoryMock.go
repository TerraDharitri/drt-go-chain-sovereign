package factory

import (
	"github.com/TerraDharitri/drt-go-chain/outport"
	"github.com/TerraDharitri/drt-go-chain/outport/process/factory"
	outportStub "github.com/TerraDharitri/drt-go-chain/testscommon/outport"
)

// OutportDataProviderFactoryMock -
type OutportDataProviderFactoryMock struct {
	CreateOutportDataProviderCalled func(arg factory.ArgOutportDataProviderFactory) (outport.DataProviderOutport, error)
}

// CreateOutportDataProvider -
func (f *OutportDataProviderFactoryMock) CreateOutportDataProvider(arg factory.ArgOutportDataProviderFactory) (outport.DataProviderOutport, error) {
	if f.CreateOutportDataProviderCalled != nil {
		return f.CreateOutportDataProviderCalled(arg)
	}

	return &outportStub.OutportDataProviderStub{}, nil
}

// IsInterfaceNil -
func (f *OutportDataProviderFactoryMock) IsInterfaceNil() bool {
	return f == nil
}
