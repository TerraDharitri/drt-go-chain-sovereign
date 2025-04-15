package mock

import "github.com/TerraDharitri/drt-go-chain/process"

// InterceptedDataVerifierMock -
type InterceptedDataVerifierMock struct {
}

// IsForCurrentShard -
func (i *InterceptedDataVerifierMock) IsForCurrentShard(_ process.InterceptedData) bool {
	return true
}

// IsInterfaceNil returns true if underlying object is
func (i *InterceptedDataVerifierMock) IsInterfaceNil() bool {
	return i == nil
}
