package storageBootstrap

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/process/block/bootstrapStorage"
)

func (st *storageBootstrapper) sovereignChainGetScheduledRootHash(headerFromStorage data.HeaderHandler, _ []byte) []byte {
	return headerFromStorage.GetRootHash()
}

func (st *storageBootstrapper) sovereignChainSetScheduledInfo(_ bootstrapStorage.BootstrapData) {
}
