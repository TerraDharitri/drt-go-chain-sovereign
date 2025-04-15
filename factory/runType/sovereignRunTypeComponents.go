package runType

import (
	"fmt"
	"math/big"
	"time"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/broadcastFactory"
	sovereignFactory "github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool/sovereign"
	requesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/requestersContainer"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/resolverscontainer"
	storageRequestFactory "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/storageRequestersContainer/factory"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/requestHandlers"
	"github.com/TerraDharitri/drt-go-chain/epochStart/bootstrap"
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/errors"
	factoryBlock "github.com/TerraDharitri/drt-go-chain/factory/block"
	"github.com/TerraDharitri/drt-go-chain/factory/epochStartTrigger"
	"github.com/TerraDharitri/drt-go-chain/factory/processing/api"
	"github.com/TerraDharitri/drt-go-chain/factory/processing/dataRetriever"
	factoryVm "github.com/TerraDharitri/drt-go-chain/factory/vm"
	"github.com/TerraDharitri/drt-go-chain/genesis"
	"github.com/TerraDharitri/drt-go-chain/genesis/parsing"
	processComp "github.com/TerraDharitri/drt-go-chain/genesis/process"
	"github.com/TerraDharitri/drt-go-chain/node/external/transactionAPI"
	trieIteratorsFactory "github.com/TerraDharitri/drt-go-chain/node/trieIterators/factory"
	outportFactory "github.com/TerraDharitri/drt-go-chain/outport/process/factory"
	"github.com/TerraDharitri/drt-go-chain/process/block"
	"github.com/TerraDharitri/drt-go-chain/process/block/preprocess"
	"github.com/TerraDharitri/drt-go-chain/process/block/sovereign"
	"github.com/TerraDharitri/drt-go-chain/process/coordinator"
	"github.com/TerraDharitri/drt-go-chain/process/factory/interceptorscontainer"
	procSovereign "github.com/TerraDharitri/drt-go-chain/process/factory/sovereign"
	"github.com/TerraDharitri/drt-go-chain/process/headerCheck"
	"github.com/TerraDharitri/drt-go-chain/process/peer"
	"github.com/TerraDharitri/drt-go-chain/process/scToProtocol"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/processorV2"
	"github.com/TerraDharitri/drt-go-chain/process/sync"
	"github.com/TerraDharitri/drt-go-chain/process/sync/storageBootstrap"
	"github.com/TerraDharitri/drt-go-chain/process/track"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	nodesCoord "github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state/factory"
	syncerFactory "github.com/TerraDharitri/drt-go-chain/state/syncer/factory"
	storageFactory "github.com/TerraDharitri/drt-go-chain/storage/factory"
	"github.com/TerraDharitri/drt-go-chain/storage/latestData"
	updateFactory "github.com/TerraDharitri/drt-go-chain/update/factory/creator"
	"github.com/TerraDharitri/drt-go-chain/vm/systemSmartContracts"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
)

type ArgsSovereignRunTypeComponents struct {
	RunTypeComponentsFactory *runTypeComponentsFactory
	Config                   config.SovereignConfig
	DataCodec                sovereign.DataCodecHandler
	TopicsChecker            sovereign.TopicsCheckerHandler
}

type sovereignRunTypeComponentsFactory struct {
	*runTypeComponentsFactory
	sovConfig     config.SovereignConfig
	dataCodec     sovereign.DataCodecHandler
	topicsChecker sovereign.TopicsCheckerHandler
}

// NewSovereignRunTypeComponentsFactory will return a new instance of runTypeComponentsFactory
func NewSovereignRunTypeComponentsFactory(args ArgsSovereignRunTypeComponents) (*sovereignRunTypeComponentsFactory, error) {
	if check.IfNil(args.RunTypeComponentsFactory) {
		return nil, errors.ErrNilRunTypeComponentsFactory
	}
	if check.IfNil(args.DataCodec) {
		return nil, errors.ErrNilDataCodec
	}
	if check.IfNil(args.TopicsChecker) {
		return nil, errors.ErrNilTopicsChecker
	}

	return &sovereignRunTypeComponentsFactory{
		runTypeComponentsFactory: args.RunTypeComponentsFactory,
		sovConfig:                args.Config,
		dataCodec:                args.DataCodec,
		topicsChecker:            args.TopicsChecker,
	}, nil
}

// Create creates the runType components
func (rcf *sovereignRunTypeComponentsFactory) Create() (*runTypeComponents, error) {
	rtc, err := rcf.runTypeComponentsFactory.Create()
	if err != nil {
		return nil, err
	}

	blockProcessorFactory, err := block.NewSovereignBlockProcessorFactory(rtc.blockProcessorCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignBlockProcessorFactory failed: %w", err)
	}

	forkDetectorFactory, err := sync.NewSovereignForkDetectorFactory(rtc.forkDetectorCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignForkDetectorFactory failed: %w", err)
	}

	blockTrackerFactory, err := track.NewSovereignBlockTrackerFactory(rtc.blockTrackerCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignBlockTrackerFactory failed: %w", err)
	}

	requestHandlerFactory, err := requestHandlers.NewSovereignResolverRequestHandlerFactory(rtc.requestHandlerCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignResolverRequestHandlerFactory failed: %w", err)
	}

	headerValidatorFactory, err := block.NewSovereignHeaderValidatorFactory(rtc.headerValidatorCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignHeaderValidatorFactory failed: %w", err)
	}

	transactionCoordinatorFactory, err := coordinator.NewSovereignTransactionCoordinatorFactory(rtc.transactionCoordinatorCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignTransactionCoordinatorFactory failed: %w", err)
	}

	validatorStatisticsProcessorFactory, err := peer.NewSovereignValidatorStatisticsProcessorFactory(rtc.validatorStatisticsProcessorCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignValidatorStatisticsProcessorFactory failed: %w", err)
	}

	scProcessorCreator, err := processorV2.NewSovereignSCProcessFactory(rtc.scProcessorCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignSCProcessFactory failed: %w", err)
	}

	scResultPreProcessorCreator, err := preprocess.NewSovereignSmartContractResultPreProcessorFactory(rtc.scResultPreProcessorCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignSmartContractResultPreProcessorFactory failed: %w", err)
	}

	sovVMContextCreator := systemSmartContracts.NewOneShardSystemVMEEICreator()
	rtc.vmContainerMetaFactory, err = factoryVm.NewVmContainerMetaFactory(sovVMContextCreator)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewVmContainerMetaFactory failed: %w", err)
	}

	sovereignVmContainerShardCreator, err := factoryVm.NewSovereignVmContainerShardFactory(rtc.vmContainerMetaFactory, rtc.vmContainerShardFactory)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignVmContainerShardFactory failed: %w", err)
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
	sovereignAccountsParser, err := parsing.NewSovereignAccountsParser(accountsParser)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignAccountsParser failed: %w", err)
	}

	accountsCreator, err := factory.NewSovereignAccountCreator(factory.ArgsSovereignAccountCreator{
		ArgsAccountCreator: factory.ArgsAccountCreator{
			Hasher:              rcf.coreComponents.Hasher(),
			Marshaller:          rcf.coreComponents.InternalMarshalizer(),
			EnableEpochsHandler: rcf.coreComponents.EnableEpochsHandler(),
		},
		BaseTokenID: rcf.sovConfig.GenesisConfig.NativeDCDT,
	})
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignAccountCreator failed: %w", err)
	}

	expiryTime := time.Second * time.Duration(rcf.sovConfig.OutgoingSubscribedEvents.TimeToWaitForUnconfirmedOutGoingOperationInSeconds)

	sovHeaderSigVerifier, err := headerCheck.NewSovereignHeaderSigVerifier(rcf.cryptoComponents.BlockSigner())
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignHeaderSigVerifier failed: %w", err)
	}
	err = rtc.extraHeaderSigVerifierHolder.RegisterExtraHeaderSigVerifier(sovHeaderSigVerifier)
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - RegisterExtraHeaderSigVerifier failed: %w", err)
	}

	apiRewardTxHandler, err := transactionAPI.NewSovereignAPIRewardsHandler(rcf.coreComponents.AddressPubKeyConverter())
	if err != nil {
		return nil, fmt.Errorf("sovereignRunTypeComponentsFactory - NewSovereignAPIRewardsHandler failed: %w", err)
	}

	headerVersionHandler, err := rcf.createHeaderVersionHandler()
	if err != nil {
		return nil, fmt.Errorf("runTypeComponentsFactory - createHeaderVersionHandler failed: %w", err)
	}
	versionedHeaderFactory, err := factoryBlock.NewSovereignShardHeaderFactory(headerVersionHandler)
	if err != nil {
		return nil, fmt.Errorf("runTypeComponentsFactory - NewSovereignShardHeaderFactory failed: %w", err)
	}

	return &runTypeComponents{
		blockChainHookHandlerCreator:            hooks.NewSovereignBlockChainHookFactory(),
		epochStartBootstrapperCreator:           bootstrap.NewSovereignEpochStartBootstrapperFactory(),
		bootstrapperFromStorageCreator:          storageBootstrap.NewSovereignShardStorageBootstrapperFactory(),
		bootstrapperCreator:                     storageBootstrap.NewSovereignShardBootstrapFactory(),
		blockProcessorCreator:                   blockProcessorFactory,
		forkDetectorCreator:                     forkDetectorFactory,
		blockTrackerCreator:                     blockTrackerFactory,
		requestHandlerCreator:                   requestHandlerFactory,
		headerValidatorCreator:                  headerValidatorFactory,
		scheduledTxsExecutionCreator:            preprocess.NewSovereignScheduledTxsExecutionFactory(),
		transactionCoordinatorCreator:           transactionCoordinatorFactory,
		validatorStatisticsProcessorCreator:     validatorStatisticsProcessorFactory,
		additionalStorageServiceCreator:         storageFactory.NewSovereignAdditionalStorageServiceFactory(),
		scProcessorCreator:                      scProcessorCreator,
		scResultPreProcessorCreator:             scResultPreProcessorCreator,
		consensusModel:                          consensus.ConsensusModelV2,
		vmContainerMetaFactory:                  rtc.vmContainerMetaFactory,
		vmContainerShardFactory:                 sovereignVmContainerShardCreator,
		accountsParser:                          sovereignAccountsParser,
		accountsCreator:                         accountsCreator,
		vmContextCreator:                        sovVMContextCreator,
		outGoingOperationsPoolHandler:           sovereignFactory.NewOutGoingOperationPool(expiryTime),
		dataCodecHandler:                        rcf.dataCodec,
		topicsCheckerHandler:                    rcf.topicsChecker,
		shardCoordinatorCreator:                 sharding.NewSovereignShardCoordinatorFactory(),
		nodesCoordinatorWithRaterFactoryCreator: nodesCoord.NewSovereignIndexHashedNodesCoordinatorWithRaterFactory(),
		requestersContainerFactoryCreator:       requesterscontainer.NewSovereignShardRequestersContainerFactoryCreator(),
		interceptorsContainerFactoryCreator:     interceptorscontainer.NewSovereignShardInterceptorsContainerFactoryCreator(),
		shardResolversContainerFactoryCreator:   resolverscontainer.NewSovereignShardResolversContainerFactoryCreator(),
		txPreProcessorCreator:                   preprocess.NewSovereignTxPreProcessorCreator(),
		extraHeaderSigVerifierHolder:            rtc.extraHeaderSigVerifierHolder,
		genesisBlockCreatorFactory:              processComp.NewSovereignGenesisBlockCreatorFactory(),
		genesisMetaBlockCheckerCreator:          processComp.NewSovereignGenesisMetaBlockChecker(),
		nodesSetupCheckerFactory:                rtc.nodesSetupCheckerFactory,
		epochStartTriggerFactory:                epochStartTrigger.NewSovereignEpochStartTriggerFactory(),
		latestDataProviderFactory:               latestData.NewSovereignLatestDataProviderFactory(),
		scToProtocolFactory:                     scToProtocol.NewSovereignStakingToPeerFactory(),
		validatorInfoCreatorFactory:             metachain.NewSovereignValidatorInfoCreatorFactory(),
		apiProcessorCompsCreatorHandler:         api.NewSovereignAPIProcessorCompsCreator(),
		endOfEpochEconomicsFactoryHandler:       metachain.NewSovereignEconomicsFactory(),
		rewardsCreatorFactory:                   metachain.NewSovereignRewardsCreatorFactory(),
		systemSCProcessorFactory:                metachain.NewSovereignSysSCFactory(),
		preProcessorsContainerFactoryCreator:    procSovereign.NewSovereignPreProcessorContainerFactoryCreator(),
		dataRetrieverContainersSetter:           dataRetriever.NewSovereignDataRetrieverContainerSetter(),
		shardMessengerFactory:                   broadcastFactory.NewSovereignShardChainMessengerFactory(),
		exportHandlerFactoryCreator:             updateFactory.NewSovereignExportHandlerFactoryCreator(),
		validatorAccountsSyncerFactoryHandler:   syncerFactory.NewSovereignValidatorAccountsSyncerFactory(),
		shardRequestersContainerCreatorHandler:  storageRequestFactory.NewSovereignShardRequestersContainerCreator(),
		apiRewardTxHandler:                      apiRewardTxHandler,
		outportDataProviderFactory:              outportFactory.NewSovereignOutportDataProviderFactory(),
		delegatedListFactoryHandler:             trieIteratorsFactory.NewSovereignDelegatedListProcessorFactory(),
		directStakedListFactoryHandler:          trieIteratorsFactory.NewSovereignDirectStakedListProcessorFactory(),
		totalStakedValueFactoryHandler:          trieIteratorsFactory.NewSovereignTotalStakedValueProcessorFactory(),
		versionedHeaderFactory:                  versionedHeaderFactory,
	}, nil
}
