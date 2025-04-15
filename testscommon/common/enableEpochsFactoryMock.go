package common

import (
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding/mock"
)

// EnableEpochsFactoryMock -
type EnableEpochsFactoryMock struct {
	CreateEnableEpochsHandlerCaller func(epochConfig config.EpochConfig, epochNotifier process.EpochNotifier) (common.EnableEpochsHandler, error)
}

// CreateEnableEpochsHandler -
func (f *EnableEpochsFactoryMock) CreateEnableEpochsHandler(epochConfig config.EpochConfig, epochNotifier process.EpochNotifier) (common.EnableEpochsHandler, error) {
	if f.CreateEnableEpochsHandlerCaller != nil {
		return f.CreateEnableEpochsHandlerCaller(epochConfig, epochNotifier)
	}

	return &mock.EnableEpochsHandlerMock{}, nil
}

// IsInterfaceNil -
func (f *EnableEpochsFactoryMock) IsInterfaceNil() bool {
	return f == nil
}
