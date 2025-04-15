package persister

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

// PersistentStatusHandler defines a persistent status handler
type PersistentStatusHandler interface {
	core.AppStatusHandler
	SetStorage(store storage.Storer) error
}
