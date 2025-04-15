package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/factory/interceptorscontainer"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

// InterceptorsContainerFactoryMock -
type InterceptorsContainerFactoryMock struct {
	CreateInterceptorsContainerFactoryCalled func(args interceptorscontainer.CommonInterceptorsContainerFactoryArgs) (process.InterceptorsContainerFactory, error)
}

// CreateInterceptorsContainerFactory -
func (i *InterceptorsContainerFactoryMock) CreateInterceptorsContainerFactory(args interceptorscontainer.CommonInterceptorsContainerFactoryArgs) (process.InterceptorsContainerFactory, error) {
	if i.CreateInterceptorsContainerFactoryCalled != nil {
		return i.CreateInterceptorsContainerFactory(args)
	}
	return &testscommon.ShardInterceptorsContainerFactoryMock{}, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (i *InterceptorsContainerFactoryMock) IsInterfaceNil() bool {
	return i == nil
}
