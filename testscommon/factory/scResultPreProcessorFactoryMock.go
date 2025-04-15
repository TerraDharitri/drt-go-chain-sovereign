package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/preprocess"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
)

// SmartContractResultPreProcessorFactoryMock -
type SmartContractResultPreProcessorFactoryMock struct {
	CreateSmartContractResultPreProcessorCalled func(args preprocess.SmartContractResultPreProcessorCreatorArgs) (process.PreProcessor, error)
}

// CreateSmartContractResultPreProcessor -
func (s *SmartContractResultPreProcessorFactoryMock) CreateSmartContractResultPreProcessor(args preprocess.SmartContractResultPreProcessorCreatorArgs) (process.PreProcessor, error) {
	if s.CreateSmartContractResultPreProcessorCalled != nil {
		return s.CreateSmartContractResultPreProcessorCalled(args)
	}
	return &mock.PreProcessorMock{}, nil
}

// IsInterfaceNil -
func (s *SmartContractResultPreProcessorFactoryMock) IsInterfaceNil() bool {
	return s == nil
}
