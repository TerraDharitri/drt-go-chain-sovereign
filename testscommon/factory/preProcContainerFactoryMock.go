package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	procMock "github.com/TerraDharitri/drt-go-chain/process/mock"
)

// PreProcessorsContainerFactoryMock -
type PreProcessorsContainerFactoryMock struct {
}

// Create -
func (mock *PreProcessorsContainerFactoryMock) Create() (process.PreProcessorsContainer, error) {
	return &procMock.PreProcessorContainerMock{}, nil
}

// IsInterfaceNil  -
func (mock *PreProcessorsContainerFactoryMock) IsInterfaceNil() bool {
	return mock == nil
}
