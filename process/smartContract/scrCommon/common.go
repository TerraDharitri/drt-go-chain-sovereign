package scrCommon

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

// ArgsNewSmartContractProcessor defines the arguments needed for new smart contract processor
type ArgsNewSmartContractProcessor struct {
	VmContainer         process.VirtualMachinesContainer
	ArgsParser          process.ArgumentsParser
	Hasher              hashing.Hasher
	Marshalizer         marshal.Marshalizer
	AccountsDB          state.AccountsAdapter
	BlockChainHook      process.BlockChainHookHandler
	BuiltInFunctions    vmcommon.BuiltInFunctionContainer
	PubkeyConv          core.PubkeyConverter
	ShardCoordinator    sharding.Coordinator
	ScrForwarder        process.IntermediateTransactionHandler
	TxFeeHandler        process.TransactionFeeHandler
	EconomicsFee        process.FeeHandler
	TxTypeHandler       process.TxTypeHandler
	GasHandler          process.GasHandler
	GasSchedule         core.GasScheduleNotifier
	TxLogsProcessor     process.TransactionLogProcessor
	BadTxForwarder      process.IntermediateTransactionHandler
	EnableRoundsHandler process.EnableRoundsHandler
	EnableEpochsHandler common.EnableEpochsHandler
	EnableEpochs        config.EnableEpochs
	VMOutputCacher      storage.Cacher
	WasmVMChangeLocker  common.Locker
	IsGenesisProcessing bool
	EpochNotifier       vmcommon.EpochNotifier
}

// ScrProcessingData is a struct placeholder for scr data to be processed after validation checks
type ScrProcessingData struct {
	Hash        []byte
	Snapshot    int
	Sender      state.UserAccountHandler
	Destination state.UserAccountHandler
}

// GetHash returns the hash
func (spd *ScrProcessingData) GetHash() []byte {
	return spd.Hash
}

// GetSnapshot returns the snapshot
func (spd *ScrProcessingData) GetSnapshot() int {
	return spd.Snapshot
}

// GetSender returns the sender
func (spd *ScrProcessingData) GetSender() state.UserAccountHandler {
	return spd.Sender
}

// GetDestination returns the destination
func (spd *ScrProcessingData) GetDestination() state.UserAccountHandler {
	return spd.Destination
}

// FindVMByScAddress is exported for use in all version of scr processors
func FindVMByScAddress(container process.VirtualMachinesContainer, scAddress []byte) (vmcommon.VMExecutionHandler, []byte, error) {
	vmType, err := vmcommon.ParseVMTypeFromContractAddress(scAddress)
	if err != nil {
		return nil, nil, err
	}

	vm, err := container.Get(vmType)
	if err != nil {
		return nil, nil, err
	}

	return vm, vmType, nil
}

// CreateExecutableCheckersMap creates a map of executable checker builtin functions
func CreateExecutableCheckersMap(builtinFunctions vmcommon.BuiltInFunctionContainer) map[string]ExecutableChecker {
	executableCheckers := make(map[string]ExecutableChecker)

	for key := range builtinFunctions.Keys() {
		builtinFunc, err := builtinFunctions.Get(key)
		if err != nil {
			continue
		}
		executableCheckerFunc, ok := builtinFunc.(ExecutableChecker)
		if !ok {
			continue
		}
		executableCheckers[key] = executableCheckerFunc
	}

	return executableCheckers
}
