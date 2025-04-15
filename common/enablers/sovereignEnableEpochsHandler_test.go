package enablers

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/testscommon/epochNotifier"
)

func TestNewSovereignEnableEpochsHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil epoch notifier should error", func(t *testing.T) {
		t.Parallel()

		sovHandler, err := NewSovereignEnableEpochsHandler(createEnableEpochsConfig(), config.SovereignEpochConfig{}, nil)
		require.Equal(t, process.ErrNilEpochNotifier, err)
		require.True(t, sovHandler.IsInterfaceNil())
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		sovHandler, err := NewSovereignEnableEpochsHandler(createEnableEpochsConfig(), config.SovereignEpochConfig{}, &epochNotifier.EpochNotifierStub{
			RegisterNotifyHandlerCalled: func(handler vmcommon.EpochSubscriberHandler) {
				wasCalled = true
			},
		})
		require.Nil(t, err)
		require.False(t, check.IfNil(sovHandler))
		require.True(t, wasCalled)
	})
}
