package mainFactoryMocks

import (
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/sposFactory"
	sovereignBlock "github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool/sovereign"
	requesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/requestersContainer"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/resolverscontainer"
	storageRequestFactory "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/storageRequestersContainer/factory"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/requestHandlers"
	"github.com/TerraDharitri/drt-go-chain/epochStart/bootstrap"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/factory/processing/api"
	factoryVm "github.com/TerraDharitri/drt-go-chain/factory/vm"
	"github.com/TerraDharitri/drt-go-chain/genesis"
	"github.com/TerraDharitri/drt-go-chain/genesis/checking"
	processGenesis "github.com/TerraDharitri/drt-go-chain/genesis/process"
	"github.com/TerraDharitri/drt-go-chain/node/external/transactionAPI"
	trieIteratorsFactory "github.com/TerraDharitri/drt-go-chain/node/trieIterators/factory"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block"
	"github.com/TerraDharitri/drt-go-chain/process/block/preprocess"
	"github.com/TerraDharitri/drt-go-chain/process/block/sovereign"
	"github.com/TerraDharitri/drt-go-chain/process/coordinator"
	"github.com/TerraDharitri/drt-go-chain/process/factory/interceptorscontainer"
	"github.com/TerraDharitri/drt-go-chain/process/factory/shard/data"
	"github.com/TerraDharitri/drt-go-chain/process/headerCheck"
	"github.com/TerraDharitri/drt-go-chain/process/peer"
	"github.com/TerraDharitri/drt-go-chain/process/scToProtocol"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/scrCommon"
	"github.com/TerraDharitri/drt-go-chain/process/sync"
	"github.com/TerraDharitri/drt-go-chain/process/sync/storageBootstrap"
	"github.com/TerraDharitri/drt-go-chain/process/track"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	nodesCoord "github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	syncerFactory "github.com/TerraDharitri/drt-go-chain/state/syncer/factory"
	"github.com/TerraDharitri/drt-go-chain/storage/latestData"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	apiTests "github.com/TerraDharitri/drt-go-chain/testscommon/api"
	testFactory "github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genesisMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/headerSigVerifier"
	sovereignMocks "github.com/TerraDharitri/drt-go-chain/testscommon/sovereign"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon/vmContext"
	"github.com/TerraDharitri/drt-go-chain/vm/systemSmartContracts"
)

// RunTypeComponentsStub -
type RunTypeComponentsStub struct {
	BlockChainHookHandlerFactory                hooks.BlockChainHookHandlerCreator
	BlockProcessorFactory                       block.BlockProcessorCreator
	BlockTrackerFactory                         track.BlockTrackerCreator
	BootstrapperFromStorageFactory              storageBootstrap.BootstrapperFromStorageCreator
	BootstrapperFactory                         storageBootstrap.BootstrapperCreator
	EpochStartBootstrapperFactory               bootstrap.EpochStartBootstrapperCreator
	ForkDetectorFactory                         sync.ForkDetectorCreator
	HeaderValidatorFactory                      block.HeaderValidatorCreator
	RequestHandlerFactory                       requestHandlers.RequestHandlerCreator
	ScheduledTxsExecutionFactory                preprocess.ScheduledTxsExecutionCreator
	TransactionCoordinatorFactory               coordinator.TransactionCoordinatorCreator
	ValidatorStatisticsProcessorFactory         peer.ValidatorStatisticsProcessorCreator
	AdditionalStorageServiceFactory             process.AdditionalStorageServiceCreator
	SCResultsPreProcessorFactory                preprocess.SmartContractResultPreProcessorCreator
	SCProcessorFactory                          scrCommon.SCProcessorCreator
	ConsensusModelType                          consensus.ConsensusModel
	VmContainerMetaFactory                      factoryVm.VmContainerCreator
	VmContainerShardFactory                     factoryVm.VmContainerCreator
	AccountParser                               genesis.AccountsParser
	AccountCreator                              state.AccountFactory
	VMContextCreatorHandler                     systemSmartContracts.VMContextCreatorHandler
	OutGoingOperationsPool                      sovereignBlock.OutGoingOperationsPool
	DataCodec                                   sovereign.DataCodecHandler
	TopicsChecker                               sovereign.TopicsCheckerHandler
	ShardCoordinatorFactory                     sharding.ShardCoordinatorFactory
	NodesCoordinatorWithRaterFactory            nodesCoord.NodesCoordinatorWithRaterFactory
	RequestersContainerFactory                  requesterscontainer.RequesterContainerFactoryCreator
	InterceptorsContainerFactory                interceptorscontainer.InterceptorsContainerFactoryCreator
	ShardResolversContainerFactory              resolverscontainer.ShardResolversContainerFactoryCreator
	TxPreProcessorFactory                       preprocess.TxPreProcessorCreator
	ExtraHeaderSigVerifier                      headerCheck.ExtraHeaderSigVerifierHolder
	GenesisBlockFactory                         processGenesis.GenesisBlockCreatorFactory
	GenesisMetaBlockChecker                     processGenesis.GenesisMetaBlockChecker
	NodesSetupCheckerFactoryField               checking.NodesSetupCheckerFactory
	EpochStartTriggerFactoryField               factory.EpochStartTriggerFactoryHandler
	LatestDataProviderFactoryField              latestData.LatestDataProviderFactory
	StakingToPeerFactoryField                   scToProtocol.StakingToPeerFactoryHandler
	ValidatorInfoCreatorFactoryField            factory.ValidatorInfoCreatorFactory
	APIProcessorCompsCreatorHandlerField        api.ApiProcessorCompsCreatorHandler
	EndOfEpochEconomicsFactoryHandlerField      factory.EndOfEpochEconomicsFactoryHandler
	RewardsCreatorFactoryField                  factory.RewardsCreatorFactory
	SystemSCProcessorFactoryField               factory.SystemSCProcessorFactory
	PreProcessorsContainerFactoryCreatorField   data.PreProcessorsContainerFactoryCreator
	DataRetrieverContainersSetterField          factory.DataRetrieverContainersSetter
	ShardMessengerFactoryField                  sposFactory.BroadCastShardMessengerFactoryHandler
	ExportHandlerFactoryCreatorField            factory.ExportHandlerFactoryCreator
	ValidatorAccountsSyncerFactoryHandlerField  syncerFactory.ValidatorAccountsSyncerFactoryHandler
	ShardRequestersContainerCreatorHandlerField storageRequestFactory.ShardRequestersContainerCreatorHandler
	APIRewardsTxHandlerField                    transactionAPI.APIRewardTxHandler
	OutportDataProviderFactoryField             factory.OutportDataProviderFactoryHandler
	DelegatedListFactoryField                   trieIteratorsFactory.DelegatedListProcessorFactoryHandler
	DirectStakedListFactoryField                trieIteratorsFactory.DirectStakedListProcessorFactoryHandler
	TotalStakedValueFactoryField                trieIteratorsFactory.TotalStakedValueProcessorFactoryHandler
	VersionedHeaderFactoryField                 genesis.VersionedHeaderFactory
}

// NewRunTypeComponentsStub -
func NewRunTypeComponentsStub() *RunTypeComponentsStub {
	return &RunTypeComponentsStub{
		BlockChainHookHandlerFactory:                &testFactory.BlockChainHookHandlerFactoryMock{},
		BlockProcessorFactory:                       &testFactory.BlockProcessorFactoryMock{},
		BlockTrackerFactory:                         &testFactory.BlockTrackerFactoryMock{},
		BootstrapperFromStorageFactory:              &testFactory.BootstrapperFromStorageFactoryMock{},
		BootstrapperFactory:                         &testFactory.BootstrapperFactoryMock{},
		EpochStartBootstrapperFactory:               &testFactory.EpochStartBootstrapperFactoryMock{},
		ForkDetectorFactory:                         &testFactory.ForkDetectorFactoryMock{},
		HeaderValidatorFactory:                      &testFactory.HeaderValidatorFactoryMock{},
		RequestHandlerFactory:                       &testFactory.RequestHandlerFactoryMock{},
		ScheduledTxsExecutionFactory:                &testFactory.ScheduledTxsExecutionFactoryMock{},
		TransactionCoordinatorFactory:               &testFactory.TransactionCoordinatorFactoryMock{},
		ValidatorStatisticsProcessorFactory:         &testFactory.ValidatorStatisticsProcessorFactoryMock{},
		AdditionalStorageServiceFactory:             &testFactory.AdditionalStorageServiceFactoryMock{},
		SCResultsPreProcessorFactory:                &testFactory.SmartContractResultPreProcessorFactoryMock{},
		SCProcessorFactory:                          &testFactory.SCProcessorFactoryMock{},
		ConsensusModelType:                          consensus.ConsensusModelV1,
		VmContainerMetaFactory:                      &testFactory.VMContainerFactoryMock{},
		VmContainerShardFactory:                     &testFactory.VMContainerFactoryMock{},
		AccountParser:                               &genesisMocks.AccountsParserStub{},
		AccountCreator:                              &stateMock.AccountsFactoryStub{},
		VMContextCreatorHandler:                     &vmContext.VMContextCreatorStub{},
		OutGoingOperationsPool:                      &sovereignMocks.OutGoingOperationsPoolMock{},
		DataCodec:                                   &sovereignMocks.DataCodecMock{},
		TopicsChecker:                               &sovereignMocks.TopicsCheckerMock{},
		ShardCoordinatorFactory:                     &testscommon.MultiShardCoordinatorFactoryMock{},
		NodesCoordinatorWithRaterFactory:            &testscommon.NodesCoordinatorFactoryMock{},
		RequestersContainerFactory:                  &testFactory.RequestersContainerFactoryMock{},
		InterceptorsContainerFactory:                &testFactory.InterceptorsContainerFactoryMock{},
		ShardResolversContainerFactory:              &testFactory.ResolversContainerFactoryMock{},
		TxPreProcessorFactory:                       &testFactory.TxPreProcessorFactoryMock{},
		ExtraHeaderSigVerifier:                      &headerSigVerifier.ExtraHeaderSigVerifierHolderMock{},
		GenesisBlockFactory:                         &testFactory.GenesisBlockCreatorFactoryMock{},
		GenesisMetaBlockChecker:                     &testFactory.GenesisMetaBlockCheckerMock{},
		NodesSetupCheckerFactoryField:               checking.NewNodesSetupCheckerFactory(),
		EpochStartTriggerFactoryField:               &testFactory.EpochStartTriggerFactoryMock{},
		LatestDataProviderFactoryField:              &testFactory.LatestDataProviderFactoryMock{},
		StakingToPeerFactoryField:                   &testFactory.StakingToPeerFactoryMock{},
		ValidatorInfoCreatorFactoryField:            &testFactory.ValidatorInfoCreatorFactoryMock{},
		APIProcessorCompsCreatorHandlerField:        &testFactory.APIProcessorCompsCreatorMock{},
		EndOfEpochEconomicsFactoryHandlerField:      &testFactory.EconomicsFactoryMock{},
		SystemSCProcessorFactoryField:               &testFactory.SysSCFactoryMock{},
		PreProcessorsContainerFactoryCreatorField:   &testFactory.PreProcessorContainerFactoryCreatorMock{},
		DataRetrieverContainersSetterField:          &testFactory.DataRetrieverContainersSetterMock{},
		ShardMessengerFactoryField:                  &testFactory.ShardChainMessengerFactoryMock{},
		ExportHandlerFactoryCreatorField:            &testFactory.ExportHandlerFactoryCreatorMock{},
		ValidatorAccountsSyncerFactoryHandlerField:  &testFactory.ValidatorAccountsSyncerFactoryMock{},
		ShardRequestersContainerCreatorHandlerField: &testFactory.ShardRequestersContainerCreatorMock{},
		APIRewardsTxHandlerField:                    &apiTests.APIRewardsHandlerStub{},
		OutportDataProviderFactoryField:             &testFactory.OutportDataProviderFactoryMock{},
		VersionedHeaderFactoryField:                 &testscommon.VersionedHeaderFactoryStub{},
	}
}

// Create -
func (r *RunTypeComponentsStub) Create() error {
	return nil
}

// Close -
func (r *RunTypeComponentsStub) Close() error {
	return nil
}

// CheckSubcomponents -
func (r *RunTypeComponentsStub) CheckSubcomponents() error {
	return nil
}

// String -
func (r *RunTypeComponentsStub) String() string {
	return ""
}

// BlockChainHookHandlerCreator -
func (r *RunTypeComponentsStub) BlockChainHookHandlerCreator() hooks.BlockChainHookHandlerCreator {
	return r.BlockChainHookHandlerFactory
}

// BlockProcessorCreator -
func (r *RunTypeComponentsStub) BlockProcessorCreator() block.BlockProcessorCreator {
	return r.BlockProcessorFactory
}

// BlockTrackerCreator -
func (r *RunTypeComponentsStub) BlockTrackerCreator() track.BlockTrackerCreator {
	return r.BlockTrackerFactory
}

// BootstrapperFromStorageCreator -
func (r *RunTypeComponentsStub) BootstrapperFromStorageCreator() storageBootstrap.BootstrapperFromStorageCreator {
	return r.BootstrapperFromStorageFactory
}

// BootstrapperCreator -
func (r *RunTypeComponentsStub) BootstrapperCreator() storageBootstrap.BootstrapperCreator {
	return r.BootstrapperFactory
}

// EpochStartBootstrapperCreator -
func (r *RunTypeComponentsStub) EpochStartBootstrapperCreator() bootstrap.EpochStartBootstrapperCreator {
	return r.EpochStartBootstrapperFactory
}

// ForkDetectorCreator -
func (r *RunTypeComponentsStub) ForkDetectorCreator() sync.ForkDetectorCreator {
	return r.ForkDetectorFactory
}

// HeaderValidatorCreator -
func (r *RunTypeComponentsStub) HeaderValidatorCreator() block.HeaderValidatorCreator {
	return r.HeaderValidatorFactory
}

// RequestHandlerCreator -
func (r *RunTypeComponentsStub) RequestHandlerCreator() requestHandlers.RequestHandlerCreator {
	return r.RequestHandlerFactory
}

// ScheduledTxsExecutionCreator -
func (r *RunTypeComponentsStub) ScheduledTxsExecutionCreator() preprocess.ScheduledTxsExecutionCreator {
	return r.ScheduledTxsExecutionFactory
}

// TransactionCoordinatorCreator -
func (r *RunTypeComponentsStub) TransactionCoordinatorCreator() coordinator.TransactionCoordinatorCreator {
	return r.TransactionCoordinatorFactory
}

// ValidatorStatisticsProcessorCreator -
func (r *RunTypeComponentsStub) ValidatorStatisticsProcessorCreator() peer.ValidatorStatisticsProcessorCreator {
	return r.ValidatorStatisticsProcessorFactory
}

// AdditionalStorageServiceCreator -
func (r *RunTypeComponentsStub) AdditionalStorageServiceCreator() process.AdditionalStorageServiceCreator {
	return r.AdditionalStorageServiceFactory
}

// SCProcessorCreator -
func (r *RunTypeComponentsStub) SCProcessorCreator() scrCommon.SCProcessorCreator {
	return r.SCProcessorFactory
}

// SCResultsPreProcessorCreator -
func (r *RunTypeComponentsStub) SCResultsPreProcessorCreator() preprocess.SmartContractResultPreProcessorCreator {
	return r.SCResultsPreProcessorFactory
}

// ConsensusModel -
func (r *RunTypeComponentsStub) ConsensusModel() consensus.ConsensusModel {
	return r.ConsensusModelType
}

// VmContainerMetaFactoryCreator -
func (r *RunTypeComponentsStub) VmContainerMetaFactoryCreator() factoryVm.VmContainerCreator {
	return r.VmContainerMetaFactory
}

// VmContainerShardFactoryCreator -
func (r *RunTypeComponentsStub) VmContainerShardFactoryCreator() factoryVm.VmContainerCreator {
	return r.VmContainerShardFactory
}

// AccountsParser -
func (r *RunTypeComponentsStub) AccountsParser() genesis.AccountsParser {
	return r.AccountParser
}

// AccountsCreator -
func (r *RunTypeComponentsStub) AccountsCreator() state.AccountFactory {
	return r.AccountCreator
}

// VMContextCreator -
func (r *RunTypeComponentsStub) VMContextCreator() systemSmartContracts.VMContextCreatorHandler {
	return r.VMContextCreatorHandler
}

// OutGoingOperationsPoolHandler -
func (r *RunTypeComponentsStub) OutGoingOperationsPoolHandler() sovereignBlock.OutGoingOperationsPool {
	return r.OutGoingOperationsPool
}

// DataCodecHandler -
func (r *RunTypeComponentsStub) DataCodecHandler() sovereign.DataCodecHandler {
	return r.DataCodec
}

// TopicsCheckerHandler -
func (r *RunTypeComponentsStub) TopicsCheckerHandler() sovereign.TopicsCheckerHandler {
	return r.TopicsChecker
}

// ShardCoordinatorCreator -
func (r *RunTypeComponentsStub) ShardCoordinatorCreator() sharding.ShardCoordinatorFactory {
	return r.ShardCoordinatorFactory
}

// NodesCoordinatorWithRaterCreator -
func (r *RunTypeComponentsStub) NodesCoordinatorWithRaterCreator() nodesCoord.NodesCoordinatorWithRaterFactory {
	return r.NodesCoordinatorWithRaterFactory
}

// RequestersContainerFactoryCreator -
func (r *RunTypeComponentsStub) RequestersContainerFactoryCreator() requesterscontainer.RequesterContainerFactoryCreator {
	return r.RequestersContainerFactory
}

// InterceptorsContainerFactoryCreator -
func (r *RunTypeComponentsStub) InterceptorsContainerFactoryCreator() interceptorscontainer.InterceptorsContainerFactoryCreator {
	return r.InterceptorsContainerFactory
}

// ShardResolversContainerFactoryCreator -
func (r *RunTypeComponentsStub) ShardResolversContainerFactoryCreator() resolverscontainer.ShardResolversContainerFactoryCreator {
	return r.ShardResolversContainerFactory
}

// TxPreProcessorCreator -
func (r *RunTypeComponentsStub) TxPreProcessorCreator() preprocess.TxPreProcessorCreator {
	return r.TxPreProcessorFactory
}

// ExtraHeaderSigVerifierHolder -
func (r *RunTypeComponentsStub) ExtraHeaderSigVerifierHolder() headerCheck.ExtraHeaderSigVerifierHolder {
	return r.ExtraHeaderSigVerifier
}

// GenesisBlockCreatorFactory -
func (r *RunTypeComponentsStub) GenesisBlockCreatorFactory() processGenesis.GenesisBlockCreatorFactory {
	return r.GenesisBlockFactory
}

// GenesisMetaBlockCheckerCreator -
func (r *RunTypeComponentsStub) GenesisMetaBlockCheckerCreator() processGenesis.GenesisMetaBlockChecker {
	return r.GenesisMetaBlockChecker
}

// NodesSetupCheckerFactory -
func (r *RunTypeComponentsStub) NodesSetupCheckerFactory() checking.NodesSetupCheckerFactory {
	return r.NodesSetupCheckerFactoryField
}

// EpochStartTriggerFactory -
func (r *RunTypeComponentsStub) EpochStartTriggerFactory() factory.EpochStartTriggerFactoryHandler {
	return r.EpochStartTriggerFactoryField
}

// LatestDataProviderFactory  -
func (r *RunTypeComponentsStub) LatestDataProviderFactory() latestData.LatestDataProviderFactory {
	return r.LatestDataProviderFactoryField
}

// StakingToPeerFactory -
func (r *RunTypeComponentsStub) StakingToPeerFactory() scToProtocol.StakingToPeerFactoryHandler {
	return r.StakingToPeerFactoryField
}

// ValidatorInfoCreatorFactory -
func (r *RunTypeComponentsStub) ValidatorInfoCreatorFactory() factory.ValidatorInfoCreatorFactory {
	return r.ValidatorInfoCreatorFactoryField
}

// ApiProcessorCompsCreatorHandler -
func (r *RunTypeComponentsStub) ApiProcessorCompsCreatorHandler() api.ApiProcessorCompsCreatorHandler {
	return r.APIProcessorCompsCreatorHandlerField
}

// EndOfEpochEconomicsFactoryHandler -
func (r *RunTypeComponentsStub) EndOfEpochEconomicsFactoryHandler() factory.EndOfEpochEconomicsFactoryHandler {
	return r.EndOfEpochEconomicsFactoryHandlerField
}

// RewardsCreatorFactory -
func (r *RunTypeComponentsStub) RewardsCreatorFactory() factory.RewardsCreatorFactory {
	return r.RewardsCreatorFactoryField
}

// SystemSCProcessorFactory -
func (r *RunTypeComponentsStub) SystemSCProcessorFactory() factory.SystemSCProcessorFactory {
	return r.SystemSCProcessorFactoryField
}

// PreProcessorsContainerFactoryCreator -
func (r *RunTypeComponentsStub) PreProcessorsContainerFactoryCreator() data.PreProcessorsContainerFactoryCreator {
	return r.PreProcessorsContainerFactoryCreatorField
}

// DataRetrieverContainersSetter -
func (r *RunTypeComponentsStub) DataRetrieverContainersSetter() factory.DataRetrieverContainersSetter {
	return r.DataRetrieverContainersSetterField
}

// BroadCastShardMessengerFactoryHandler -
func (r *RunTypeComponentsStub) BroadCastShardMessengerFactoryHandler() sposFactory.BroadCastShardMessengerFactoryHandler {
	return r.ShardMessengerFactoryField
}

// ExportHandlerFactoryCreator -
func (r *RunTypeComponentsStub) ExportHandlerFactoryCreator() factory.ExportHandlerFactoryCreator {
	return r.ExportHandlerFactoryCreatorField
}

// ValidatorAccountsSyncerFactoryHandler -
func (r *RunTypeComponentsStub) ValidatorAccountsSyncerFactoryHandler() syncerFactory.ValidatorAccountsSyncerFactoryHandler {
	return r.ValidatorAccountsSyncerFactoryHandlerField
}

// ShardRequestersContainerCreatorHandler -
func (r *RunTypeComponentsStub) ShardRequestersContainerCreatorHandler() storageRequestFactory.ShardRequestersContainerCreatorHandler {
	return r.ShardRequestersContainerCreatorHandlerField
}

// APIRewardsTxHandler -
func (r *RunTypeComponentsStub) APIRewardsTxHandler() transactionAPI.APIRewardTxHandler {
	return r.APIRewardsTxHandlerField
}

// OutportDataProviderFactory -
func (r *RunTypeComponentsStub) OutportDataProviderFactory() factory.OutportDataProviderFactoryHandler {
	return r.OutportDataProviderFactoryField
}

// DelegatedListFactoryHandler -
func (r *RunTypeComponentsStub) DelegatedListFactoryHandler() trieIteratorsFactory.DelegatedListProcessorFactoryHandler {
	return r.DelegatedListFactoryField
}

// DirectStakedListFactoryHandler -
func (r *RunTypeComponentsStub) DirectStakedListFactoryHandler() trieIteratorsFactory.DirectStakedListProcessorFactoryHandler {
	return r.DirectStakedListFactoryField
}

// TotalStakedValueFactoryHandler -
func (r *RunTypeComponentsStub) TotalStakedValueFactoryHandler() trieIteratorsFactory.TotalStakedValueProcessorFactoryHandler {
	return r.TotalStakedValueFactoryField
}

// VersionedHeaderFactory  -
func (r *RunTypeComponentsStub) VersionedHeaderFactory() genesis.VersionedHeaderFactory {
	return r.VersionedHeaderFactoryField
}

// IsInterfaceNil -
func (r *RunTypeComponentsStub) IsInterfaceNil() bool {
	return r == nil
}
