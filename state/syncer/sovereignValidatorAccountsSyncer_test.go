package syncer_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/state/syncer"
	"github.com/stretchr/testify/require"
)

func TestNewSovereignValidatorAccountsSyncer(t *testing.T) {
	t.Parallel()

	args := syncer.ArgsNewValidatorAccountsSyncer{
		ArgsNewBaseAccountsSyncer: getDefaultBaseAccSyncerArgs(),
	}

	sovSyncer, err := syncer.NewSovereignValidatorAccountsSyncer(args)
	require.Nil(t, err)
	require.NotNil(t, sovSyncer)
	require.Equal(t, core.SovereignChainShardId, sovSyncer.GetShardID())
}
