package interceptorscontainer_test

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/process/factory/interceptorscontainer"
	"github.com/TerraDharitri/drt-go-chain/testscommon/sovereign"
	"github.com/stretchr/testify/require"
)

func TestNewSovereignShardInterceptorsContainerFactoryCreator(t *testing.T) {
	t.Parallel()

	factoryCreator := interceptorscontainer.NewSovereignShardInterceptorsContainerFactoryCreator()
	require.Implements(t, new(interceptorscontainer.InterceptorsContainerFactoryCreator), factoryCreator)
	require.False(t, check.IfNil(factoryCreator))
}

func TestSovereignShardInterceptorsContainerFactoryCreator_CreateInterceptorsContainerFactory(t *testing.T) {
	t.Parallel()

	factoryCreator := interceptorscontainer.NewSovereignShardInterceptorsContainerFactoryCreator()
	coreComp, cryptoComp := createMockComponentHolders()
	args := getArgumentsShard(coreComp, cryptoComp)
	args.IncomingHeaderSubscriber = &sovereign.IncomingHeaderSubscriberStub{}
	factory, err := factoryCreator.CreateInterceptorsContainerFactory(args)
	require.Nil(t, err)
	require.Equal(t, "*interceptorscontainer.sovereignShardInterceptorsContainerFactory", fmt.Sprintf("%T", factory))
}
