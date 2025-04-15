package factory_test

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	trieIteratorsFactory "github.com/TerraDharitri/drt-go-chain/node/trieIterators/factory"
	"github.com/stretchr/testify/require"
)

func TestNewSovereignDelegatedListProcessorFactory(t *testing.T) {
	t.Parallel()

	sovereignDelegatedListHandlerFactory := trieIteratorsFactory.NewSovereignDelegatedListProcessorFactory()
	require.False(t, sovereignDelegatedListHandlerFactory.IsInterfaceNil())
}

func TestSovereignDelegatedListProcessorFactory_CreateDelegatedListProcessorHandler(t *testing.T) {
	t.Parallel()

	args := createMockArgs(core.SovereignChainShardId)

	sovereignDelegatedListHandler, err := trieIteratorsFactory.NewSovereignDelegatedListProcessorFactory().CreateDelegatedListProcessorHandler(args)
	require.Nil(t, err)
	require.Equal(t, "*trieIterators.delegatedListProcessor", fmt.Sprintf("%T", sovereignDelegatedListHandler))
}
