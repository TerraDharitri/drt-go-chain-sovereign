package sovereign

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain/process/factory/shard/data"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	mockCommon "github.com/TerraDharitri/drt-go-chain/testscommon/common"
	dataRetrieverMock "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/processMocks"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	storageStubs "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/stretchr/testify/require"
)

func createMockPreProcessorsContainerFactoryArguments() data.ArgPreProcessorsContainerFactory {
	return data.ArgPreProcessorsContainerFactory{
		ShardCoordinator:             mock.NewMultiShardsCoordinatorMock(3),
		Store:                        &storageStubs.ChainStorerStub{},
		Marshaller:                   &mock.MarshalizerMock{},
		Hasher:                       &hashingMocks.HasherMock{},
		DataPool:                     dataRetrieverMock.NewPoolsHolderMock(),
		PubkeyConverter:              testscommon.NewPubkeyConverterMock(32),
		Accounts:                     &stateMock.AccountsStub{},
		RequestHandler:               &testscommon.RequestHandlerStub{},
		TxProcessor:                  &testscommon.TxProcessorMock{},
		ScProcessor:                  &testscommon.SCProcessorMock{},
		ScResultProcessor:            &testscommon.SmartContractResultsProcessorMock{},
		RewardsTxProcessor:           &testscommon.RewardTxProcessorMock{},
		EconomicsFee:                 &economicsmocks.EconomicsHandlerStub{},
		GasHandler:                   &testscommon.GasHandlerStub{},
		BlockTracker:                 &mock.BlockTrackerMock{},
		BlockSizeComputation:         &testscommon.BlockSizeComputationStub{},
		BalanceComputation:           &testscommon.BalanceComputationStub{},
		EnableEpochsHandler:          &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		TxTypeHandler:                &testscommon.TxTypeHandlerMock{},
		ScheduledTxsExecutionHandler: &testscommon.ScheduledTxsExecutionStub{},
		ProcessedMiniBlocksTracker:   &testscommon.ProcessedMiniBlocksTrackerStub{},
		TxExecutionOrderHandler:      &mockCommon.TxExecutionOrderHandlerStub{},
		RunTypeComponents:            processMocks.NewRunTypeComponentsStub(),
	}
}

func TestSovereignPreProcessorContainerFactoryCreator_CreatePreProcessorContainerFactory(t *testing.T) {
	t.Parallel()

	f := NewSovereignPreProcessorContainerFactoryCreator()
	require.False(t, f.IsInterfaceNil())

	args := createMockPreProcessorsContainerFactoryArguments()
	containerFactory, err := f.CreatePreProcessorContainerFactory(args)
	require.Nil(t, err)
	require.Equal(t, "*sovereign.sovereignPreProcessorsContainerFactory", fmt.Sprintf("%T", containerFactory))
}
