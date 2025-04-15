package factory

import (
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/requestHandlers"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

// RequestHandlerFactoryMock -
type RequestHandlerFactoryMock struct {
	CreateRequestHandlerCalled func(args requestHandlers.RequestHandlerArgs) (process.RequestHandler, error)
}

// CreateRequestHandler -
func (r *RequestHandlerFactoryMock) CreateRequestHandler(args requestHandlers.RequestHandlerArgs) (process.RequestHandler, error) {
	if r.CreateRequestHandlerCalled != nil {
		return r.CreateRequestHandlerCalled(args)
	}
	return &testscommon.ExtendedShardHeaderRequestHandlerStub{}, nil
}

// IsInterfaceNil -
func (r *RequestHandlerFactoryMock) IsInterfaceNil() bool {
	return r == nil
}
