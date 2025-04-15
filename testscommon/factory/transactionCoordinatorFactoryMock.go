package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/coordinator"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

// TransactionCoordinatorFactoryMock -
type TransactionCoordinatorFactoryMock struct {
	CreateTransactionCoordinatorCalled func(args coordinator.ArgTransactionCoordinator) (process.TransactionCoordinator, error)
}

// CreateTransactionCoordinator -
func (t *TransactionCoordinatorFactoryMock) CreateTransactionCoordinator(args coordinator.ArgTransactionCoordinator) (process.TransactionCoordinator, error) {
	if t.CreateTransactionCoordinatorCalled != nil {
		return t.CreateTransactionCoordinatorCalled(args)
	}
	return &testscommon.TransactionCoordinatorMock{}, nil
}

// IsInterfaceNil -
func (t *TransactionCoordinatorFactoryMock) IsInterfaceNil() bool {
	return t == nil
}
