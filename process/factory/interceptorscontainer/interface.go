package interceptorscontainer

import "github.com/TerraDharitri/drt-go-chain/process"

// InterceptorsContainerFactoryCreator defines an interceptor container factory creator
type InterceptorsContainerFactoryCreator interface {
	CreateInterceptorsContainerFactory(args CommonInterceptorsContainerFactoryArgs) (process.InterceptorsContainerFactory, error)
	IsInterfaceNil() bool
}
