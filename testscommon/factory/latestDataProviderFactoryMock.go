package factory

import (
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/storage/latestData"
	storageMock "github.com/TerraDharitri/drt-go-chain/storage/mock"
)

// LatestDataProviderFactoryMock -
type LatestDataProviderFactoryMock struct {
	CreateLatestDataProviderCalled func(args latestData.ArgsLatestDataProvider) (storage.LatestStorageDataProviderHandler, error)
}

// CreateLatestDataProvider -
func (mock *LatestDataProviderFactoryMock) CreateLatestDataProvider(args latestData.ArgsLatestDataProvider) (storage.LatestStorageDataProviderHandler, error) {
	if mock.CreateLatestDataProviderCalled != nil {
		return mock.CreateLatestDataProviderCalled(args)
	}

	return &storageMock.LatestStorageDataProviderStub{}, nil
}

// IsInterfaceNil -
func (mock *LatestDataProviderFactoryMock) IsInterfaceNil() bool {
	return mock == nil
}
