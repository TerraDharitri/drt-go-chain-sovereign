package node_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/factory/mock"
	"github.com/TerraDharitri/drt-go-chain/node"
	componentsMock "github.com/TerraDharitri/drt-go-chain/testscommon/components"
	"github.com/TerraDharitri/drt-go-chain/testscommon/consensus/factoryMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/mainFactoryMocks"
)

func TestCreateNode(t *testing.T) {
	t.Parallel()

	t.Run("nil node factory should not work", func(t *testing.T) {
		t.Parallel()

		nodeHandler, err := node.CreateNode(
			&config.Config{},
			componentsMock.GetRunTypeComponents(),
			&factory.StatusCoreComponentsStub{},
			getDefaultBootstrapComponents(),
			getDefaultCoreComponents(),
			getDefaultCryptoComponents(),
			getDefaultDataComponents(),
			getDefaultNetworkComponents(),
			getDefaultProcessComponents(),
			getDefaultStateComponents(),
			&mainFactoryMocks.StatusComponentsStub{},
			&mock.HeartbeatV2ComponentsStub{},
			&factoryMocks.ConsensusComponentsStub{
				GroupSize: 1,
			},
			0,
			false,
			nil)

		require.NotNil(t, err)
		require.Equal(t, errors.ErrNilNode, err)
		require.Nil(t, nodeHandler)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		nodeHandler, err := node.CreateNode(
			&config.Config{},
			componentsMock.GetRunTypeComponents(),
			&factory.StatusCoreComponentsStub{},
			getDefaultBootstrapComponents(),
			getDefaultCoreComponents(),
			getDefaultCryptoComponents(),
			getDefaultDataComponents(),
			getDefaultNetworkComponents(),
			getDefaultProcessComponents(),
			getDefaultStateComponents(),
			&mainFactoryMocks.StatusComponentsStub{},
			&mock.HeartbeatV2ComponentsStub{},
			&factoryMocks.ConsensusComponentsStub{
				GroupSize: 1,
			},
			0,
			false,
			node.NewSovereignNodeFactory(nativeDCDT))

		require.Nil(t, err)
		require.NotNil(t, nodeHandler)
	})
}
