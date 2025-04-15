package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/sync/storageBootstrap"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

// BootstrapperFromStorageFactoryMock -
type BootstrapperFromStorageFactoryMock struct {
	CreateBootstrapperFromStorageCalled func(args storageBootstrap.ArgsShardStorageBootstrapper) (process.BootstrapperFromStorage, error)
}

// CreateBootstrapperFromStorage -
func (b *BootstrapperFromStorageFactoryMock) CreateBootstrapperFromStorage(args storageBootstrap.ArgsShardStorageBootstrapper) (process.BootstrapperFromStorage, error) {
	if b.CreateBootstrapperFromStorageCalled != nil {
		return b.CreateBootstrapperFromStorageCalled(args)
	}
	return &testscommon.StorageBootstrapperMock{}, nil
}

// IsInterfaceNil -
func (b *BootstrapperFromStorageFactoryMock) IsInterfaceNil() bool {
	return b == nil
}
