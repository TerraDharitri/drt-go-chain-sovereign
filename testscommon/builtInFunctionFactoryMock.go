package testscommon

import (
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// BuiltInFunctionFactoryMock -
type BuiltInFunctionFactoryMock struct {
	DCDTGlobalSettingsHandlerCalled      func() vmcommon.DCDTGlobalSettingsHandler
	NFTStorageHandlerCalled              func() vmcommon.SimpleDCDTNFTStorageHandler
	BuiltInFunctionContainerCalled       func() vmcommon.BuiltInFunctionContainer
	SetPayableHandlerCalled              func(handler vmcommon.PayableHandler) error
	CreateBuiltInFunctionContainerCalled func() error
}

// DCDTGlobalSettingsHandler -
func (b *BuiltInFunctionFactoryMock) DCDTGlobalSettingsHandler() vmcommon.DCDTGlobalSettingsHandler {
	if b.DCDTGlobalSettingsHandlerCalled != nil {
		return b.DCDTGlobalSettingsHandlerCalled()
	}
	return &DCDTGlobalSettingsHandlerStub{}
}

// NFTStorageHandler -
func (b *BuiltInFunctionFactoryMock) NFTStorageHandler() vmcommon.SimpleDCDTNFTStorageHandler {
	if b.NFTStorageHandlerCalled != nil {
		return b.NFTStorageHandlerCalled()
	}
	return &SimpleNFTStorageHandlerStub{}
}

// BuiltInFunctionContainer -
func (b *BuiltInFunctionFactoryMock) BuiltInFunctionContainer() vmcommon.BuiltInFunctionContainer {
	if b.BuiltInFunctionContainerCalled != nil {
		return b.BuiltInFunctionContainerCalled()
	}
	return &BuiltInFunctionContainerStub{}
}

// SetPayableHandler -
func (b *BuiltInFunctionFactoryMock) SetPayableHandler(handler vmcommon.PayableHandler) error {
	if b.SetPayableHandlerCalled != nil {
		return b.SetPayableHandlerCalled(handler)
	}
	return nil
}

// CreateBuiltInFunctionContainer -
func (b *BuiltInFunctionFactoryMock) CreateBuiltInFunctionContainer() error {
	if b.CreateBuiltInFunctionContainerCalled != nil {
		return b.CreateBuiltInFunctionContainerCalled()
	}
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (b *BuiltInFunctionFactoryMock) IsInterfaceNil() bool {
	return b == nil
}
