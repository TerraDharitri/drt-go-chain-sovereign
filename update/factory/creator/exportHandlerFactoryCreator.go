package creator

import (
	drtFactory "github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/update"
	"github.com/TerraDharitri/drt-go-chain/update/factory"
)

type exportHandlerFactoryCreator struct {
}

// NewExportHandlerFactoryCreator creates an export handler factory creator
func NewExportHandlerFactoryCreator() *exportHandlerFactoryCreator {
	return &exportHandlerFactoryCreator{}
}

// CreateExportFactoryHandler creates an export factory handler
func (f *exportHandlerFactoryCreator) CreateExportFactoryHandler(args drtFactory.ArgsExporter) (update.ExportFactoryHandler, error) {
	return factory.NewExportHandlerFactory(args)
}

// IsInterfaceNil checks if the underlying pointer is nil
func (f *exportHandlerFactoryCreator) IsInterfaceNil() bool {
	return f == nil
}
