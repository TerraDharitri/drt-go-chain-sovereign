package factory

import (
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
)

const (
	// BootstrapComponentsName is the bootstrap components identifier
	BootstrapComponentsName = "managedBootstrapComponents"
	// ConsensusComponentsName is the consensus components identifier
	ConsensusComponentsName = "managedConsensusComponents"
	// CoreComponentsName is the core components identifier
	CoreComponentsName = "managedCoreComponents"
	// StatusCoreComponentsName is the status core components identifier
	StatusCoreComponentsName = "managedStatusCoreComponents"
	// CryptoComponentsName is the crypto components identifier
	CryptoComponentsName = "managedCryptoComponents"
	// DataComponentsName is the data components identifier
	DataComponentsName = "managedDataComponents"
	// HeartbeatV2ComponentsName is the heartbeat V2 components identifier
	HeartbeatV2ComponentsName = "managedHeartbeatV2Components"
	// NetworkComponentsName is the network components identifier
	NetworkComponentsName = "managedNetworkComponents"
	// ProcessComponentsName is the process components identifier
	ProcessComponentsName = "managedProcessComponents"
	// StateComponentsName is the state components identifier
	StateComponentsName = "managedStateComponents"
	// StatusComponentsName is the status components identifier
	StatusComponentsName = "managedStatusComponents"
	// RunTypeComponentsName is the runType components identifier
	RunTypeComponentsName = "managedRunTypeComponents"
	// RunTypeCoreComponentsName is the runType core components identifier
	RunTypeCoreComponentsName = "managedRunTypeCoreComponents"
)

// ArgsEpochStartTrigger is a struct placeholder for arguments needed to create an epoch start trigger
type ArgsEpochStartTrigger struct {
	RequestHandler             epochStart.RequestHandler
	CoreData                   CoreComponentsHolder
	BootstrapComponents        BootstrapComponentsHolder
	DataComps                  DataComponentsHolder
	StatusCoreComponentsHolder StatusCoreComponentsHolder
	RunTypeComponentsHolder    RunTypeComponentsHolder
	Config                     config.Config
}

// ArgsExporter is the argument structure to create a new exporter
type ArgsExporter struct {
	CoreComponents                   process.CoreComponentsHolder
	CryptoComponents                 process.CryptoComponentsHolder
	StatusCoreComponents             process.StatusCoreComponentsHolder
	NetworkComponents                NetworkComponentsHolder
	HeaderValidator                  epochStart.HeaderValidator
	DataPool                         dataRetriever.PoolsHolder
	StorageService                   dataRetriever.StorageService
	RequestHandler                   process.RequestHandler
	ShardCoordinator                 sharding.Coordinator
	ActiveAccountsDBs                map[state.AccountsDbIdentifier]state.AccountsAdapter
	ExistingResolvers                dataRetriever.ResolversContainer
	ExistingRequesters               dataRetriever.RequestersContainer
	ExportFolder                     string
	ExportTriesStorageConfig         config.StorageConfig
	ExportStateStorageConfig         config.StorageConfig
	ExportStateKeysConfig            config.StorageConfig
	MaxTrieLevelInMemory             uint
	WhiteListHandler                 process.WhiteListHandler
	WhiteListerVerifiedTxs           process.WhiteListHandler
	MainInterceptorsContainer        process.InterceptorsContainer
	FullArchiveInterceptorsContainer process.InterceptorsContainer
	NodesCoordinator                 nodesCoordinator.NodesCoordinator
	HeaderSigVerifier                process.InterceptedHeaderSigVerifier
	HeaderIntegrityVerifier          process.HeaderIntegrityVerifier
	ValidityAttester                 process.ValidityAttester
	RoundHandler                     process.RoundHandler
	InterceptorDebugConfig           config.InterceptorResolverDebugConfig
	MaxHardCapForMissingNodes        int
	NumConcurrentTrieSyncers         int
	TrieSyncerVersion                int
	CheckNodesOnDisk                 bool
	NodeOperationMode                common.NodeOperation

	ShardCoordinatorFactory sharding.ShardCoordinatorFactory
}
