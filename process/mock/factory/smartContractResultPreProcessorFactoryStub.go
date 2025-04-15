package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/preprocess"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
)

// SmartContractResultPreProcessorFactoryStub -
type SmartContractResultPreProcessorFactoryStub struct {
}

// CreateSmartContractResultPreProcessor -
func (s *SmartContractResultPreProcessorFactoryStub) CreateSmartContractResultPreProcessor(_ preprocess.SmartContractResultPreProcessorCreatorArgs) (process.PreProcessor, error) {
	return &mock.PreProcessorMock{}, nil
}

// IsInterfaceNil -
func (s *SmartContractResultPreProcessorFactoryStub) IsInterfaceNil() bool {
	return s == nil
}
