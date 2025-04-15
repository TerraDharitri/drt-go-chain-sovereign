package creator

import (
	drtFactory "github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/update"
	"github.com/TerraDharitri/drt-go-chain/update/disabled"
)

type sovereignExportHandlerFactoryCreator struct {
}

// NewSovereignExportHandlerFactoryCreator creates a sovereign export handler factory creator
func NewSovereignExportHandlerFactoryCreator() *sovereignExportHandlerFactoryCreator {
	return &sovereignExportHandlerFactoryCreator{}
}

// CreateExportFactoryHandler creates a disabled export factory handler
func (f *sovereignExportHandlerFactoryCreator) CreateExportFactoryHandler(_ drtFactory.ArgsExporter) (update.ExportFactoryHandler, error) {
	return &disabled.ExportFactoryHandler{}, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (f *sovereignExportHandlerFactoryCreator) IsInterfaceNil() bool {
	return f == nil
}
