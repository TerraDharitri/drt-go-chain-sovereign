package epochStartTrigger

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	retriever "github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/factory"
	nodeFactoryMock "github.com/TerraDharitri/drt-go-chain/node/mock/factory"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	shardingMock "github.com/TerraDharitri/drt-go-chain/sharding/mock"
	chainStorage "github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/bootstrapMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	testsFactory "github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genesisMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/mainFactoryMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/pool"
	"github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
	"github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	validatorInfoCacherStub "github.com/TerraDharitri/drt-go-chain/testscommon/validatorInfoCacher"
	updateMock "github.com/TerraDharitri/drt-go-chain/update/mock"
)

func createArgs(shardID uint32) factory.ArgsEpochStartTrigger {
	return factory.ArgsEpochStartTrigger{
		RequestHandler: &testscommon.RequestHandlerStub{},
		CoreData: &testsFactory.CoreComponentsHolderMock{
			HasherCalled: func() hashing.Hasher {
				return &testscommon.HasherStub{}
			},
			InternalMarshalizerCalled: func() marshal.Marshalizer {
				return &testscommon.MarshallerStub{}
			},
			Uint64ByteSliceConverterCalled: func() typeConverters.Uint64ByteSliceConverter {
				return &testscommon.Uint64ByteSliceConverterStub{}
			},
			EpochStartNotifierWithConfirmCalled: func() factory.EpochStartNotifierWithConfirm {
				return &updateMock.EpochStartNotifierStub{}
			},
			RoundHandlerCalled: func() consensus.RoundHandler {
				return &testscommon.RoundHandlerMock{}
			},
			EnableEpochsHandlerCalled: func() common.EnableEpochsHandler {
				return &shardingMock.EnableEpochsHandlerMock{}
			},
			GenesisNodesSetupCalled: func() sharding.GenesisNodesSetupHandler {
				return &genesisMocks.NodesSetupStub{}
			},
		},
		BootstrapComponents: createBootstrapComps(shardID),
		DataComps:           createDataCompsMock(),
		StatusCoreComponentsHolder: &testsFactory.StatusCoreComponentsStub{
			AppStatusHandlerField: &statusHandler.AppStatusHandlerStub{},
		},
		RunTypeComponentsHolder: mainFactoryMocks.NewRunTypeComponentsStub(),
		Config: config.Config{
			EpochStartConfig: config.EpochStartConfig{
				RoundsPerEpoch:         22,
				MinRoundsBetweenEpochs: 22,
			},
		},
	}
}

func createBootstrapComps(shardID uint32) *mainFactoryMocks.BootstrapComponentsStub {
	return &mainFactoryMocks.BootstrapComponentsStub{
		ShardCoordinatorCalled: func() sharding.Coordinator {
			return &testscommon.ShardsCoordinatorMock{
				NoShards: 1,
				SelfIDCalled: func() uint32 {
					return shardID
				},
			}
		},
		BootstrapParams:      &bootstrapMocks.BootstrapParamsHandlerMock{},
		HdrIntegrityVerifier: &mock.HeaderIntegrityVerifierStub{},
		Bootstrapper:         &bootstrapMocks.EpochStartBootstrapperStub{},
	}
}

func createDataCompsMock() *nodeFactoryMock.DataComponentsMock {
	return &nodeFactoryMock.DataComponentsMock{
		DataPool: createDataPoolMock(),
		Store: &storage.ChainStorerStub{
			GetStorerCalled: func(unitType retriever.UnitType) (chainStorage.Storer, error) {
				return &storage.StorerStub{}, nil
			},
		},
		BlockChain: &testscommon.ChainHandlerStub{
			GetGenesisHeaderCalled: func() data.HeaderHandler {
				return &block.HeaderV2{}
			},
		},
	}
}

func createDataPoolMock() *dataRetriever.PoolsHolderStub {
	return &dataRetriever.PoolsHolderStub{
		MetaBlocksCalled: func() chainStorage.Cacher {
			return &testscommon.CacherStub{}
		},
		HeadersCalled: func() retriever.HeadersPool {
			return &pool.HeadersPoolStub{}
		},
		ValidatorsInfoCalled: func() retriever.ShardedDataCacherNotifier {
			return &testscommon.ShardedDataCacheNotifierMock{}
		},
		CurrEpochValidatorInfoCalled: func() retriever.ValidatorInfoCacher {
			return &validatorInfoCacherStub.ValidatorInfoCacherStub{}
		},
	}
}

func TestNewEpochStartTriggerFactory(t *testing.T) {
	t.Parallel()

	f := NewEpochStartTriggerFactory()
	require.False(t, f.IsInterfaceNil())

	t.Run("create for shard", func(t *testing.T) {
		args := createArgs(0)
		trigger, err := f.CreateEpochStartTrigger(args)
		require.Nil(t, err)
		require.Equal(t, "*shardchain.trigger", fmt.Sprintf("%T", trigger))
	})
	t.Run("create for meta", func(t *testing.T) {
		args := createArgs(core.MetachainShardId)
		trigger, err := f.CreateEpochStartTrigger(args)
		require.Nil(t, err)
		require.Equal(t, "*metachain.trigger", fmt.Sprintf("%T", trigger))
	})
	t.Run("invalid shard id", func(t *testing.T) {
		args := createArgs(444)
		trigger, err := f.CreateEpochStartTrigger(args)
		require.ErrorIs(t, err, process.ErrInvalidShardId)
		require.Nil(t, trigger)
	})
}
