package storageBootstrap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/process/sync"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	testsDataRetriever "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dblookupext"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genericMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/outport"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	statusHandlerMock "github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
)

func TestNewShardBootstrapFactory(t *testing.T) {
	t.Parallel()

	sbf := NewShardBootstrapFactory()
	require.NotNil(t, sbf)
	require.False(t, sbf.IsInterfaceNil())
}

func TestShardBootstrapFactory_CreateShardBootstrapFactory(t *testing.T) {
	t.Parallel()

	sbf := NewShardBootstrapFactory()
	bootStrapper, err := sbf.CreateBootstrapper(getDefaultArgs())

	require.NotNil(t, bootStrapper)
	require.Nil(t, err)
}

func getDefaultArgs() sync.ArgShardBootstrapper {
	bootStorer := genericMocks.NewStorerMock()
	argBaseBoostrapper := sync.ArgBaseBootstrapper{
		PoolsHolder: testsDataRetriever.NewPoolsHolderMock(),
		Store: &storageStubs.ChainStorerStub{
			GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
				return bootStorer, nil
			},
		},
		ChainHandler:                 &testscommon.ChainHandlerStub{},
		RoundHandler:                 &testscommon.RoundHandlerMock{},
		BlockProcessor:               &testscommon.BlockProcessorStub{},
		WaitTime:                     time.Second,
		Hasher:                       &hashingMocks.HasherMock{},
		Marshalizer:                  &testscommon.ProtoMarshalizerMock{},
		ForkDetector:                 &mock.ForkDetectorMock{},
		RequestHandler:               &testscommon.RequestHandlerStub{},
		ShardCoordinator:             &testscommon.ShardsCoordinatorMock{},
		Accounts:                     &stateMock.AccountsStub{},
		BlackListHandler:             &testscommon.TimeCacheStub{},
		NetworkWatcher:               &mock.NetworkConnectionWatcherStub{},
		BootStorer:                   &mock.BoostrapStorerMock{},
		StorageBootstrapper:          &mock.StorageBootstrapperMock{},
		EpochHandler:                 &mock.EpochStartTriggerStub{},
		MiniblocksProvider:           &mock.MiniBlocksProviderStub{},
		Uint64Converter:              &mock.Uint64ByteSliceConverterMock{},
		AppStatusHandler:             &statusHandlerMock.AppStatusHandlerStub{},
		OutportHandler:               &outport.OutportStub{},
		AccountsDBSyncer:             &mock.AccountsDBSyncerStub{},
		CurrentEpochProvider:         &testscommon.CurrentEpochProviderStub{},
		HistoryRepo:                  &dblookupext.HistoryRepositoryStub{},
		ScheduledTxsExecutionHandler: &testscommon.ScheduledTxsExecutionStub{},
		ProcessWaitTime:              time.Second,
		RepopulateTokensSupplies:     false,
		ValidatorDBSyncer:            &mock.AccountsDBSyncerStub{},
	}

	return sync.ArgShardBootstrapper{
		ArgBaseBootstrapper: argBaseBoostrapper,
	}
}
