package runType

import (
	"fmt"
	"math/big"

	nodeFactory "github.com/TerraDharitri/drt-go-chain/cmd/node/factory"
	"github.com/TerraDharitri/drt-go-chain/common/disabled"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/broadcastFactory"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/sposFactory"
	sovereignBlock "github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool/sovereign"
	requesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/requestersContainer"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/resolverscontainer"
	storageRequestFactory "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/storageRequestersContainer/factory"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/requestHandlers"
	"github.com/TerraDharitri/drt-go-chain/epochStart/bootstrap"
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/errors"
	mainFactory "github.com/TerraDharitri/drt-go-chain/factory"
	factoryBlock "github.com/TerraDharitri/drt-go-chain/factory/block"
	"github.com/TerraDharitri/drt-go-chain/factory/epochStartTrigger"
	"github.com/TerraDharitri/drt-go-chain/factory/processing/api"
	"github.com/TerraDharitri/drt-go-chain/factory/processing/dataRetriever"
	factoryVm "github.com/TerraDharitri/drt-go-chain/factory/vm"
	"github.com/TerraDharitri/drt-go-chain/genesis"
	"github.com/TerraDharitri/drt-go-chain/genesis/checking"
	"github.com/TerraDharitri/drt-go-chain/genesis/parsing"
	processGenesis "github.com/TerraDharitri/drt-go-chain/genesis/process"
	"github.com/TerraDharitri/drt-go-chain/node/external/transactionAPI"
	trieIteratorsFactory "github.com/TerraDharitri/drt-go-chain/node/trieIterators/factory"
	outportFactory "github.com/TerraDharitri/drt-go-chain/outport/process/factory"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block"
	processBlock "github.com/TerraDharitri/drt-go-chain/process/block"
	"github.com/TerraDharitri/drt-go-chain/process/block/preprocess"
	"github.com/TerraDharitri/drt-go-chain/process/block/sovereign"
	"github.com/TerraDharitri/drt-go-chain/process/coordinator"
	"github.com/TerraDharitri/drt-go-chain/process/factory/interceptorscontainer"
	"github.com/TerraDharitri/drt-go-chain/process/factory/shard"
	"github.com/TerraDharitri/drt-go-chain/process/factory/shard/data"
	"github.com/TerraDharitri/drt-go-chain/process/headerCheck"
	"github.com/TerraDharitri/drt-go-chain/process/peer"
	"github.com/TerraDharitri/drt-go-chain/process/scToProtocol"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/processProxy"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/scrCommon"
	"github.com/TerraDharitri/drt-go-chain/process/sync"
	"github.com/TerraDharitri/drt-go-chain/process/sync/storageBootstrap"
	"github.com/TerraDharitri/drt-go-chain/process/track"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	nodesCoord "github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/factory"
	syncerFactory "github.com/TerraDharitri/drt-go-chain/state/syncer/factory"
	storageFactory "github.com/TerraDharitri/drt-go-chain/storage/factory"
	"github.com/TerraDharitri/drt-go-chain/storage/latestData"
	"github.com/TerraDharitri/drt-go-chain/storage/storageunit"
	updateFactory "github.com/TerraDharitri/drt-go-chain/update/factory/creator"
	"github.com/TerraDharitri/drt-go-chain/vm/systemSmartContracts"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
)

// ArgsRunTypeComponents struct holds the arguments for run type component
type ArgsRunTypeComponents struct {
	CoreComponents   process.CoreComponentsHolder
	CryptoComponents process.CryptoComponentsHolder
	Configs          config.Configs
	InitialAccounts  []genesis.InitialAccountHandler
}

type runTypeComponentsFactory struct {
	coreComponents   process.CoreComponentsHolder
	cryptoComponents process.CryptoComponentsHolder
	configs          config.Configs
	initialAccounts  []genesis.InitialAccountHandler
}

// runTypeComponents struct holds the components needed for a run type
type runTypeComponents struct {
	blockChainHookHandlerCreator            hooks.BlockChainHookHandlerCreator
	epochStartBootstrapperCreator           bootstrap.EpochStartBootstrapperCreator
	bootstrapperFromStorageCreator          storageBootstrap.BootstrapperFromStorageCreator
	bootstrapperCreator                     storageBootstrap.BootstrapperCreator
	blockProcessorCreator                   processBlock.BlockProcessorCreator
	forkDetectorCreator                     sync.ForkDetectorCreator
	blockTrackerCreator                     track.BlockTrackerCreator
	requestHandlerCreator                   requestHandlers.RequestHandlerCreator
	headerValidatorCreator                  processBlock.HeaderValidatorCreator
	scheduledTxsExecutionCreator            preprocess.ScheduledTxsExecutionCreator
	transactionCoordinatorCreator           coordinator.TransactionCoordinatorCreator
	validatorStatisticsProcessorCreator     peer.ValidatorStatisticsProcessorCreator
	additionalStorageServiceCreator         process.AdditionalStorageServiceCreator
	scProcessorCreator                      scrCommon.SCProcessorCreator
	scResultPreProcessorCreator             preprocess.SmartContractResultPreProcessorCreator
	consensusModel                          consensus.ConsensusModel
	vmContainerMetaFactory                  factoryVm.VmContainerCreator
	vmContainerShardFactory                 factoryVm.VmContainerCreator
	accountsParser                          genesis.AccountsParser
	accountsCreator                         state.AccountFactory
	vmContextCreator                        systemSmartContracts.VMContextCreatorHandler
	outGoingOperationsPoolHandler           sovereignBlock.OutGoingOperationsPool
	dataCodecHandler                        sovereign.DataCodecHandler
	topicsCheckerHandler                    sovereign.TopicsCheckerHandler
	shardCoordinatorCreator                 sharding.ShardCoordinatorFactory
	nodesCoordinatorWithRaterFactoryCreator nodesCoord.NodesCoordinatorWithRaterFactory
	requestersContainerFactoryCreator       requesterscontainer.RequesterContainerFactoryCreator
	interceptorsContainerFactoryCreator     interceptorscontainer.InterceptorsContainerFactoryCreator
	shardResolversContainerFactoryCreator   resolverscontainer.ShardResolversContainerFactoryCreator
	txPreProcessorCreator                   preprocess.TxPreProcessorCreator
	extraHeaderSigVerifierHolder            headerCheck.ExtraHeaderSigVerifierHolder
	genesisBlockCreatorFactory              processGenesis.GenesisBlockCreatorFactory
	genesisMetaBlockCheckerCreator          processGenesis.GenesisMetaBlockChecker
	nodesSetupCheckerFactory                checking.NodesSetupCheckerFactory
	epochStartTriggerFactory                mainFactory.EpochStartTriggerFactoryHandler
	latestDataProviderFactory               latestData.LatestDataProviderFactory
	scToProtocolFactory                     scToProtocol.StakingToPeerFactoryHandler
	validatorInfoCreatorFactory             mainFactory.ValidatorInfoCreatorFactory
	apiProcessorCompsCreatorHandler         api.ApiProcessorCompsCreatorHandler
	endOfEpochEconomicsFactoryHandler       mainFactory.EndOfEpochEconomicsFactoryHandler
	rewardsCreatorFactory                   mainFactory.RewardsCreatorFactory
	systemSCProcessorFactory                mainFactory.SystemSCProcessorFactory
	preProcessorsContainerFactoryCreator    data.PreProcessorsContainerFactoryCreator
	dataRetrieverContainersSetter           mainFactory.DataRetrieverContainersSetter
	shardMessengerFactory                   sposFactory.BroadCastShardMessengerFactoryHandler
	exportHandlerFactoryCreator             mainFactory.ExportHandlerFactoryCreator
	validatorAccountsSyncerFactoryHandler   syncerFactory.ValidatorAccountsSyncerFactoryHandler
	shardRequestersContainerCreatorHandler  storageRequestFactory.ShardRequestersContainerCreatorHandler
	apiRewardTxHandler                      transactionAPI.APIRewardTxHandler
	outportDataProviderFactory              mainFactory.OutportDataProviderFactoryHandler
	delegatedListFactoryHandler             trieIteratorsFactory.DelegatedListProcessorFactoryHandler
	directStakedListFactoryHandler          trieIteratorsFactory.DirectStakedListProcessorFactoryHandler
	totalStakedValueFactoryHandler          trieIteratorsFactory.TotalStakedValueProcessorFactoryHandler
	versionedHeaderFactory                  genesis.VersionedHeaderFactory
}

// NewRunTypeComponentsFactory will return a new instance of runTypeComponentsFactory
func NewRunTypeComponentsFactory(args ArgsRunTypeComponents) (*runTypeComponentsFactory, error) {
	if check.IfNil(args.CoreComponents) {
		return nil, errors.ErrNilCoreComponents
	}
	if check.IfNil(args.CryptoComponents) {
		return nil, errors.ErrNilCryptoComponents
	}
	if args.InitialAccounts == nil {
		return nil, errors.ErrNilInitialAccounts
	}

	return &runTypeComponentsFactory{
		coreComponents:   args.CoreComponents,
		cryptoComponents: args.CryptoComponents,
		configs:          args.Configs,
		initialAccounts:  args.InitialAccounts,
	}, nil
}

// Create creates the runType components
func (rcf *runTypeComponentsFactory) Create() (*runTypeComponents, error) {
	vmContextCreator := systemSmartContracts.NewVMContextCreator()
	vmContainerMetaCreator, err := factoryVm.NewVmContainerMetaFactory(vmContextCreator)
	if err != nil {
		return nil, fmt.Errorf("runTypeComponentsFactory - NewVmContainerMetaFactory failed: %w", err)
	}

	totalSupply, ok := big.NewInt(0).SetString(rcf.configs.EconomicsConfig.GlobalSettings.GenesisTotalSupply, 10)
	if !ok {
		return nil, fmt.Errorf("can not parse total suply from economics.toml, %s is not a valid value",
			rcf.configs.EconomicsConfig.GlobalSettings.GenesisTotalSupply)
	}

	accountsParserArgs := genesis.AccountsParserArgs{
		InitialAccounts: rcf.initialAccounts,
		EntireSupply:    totalSupply,
		MinterAddress:   rcf.configs.EconomicsConfig.GlobalSettings.GenesisMintingSenderAddress,
		PubkeyConverter: rcf.coreComponents.AddressPubKeyConverter(),
		KeyGenerator:    rcf.cryptoComponents.TxSignKeyGen(),
		Hasher:          rcf.coreComponents.Hasher(),
		Marshalizer:     rcf.coreComponents.InternalMarshalizer(),
	}
	accountsParser, err := parsing.NewAccountsParser(accountsParserArgs)
	if err != nil {
		return nil, fmt.Errorf("runTypeComponentsFactory - NewAccountsParser failed: %w", err)
	}

	accountsCreator, err := factory.NewAccountCreator(factory.ArgsAccountCreator{
		Hasher:              rcf.coreComponents.Hasher(),
		Marshaller:          rcf.coreComponents.InternalMarshalizer(),
		EnableEpochsHandler: rcf.coreComponents.EnableEpochsHandler(),
	})
	if err != nil {
		return nil, fmt.Errorf("runTypeComponentsFactory - NewAccountCreator failed: %w", err)
	}
	apiRewardTxHandler, err := transactionAPI.NewAPIRewardsHandler(rcf.coreComponents.AddressPubKeyConverter())
	if err != nil {
		return nil, fmt.Errorf("runTypeComponentsFactory - NewAPIRewardsHandler failed: %w", err)
	}

	headerVersionHandler, err := rcf.createHeaderVersionHandler()
	if err != nil {
		return nil, fmt.Errorf("runTypeComponentsFactory - createHeaderVersionHandler failed: %w", err)
	}

	versionedHeaderFactory, err := factoryBlock.NewShardHeaderFactory(headerVersionHandler)
	if err != nil {
		return nil, fmt.Errorf("runTypeComponentsFactory - NewShardHeaderFactory failed: %w", err)
	}

	return &runTypeComponents{
		blockChainHookHandlerCreator:            hooks.NewBlockChainHookFactory(),
		epochStartBootstrapperCreator:           bootstrap.NewEpochStartBootstrapperFactory(),
		bootstrapperFromStorageCreator:          storageBootstrap.NewShardStorageBootstrapperFactory(),
		bootstrapperCreator:                     storageBootstrap.NewShardBootstrapFactory(),
		blockProcessorCreator:                   block.NewShardBlockProcessorFactory(),
		forkDetectorCreator:                     sync.NewShardForkDetectorFactory(),
		blockTrackerCreator:                     track.NewShardBlockTrackerFactory(),
		requestHandlerCreator:                   requestHandlers.NewResolverRequestHandlerFactory(),
		headerValidatorCreator:                  block.NewShardHeaderValidatorFactory(),
		scheduledTxsExecutionCreator:            preprocess.NewShardScheduledTxsExecutionFactory(),
		transactionCoordinatorCreator:           coordinator.NewShardTransactionCoordinatorFactory(),
		validatorStatisticsProcessorCreator:     peer.NewValidatorStatisticsProcessorFactory(),
		additionalStorageServiceCreator:         storageFactory.NewShardAdditionalStorageServiceFactory(),
		scProcessorCreator:                      processProxy.NewSCProcessProxyFactory(),
		scResultPreProcessorCreator:             preprocess.NewSmartContractResultPreProcessorFactory(),
		consensusModel:                          consensus.ConsensusModelV1,
		vmContainerMetaFactory:                  vmContainerMetaCreator,
		vmContainerShardFactory:                 factoryVm.NewVmContainerShardFactory(),
		accountsParser:                          accountsParser,
		accountsCreator:                         accountsCreator,
		vmContextCreator:                        vmContextCreator,
		outGoingOperationsPoolHandler:           disabled.NewDisabledOutGoingOperationPool(),
		dataCodecHandler:                        disabled.NewDisabledDataCodec(),
		topicsCheckerHandler:                    disabled.NewDisabledTopicsChecker(),
		shardCoordinatorCreator:                 sharding.NewMultiShardCoordinatorFactory(),
		nodesCoordinatorWithRaterFactoryCreator: nodesCoord.NewIndexHashedNodesCoordinatorWithRaterFactory(),
		requestersContainerFactoryCreator:       requesterscontainer.NewShardRequestersContainerFactoryCreator(),
		interceptorsContainerFactoryCreator:     interceptorscontainer.NewShardInterceptorsContainerFactoryCreator(),
		shardResolversContainerFactoryCreator:   resolverscontainer.NewShardResolversContainerFactoryCreator(),
		txPreProcessorCreator:                   preprocess.NewTxPreProcessorCreator(),
		extraHeaderSigVerifierHolder:            headerCheck.NewExtraHeaderSigVerifierHolder(),
		genesisBlockCreatorFactory:              processGenesis.NewGenesisBlockCreatorFactory(),
		genesisMetaBlockCheckerCreator:          processGenesis.NewGenesisMetaBlockChecker(),
		nodesSetupCheckerFactory:                checking.NewNodesSetupCheckerFactory(),
		epochStartTriggerFactory:                epochStartTrigger.NewEpochStartTriggerFactory(),
		latestDataProviderFactory:               latestData.NewLatestDataProviderFactory(),
		scToProtocolFactory:                     scToProtocol.NewStakingToPeerFactory(),
		validatorInfoCreatorFactory:             metachain.NewValidatorInfoCreatorFactory(),
		apiProcessorCompsCreatorHandler:         api.NewAPIProcessorCompsCreator(),
		endOfEpochEconomicsFactoryHandler:       metachain.NewEconomicsFactory(),
		rewardsCreatorFactory:                   metachain.NewRewardsCreatorFactory(),
		systemSCProcessorFactory:                metachain.NewSysSCFactory(),
		preProcessorsContainerFactoryCreator:    shard.NewPreProcessorContainerFactoryCreator(),
		dataRetrieverContainersSetter:           dataRetriever.NewDataRetrieverContainerSetter(),
		shardMessengerFactory:                   broadcastFactory.NewShardChainMessengerFactory(),
		exportHandlerFactoryCreator:             updateFactory.NewExportHandlerFactoryCreator(),
		validatorAccountsSyncerFactoryHandler:   syncerFactory.NewValidatorAccountsSyncerFactory(),
		shardRequestersContainerCreatorHandler:  storageRequestFactory.NewShardRequestersContainerCreator(),
		apiRewardTxHandler:                      apiRewardTxHandler,
		outportDataProviderFactory:              outportFactory.NewOutportDataProviderFactory(),
		delegatedListFactoryHandler:             trieIteratorsFactory.NewDelegatedListProcessorFactory(),
		directStakedListFactoryHandler:          trieIteratorsFactory.NewDirectStakedListProcessorFactory(),
		totalStakedValueFactoryHandler:          trieIteratorsFactory.NewTotalStakedListProcessorFactory(),
		versionedHeaderFactory:                  versionedHeaderFactory,
	}, nil
}

func (rcf *runTypeComponentsFactory) createHeaderVersionHandler() (nodeFactory.HeaderVersionHandler, error) {
	cacheConfig := storageFactory.GetCacherFromConfig(rcf.configs.GeneralConfig.Versions.Cache)
	cache, err := storageunit.NewCache(cacheConfig)
	if err != nil {
		return nil, err
	}

	return factoryBlock.NewHeaderVersionHandler(
		rcf.configs.GeneralConfig.Versions.VersionsByEpochs,
		rcf.configs.GeneralConfig.Versions.DefaultVersion,
		cache,
	)
}

// IsInterfaceNil returns true if there is no value under the interface
func (rc *runTypeComponentsFactory) IsInterfaceNil() bool {
	return rc == nil
}

// Close does nothing
func (rc *runTypeComponents) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (rc *runTypeComponents) IsInterfaceNil() bool {
	return rc == nil
}
