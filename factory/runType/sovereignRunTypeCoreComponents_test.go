package runType_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/factory/runType"
)

func TestSovereignRunTypeCoreComponentsFactory_CreateAndClose(t *testing.T) {
	t.Parallel()

	srccf := runType.NewSovereignRunTypeCoreComponentsFactory(config.SovereignEpochConfig{})
	require.NotNil(t, srccf)

	rcc := srccf.Create()
	require.NotNil(t, rcc)

	require.NoError(t, rcc.Close())
}
