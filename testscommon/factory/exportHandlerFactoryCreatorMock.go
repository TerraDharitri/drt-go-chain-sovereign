package factory

import (
	drtFactory "github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/update"
	updateMock "github.com/TerraDharitri/drt-go-chain/update/mock"
)

// ExportHandlerFactoryCreatorMock -
type ExportHandlerFactoryCreatorMock struct {
	CreateExportFactoryHandlerCalled func(args drtFactory.ArgsExporter) (update.ExportFactoryHandler, error)
}

// CreateExportFactoryHandler -
func (mock *ExportHandlerFactoryCreatorMock) CreateExportFactoryHandler(args drtFactory.ArgsExporter) (update.ExportFactoryHandler, error) {
	if mock.CreateExportFactoryHandlerCalled != nil {
		return mock.CreateExportFactoryHandlerCalled(args)
	}

	return &updateMock.ExportFactoryHandlerStub{}, nil
}

// IsInterfaceNil -
func (mock *ExportHandlerFactoryCreatorMock) IsInterfaceNil() bool {
	return mock == nil
}
