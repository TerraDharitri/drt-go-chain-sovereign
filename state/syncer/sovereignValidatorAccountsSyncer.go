package syncer

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// NewSovereignValidatorAccountsSyncer creates a validator account syncer for sovereign
func NewSovereignValidatorAccountsSyncer(args ArgsNewValidatorAccountsSyncer) (*validatorAccountsSyncer, error) {
	return newValidatorAccountsSyncer(args, core.SovereignChainShardId)
}
