package enablers

import (
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// EnableEpochsFactory defines enable epochs handler factory behavior
type EnableEpochsFactory interface {
	CreateEnableEpochsHandler(enableEpochs config.EnableEpochs, epochNotifier process.EpochNotifier) (common.EnableEpochsHandler, error)
	IsInterfaceNil() bool
}
