package node_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/node"
)

var nativeDCDT = "WREWA-bd4d79"

func TestNewSovereignNode(t *testing.T) {
	t.Parallel()

	t.Run("valid node should work", func(t *testing.T) {
		t.Parallel()

		n, err := node.NewNode()
		require.Nil(t, err)
		require.NotNil(t, n)

		sn, err := node.NewSovereignNode(n, nativeDCDT)
		require.Nil(t, err)
		require.False(t, sn.IsInterfaceNil())
	})
	t.Run("nil node should error", func(t *testing.T) {
		t.Parallel()

		sn, err := node.NewSovereignNode(nil, nativeDCDT)
		require.NotNil(t, err)
		require.Equal(t, errors.ErrNilNode, err)
		require.Nil(t, sn)
	})
	t.Run("empty native dcdt should error", func(t *testing.T) {
		t.Parallel()

		n, err := node.NewNode()
		require.Nil(t, err)
		require.NotNil(t, n)

		sn, err := node.NewSovereignNode(n, "")
		require.NotNil(t, err)
		require.Equal(t, node.ErrEmptyNativeDcdt, err)
		require.Nil(t, sn)
	})
}

func TestSovereignNode_GetAllDCDTTokens(t *testing.T) {
	t.Parallel()

	testNodeGetAllIssuedDCDTs(t, node.NewSovereignNodeFactory(nativeDCDT), core.SovereignChainShardId)
	testNodeGetAllIssuedDCDTs(t, node.NewSovereignNodeFactory(nativeDCDT), core.MainChainShardId)
}

func TestSovereignNode_GetNFTTokenIDsRegisteredByAddress(t *testing.T) {
	t.Parallel()

	testNodeGetNFTTokenIDsRegisteredByAddress(t, node.NewSovereignNodeFactory(nativeDCDT), core.SovereignChainShardId)
	testNodeGetNFTTokenIDsRegisteredByAddress(t, node.NewSovereignNodeFactory(nativeDCDT), core.MainChainShardId)
}

func TestSovereignNode_GetDCDTsWithRole(t *testing.T) {
	t.Parallel()

	testNodeGetDCDTsWithRole(t, node.NewSovereignNodeFactory(nativeDCDT), core.SovereignChainShardId)
	testNodeGetDCDTsWithRole(t, node.NewSovereignNodeFactory(nativeDCDT), core.MainChainShardId)
}

func TestSovereignNode_GetDCDTsRoles(t *testing.T) {
	t.Parallel()

	testNodeGetDCDTsRoles(t, node.NewSovereignNodeFactory(nativeDCDT), core.SovereignChainShardId)
	testNodeGetDCDTsRoles(t, node.NewSovereignNodeFactory(nativeDCDT), core.MainChainShardId)
}
