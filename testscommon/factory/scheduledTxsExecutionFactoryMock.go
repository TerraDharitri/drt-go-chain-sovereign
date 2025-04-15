package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/preprocess"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

// ScheduledTxsExecutionFactoryMock -
type ScheduledTxsExecutionFactoryMock struct {
	CreateScheduledTxsExecutionHandlerCalled func(args preprocess.ScheduledTxsExecutionFactoryArgs) (process.ScheduledTxsExecutionHandler, error)
}

// CreateScheduledTxsExecutionHandler -
func (s *ScheduledTxsExecutionFactoryMock) CreateScheduledTxsExecutionHandler(args preprocess.ScheduledTxsExecutionFactoryArgs) (process.ScheduledTxsExecutionHandler, error) {
	if s.CreateScheduledTxsExecutionHandlerCalled != nil {
		return s.CreateScheduledTxsExecutionHandlerCalled(args)
	}
	return &testscommon.ScheduledTxsExecutionStub{}, nil
}

// IsInterfaceNil -
func (s *ScheduledTxsExecutionFactoryMock) IsInterfaceNil() bool {
	return s == nil
}
