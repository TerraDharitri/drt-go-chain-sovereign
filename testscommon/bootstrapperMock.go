package testscommon

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/sync/disabled"
)

// BootstrapperMock -
type BootstrapperMock struct {
	process.Bootstrapper
}

// NewTestBootstrapperMock -
func NewTestBootstrapperMock() *BootstrapperMock {
	return &BootstrapperMock{
		Bootstrapper: disabled.NewDisabledBootstrapper(),
	}
}

// RollBack -
func (tbm *BootstrapperMock) RollBack(_ bool) error {
	return nil
}

// SetProbableHighestNonce -
func (tbm *BootstrapperMock) SetProbableHighestNonce(_ uint64) {
}
