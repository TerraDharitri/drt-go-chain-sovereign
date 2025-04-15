package scrCommon

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// TestSmartContractProcessor is a SmartContractProcessor used in integration tests
type TestSmartContractProcessor interface {
	process.SmartContractProcessorFacade
	GetCompositeTestError() error
	GetGasRemaining() uint64
	GetAllSCRs() []data.TransactionHandler
	CleanGasRefunded()
}

// ExecutableChecker is an interface for checking if a builtin function is executable
type ExecutableChecker interface {
	CheckIsExecutable(senderAddr []byte, value *big.Int, receiverAddr []byte, gasProvidedForCall uint64, arguments [][]byte) error
}

// SCRProcessorHandler defines a scr processor handler
type SCRProcessorHandler interface {
	process.SmartContractProcessor
	process.SmartContractResultProcessor
}

// SCProcessorCreator defines a scr processor creator
type SCProcessorCreator interface {
	CreateSCProcessor(args ArgsNewSmartContractProcessor) (SCRProcessorHandler, error)
	IsInterfaceNil() bool
}
