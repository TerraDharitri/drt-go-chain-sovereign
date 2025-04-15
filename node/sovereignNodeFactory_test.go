package node_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/node"
)

func TestNewSovereignNodeFactory(t *testing.T) {
	t.Parallel()

	sovereignNodeFactory := node.NewSovereignNodeFactory(nativeDCDT)
	require.False(t, sovereignNodeFactory.IsInterfaceNil())
}

func TestSovereignNodeFactory_CreateNewNode(t *testing.T) {
	t.Parallel()

	sovereignNodeFactory := node.NewSovereignNodeFactory(nativeDCDT)

	sn, err := sovereignNodeFactory.CreateNewNode()
	require.Nil(t, err)
	require.NotNil(t, sn)
	require.Equal(t, "*node.sovereignNode", fmt.Sprintf("%T", sn))
}

func TestSovereignNodeFactory_CreateNewNodeFail(t *testing.T) {
	t.Parallel()

	sovereignNodeFactory := node.NewSovereignNodeFactory(nativeDCDT)

	options := []node.Option{
		node.WithStatusCoreComponents(nil),
	}

	sn, err := sovereignNodeFactory.CreateNewNode(options...)
	require.Contains(t, err.Error(), node.ErrNilStatusCoreComponents.Error())
	require.Nil(t, sn)
}
