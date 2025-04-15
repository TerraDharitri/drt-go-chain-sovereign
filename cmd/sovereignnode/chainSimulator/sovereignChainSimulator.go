package chainSimulator

import (
	"path"

	sovCommon "github.com/TerraDharitri/drt-go-chain/cmd/sovereignnode/chainSimulator/common"
	sovChainSimConfig "github.com/TerraDharitri/drt-go-chain/cmd/sovereignnode/chainSimulator/configs"
	sovereignConfig "github.com/TerraDharitri/drt-go-chain/cmd/sovereignnode/config"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/factory/runType"
	chainSimulatorIntegrationTests "github.com/TerraDharitri/drt-go-chain/integrationTests/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator"
	chainSimulatorConfigs "github.com/TerraDharitri/drt-go-chain/node/chainSimulator/configs"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/dtos"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/sovereign/incomingHeader"
)

const (
	numOfShards = 1
)

// ArgsSovereignChainSimulator holds the arguments for sovereign chain simulator
type ArgsSovereignChainSimulator struct {
	SovereignConfigPath string
	*chainSimulator.ArgsChainSimulator
}

// NewSovereignChainSimulator will create a new instance of sovereign chain simulator
func NewSovereignChainSimulator(args ArgsSovereignChainSimulator) (chainSimulatorIntegrationTests.ChainSimulator, error) {
	args.NumOfShards = numOfShards

	alterConfigs := args.AlterConfigsFunction
	configs, err := loadSovereignConfigs(args.SovereignConfigPath)
	if err != nil {
		return nil, err
	}

	args.AlterConfigsFunction = func(cfg *config.Configs) {
		cfg.EconomicsConfig = configs.EconomicsConfig
		cfg.EpochConfig = configs.EpochConfig
		cfg.GeneralConfig.SovereignConfig = *configs.SovereignExtraConfig
		cfg.GeneralConfig.VirtualMachine.Execution.WasmVMVersions = []config.WasmVMVersionByEpoch{{StartEpoch: 0, Version: "v1.5"}}
		cfg.GeneralConfig.VirtualMachine.Querying.WasmVMVersions = []config.WasmVMVersionByEpoch{{StartEpoch: 0, Version: "v1.5"}}
		cfg.SystemSCConfig.DCDTSystemSCConfig.DCDTPrefix = "sov"
		cfg.GeneralConfig.Versions.VersionsByEpochs = []config.VersionByEpochs{{StartEpoch: 0, Version: string(process.SovereignHeaderVersion)}}

		if alterConfigs != nil {
			alterConfigs(cfg)
			configs.SovereignExtraConfig = &cfg.GeneralConfig.SovereignConfig
		}
	}

	args.CreateRunTypeCoreComponents = func() (factory.RunTypeCoreComponentsHolder, error) {
		return createSovereignRunTypeCoreComponents(*configs.SovereignEpochConfig)
	}
	args.CreateIncomingHeaderSubscriber = func(config config.WebSocketConfig, dataPool dataRetriever.PoolsHolder, mainChainNotarizationStartRound uint64, runTypeComponents factory.RunTypeComponentsHolder) (process.IncomingHeaderSubscriber, error) {
		return incomingHeader.CreateIncomingHeaderProcessor(config, dataPool, mainChainNotarizationStartRound, runTypeComponents)
	}
	if args.CreateRunTypeComponents == nil {
		args.CreateRunTypeComponents = func(args runType.ArgsRunTypeComponents) (factory.RunTypeComponentsHolder, error) {
			return sovCommon.CreateSovereignRunTypeComponents(args, *configs.SovereignExtraConfig)
		}
	}
	args.NodeFactory = node.NewSovereignNodeFactory(configs.SovereignExtraConfig.GenesisConfig.NativeDCDT)
	args.ChainProcessorFactory = NewSovereignChainHandlerFactory()
	args.GenerateGenesisFile = func(args chainSimulatorConfigs.ArgsChainSimulatorConfigs, configs *config.Configs) (*dtos.InitialWalletKeys, error) {
		return sovChainSimConfig.GenerateSovereignGenesisFile(args, configs)
	}

	return chainSimulator.NewSovereignChainSimulator(*args.ArgsChainSimulator)
}

// loadSovereignConfigs loads sovereign configs
func loadSovereignConfigs(configsPath string) (*sovereignConfig.SovereignConfig, error) {
	epochConfig, err := common.LoadEpochConfig(path.Join(configsPath, "enableEpochs.toml"))
	if err != nil {
		return nil, err
	}

	economicsConfig, err := common.LoadEconomicsConfig(path.Join(configsPath, "economics.toml"))
	if err != nil {
		return nil, err
	}

	sovereignExtraConfig, err := sovereignConfig.LoadSovereignGeneralConfig(path.Join(configsPath, "sovereignConfig.toml"))
	if err != nil {
		return nil, err
	}

	sovereignEpochConfig, err := sovereignConfig.LoadSovereignEpochConfig(path.Join(configsPath, "enableEpochs.toml"))
	if err != nil {
		return nil, err
	}

	return &sovereignConfig.SovereignConfig{
		Configs: &config.Configs{
			EpochConfig:     epochConfig,
			EconomicsConfig: economicsConfig,
		},
		SovereignExtraConfig: sovereignExtraConfig,
		SovereignEpochConfig: sovereignEpochConfig,
	}, nil
}

func createSovereignRunTypeCoreComponents(sovereignEpochConfig config.SovereignEpochConfig) (factory.RunTypeCoreComponentsHolder, error) {
	sovereignRunTypeCoreComponentsFactory := runType.NewSovereignRunTypeCoreComponentsFactory(sovereignEpochConfig)
	managedRunTypeCoreComponents, err := runType.NewManagedRunTypeCoreComponents(sovereignRunTypeCoreComponentsFactory)
	if err != nil {
		return nil, err
	}
	err = managedRunTypeCoreComponents.Create()
	if err != nil {
		return nil, err
	}

	return managedRunTypeCoreComponents, nil
}
