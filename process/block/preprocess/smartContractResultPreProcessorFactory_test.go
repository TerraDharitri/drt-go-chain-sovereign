package preprocess_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/data/smartContractResult"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/preprocess"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/common"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
)

func TestNewSmartContractResultPreProcessorFactory(t *testing.T) {
	t.Parallel()

	fact := preprocess.NewSmartContractResultPreProcessorFactory()
	require.NotNil(t, fact)
	require.Implements(t, new(preprocess.SmartContractResultPreProcessorCreator), fact)
}

func TestSmartContractResultPreProcessorFactory_CreateSmartContractResultPreProcessor(t *testing.T) {
	t.Parallel()

	fact := preprocess.NewSmartContractResultPreProcessorFactory()

	args := preprocess.SmartContractResultPreProcessorCreatorArgs{}
	preProcessor, err := fact.CreateSmartContractResultPreProcessor(args)
	require.NotNil(t, err)
	require.Nil(t, preProcessor)

	args = getDefaultSmartContractResultPreProcessorCreatorArgs()
	preProcessor, err = fact.CreateSmartContractResultPreProcessor(args)
	require.Nil(t, err)
	require.NotNil(t, preProcessor)
	require.Implements(t, new(process.PreProcessor), preProcessor)
}

func TestSmartContractResultPreProcessorFactory_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	fact := preprocess.NewSmartContractResultPreProcessorFactory()
	require.False(t, fact.IsInterfaceNil())
}

func getDefaultSmartContractResultPreProcessorCreatorArgs() preprocess.SmartContractResultPreProcessorCreatorArgs {
	requestTransaction := func(shardID uint32, txHashes [][]byte) {}
	args := preprocess.SmartContractResultPreProcessorCreatorArgs{
		ScrDataPool: &testscommon.ShardedDataStub{},
		Store:       &storageStubs.ChainStorerStub{},
		Hasher:      &hashingMocks.HasherMock{},
		Marshalizer: &mock.MarshalizerMock{},
		ScrProcessor: &testscommon.TxProcessorMock{
			ProcessSmartContractResultCalled: func(scr *smartContractResult.SmartContractResult) (vmcommon.ReturnCode, error) {
				return 0, nil
			},
		},
		ShardCoordinator:             mock.NewMultiShardsCoordinatorMock(3),
		Accounts:                     &stateMock.AccountsStub{},
		OnRequestSmartContractResult: requestTransaction,
		GasHandler:                   &testscommon.GasHandlerStub{},
		EconomicsFee:                 &economicsmocks.EconomicsHandlerStub{},
		PubkeyConverter:              testscommon.NewPubkeyConverterMock(32),
		BlockSizeComputation:         &testscommon.BlockSizeComputationStub{},
		BalanceComputation:           &testscommon.BalanceComputationStub{},
		EnableEpochsHandler:          &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		ProcessedMiniBlocksTracker:   &testscommon.ProcessedMiniBlocksTrackerStub{},
		TxExecutionOrderHandler:      &common.TxExecutionOrderHandlerStub{},
	}
	return args
}
