package storageBootstrap

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/epochNotifier"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genericMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/stretchr/testify/require"
)

func TestNewShardStorageBootstrapperFactory(t *testing.T) {
	t.Parallel()

	sbf := NewShardStorageBootstrapperFactory()
	require.NotNil(t, sbf)
}

func TestShardStorageBootstrapperFactory_CreateShardStorageBootstrapper(t *testing.T) {
	t.Parallel()

	sbf := NewShardStorageBootstrapperFactory()
	bootStrapper, err := sbf.CreateBootstrapperFromStorage(getDefaultArgShardBootstrapper())

	require.NotNil(t, bootStrapper)
	require.Nil(t, err)
}

func TestShardStorageBootstrapperFactory_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	sbf := NewShardStorageBootstrapperFactory()
	require.False(t, sbf.IsInterfaceNil())

	sbf = nil
	require.True(t, sbf.IsInterfaceNil())
}

func getDefaultArgShardBootstrapper() ArgsShardStorageBootstrapper {
	bootStorer := genericMocks.NewStorerMock()
	argBaseBoostrapper := ArgsShardStorageBootstrapper{
		ArgsBaseStorageBootstrapper{
			BootStorer:     &mock.BoostrapStorerMock{},
			ForkDetector:   &mock.ForkDetectorMock{},
			BlockProcessor: &testscommon.BlockProcessorStub{},
			ChainHandler:   &testscommon.ChainHandlerStub{},
			Marshalizer:    &testscommon.ProtoMarshalizerMock{},
			Store: &storageStubs.ChainStorerStub{
				GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
					return bootStorer, nil
				},
			},
			Uint64Converter:              &mock.Uint64ByteSliceConverterMock{},
			BootstrapRoundIndex:          0,
			ShardCoordinator:             &testscommon.ShardsCoordinatorMock{},
			NodesCoordinator:             &shardingMocks.NodesCoordinatorMock{},
			EpochStartTrigger:            &testscommon.EpochStartTriggerStub{},
			BlockTracker:                 &mock.BlockTrackerMock{},
			ChainID:                      "2",
			ScheduledTxsExecutionHandler: &testscommon.ScheduledTxsExecutionStub{},
			MiniblocksProvider:           &mock.MiniBlocksProviderStub{},
			EpochNotifier:                &epochNotifier.EpochNotifierStub{},
			ProcessedMiniBlocksTracker:   &testscommon.ProcessedMiniBlocksTrackerStub{},
			AppStatusHandler:             statusHandler.NewAppStatusHandlerMock(),
		},
	}

	return argBaseBoostrapper
}
