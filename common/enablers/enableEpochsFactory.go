package enablers

import (
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/process"
)

type enableEpochsFactory struct{}

// NewEnableEpochsFactory creates an enable epochs factory for regular chain
func NewEnableEpochsFactory() EnableEpochsFactory {
	return &enableEpochsFactory{}
}

// CreateEnableEpochsHandler creates an enable epochs handler for regular chain
func (eef *enableEpochsFactory) CreateEnableEpochsHandler(enableEpochs config.EnableEpochs, epochNotifier process.EpochNotifier) (common.EnableEpochsHandler, error) {
	return NewEnableEpochsHandler(enableEpochs, epochNotifier)
}

// IsInterfaceNil checks if the underlying pointer is nil
func (eef *enableEpochsFactory) IsInterfaceNil() bool {
	return eef == nil
}
