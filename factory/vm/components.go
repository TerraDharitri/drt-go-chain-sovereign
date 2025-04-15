package vm

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

// ArgsVmContainerFactory hold the argument needed for creating vm container
type ArgsVmContainerFactory struct {
	Config                     config.VirtualMachineConfig
	BlockGasLimit              uint64
	GasSchedule                core.GasScheduleNotifier
	EpochNotifier              process.EpochNotifier
	EnableEpochsHandler        common.EnableEpochsHandler
	WasmVMChangeLocker         common.Locker
	DCDTTransferParser         vmcommon.DCDTTransferParser
	BuiltInFunctions           vmcommon.BuiltInFunctionContainer
	BlockChainHook             process.BlockChainHookWithAccountsAdapter
	Hasher                     hashing.Hasher
	Economics                  process.EconomicsDataHandler
	MessageSignVerifier        vm.MessageSignVerifier
	NodesConfigProvider        vm.NodesConfigProvider
	Marshalizer                marshal.Marshalizer
	SystemSCConfig             *config.SystemSmartContractsConfig
	ValidatorAccountsDB        state.AccountsAdapter
	UserAccountsDB             state.AccountsAdapter
	ChanceComputer             nodesCoordinator.ChanceComputer
	ShardCoordinator           sharding.Coordinator
	PubkeyConv                 core.PubkeyConverter
	IsInHistoricalBalancesMode bool
	NodesCoordinator           vm.NodesCoordinator
}
