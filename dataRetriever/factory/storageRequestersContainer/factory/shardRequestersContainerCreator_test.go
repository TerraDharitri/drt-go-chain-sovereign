package factory

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	"github.com/TerraDharitri/drt-go-chain/common/statistics/disabled"
	storagerequesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/storageRequestersContainer"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/stretchr/testify/require"
)

func createFactoryArgs() storagerequesterscontainer.FactoryArgs {
	return storagerequesterscontainer.FactoryArgs{
		ChainID:                  "T",
		WorkingDirectory:         "",
		Hasher:                   &hashingMocks.HasherMock{},
		ShardCoordinator:         mock.NewOneShardCoordinatorMock(),
		Messenger:                &p2pmocks.MessengerStub{},
		Store:                    &storage.ChainStorerStub{},
		Marshalizer:              &mock.MarshalizerMock{},
		Uint64ByteSliceConverter: &mock.Uint64ByteSliceConverterMock{},
		DataPacker:               &mock.DataPackerStub{},
		ManualEpochStartNotifier: &mock.ManualEpochStartNotifierStub{},
		ChanGracefullyClose:      make(chan endProcess.ArgEndProcess),
		EnableEpochsHandler:      &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		StateStatsHandler:        disabled.NewStateStatistics(),
	}
}

func TestShardRequestersContainerCreator_CreateShardRequestersContainerFactory(t *testing.T) {
	t.Parallel()

	creator := NewShardRequestersContainerCreator()
	require.False(t, creator.IsInterfaceNil())

	args := createFactoryArgs()
	container, err := creator.CreateShardRequestersContainerFactory(args)
	require.Nil(t, err)
	require.Equal(t, "*storagerequesterscontainer.shardRequestersContainerFactory", fmt.Sprintf("%T", container))
}
