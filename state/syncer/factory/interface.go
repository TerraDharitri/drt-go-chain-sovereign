package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/state/syncer"
)

// ValidatorAccountsSyncerFactoryHandler defines a factory able to create a validator accounts db syncer
type ValidatorAccountsSyncerFactoryHandler interface {
	CreateValidatorAccountsSyncer(args syncer.ArgsNewValidatorAccountsSyncer) (process.AccountsDBSyncer, error)
	IsInterfaceNil() bool
}
