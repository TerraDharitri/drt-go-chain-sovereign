package transactionEvaluator

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	datafield "github.com/TerraDharitri/drt-go-chain-vm-common/parsers/dataField"
	"github.com/TerraDharitri/drt-go-chain/state"
)

// TransactionProcessor defines the operations needed to be done by a transaction processor
type TransactionProcessor interface {
	ProcessTransaction(transaction *transaction.Transaction) (vmcommon.ReturnCode, error)
	VerifyTransaction(transaction *transaction.Transaction) error
	GetSenderAndReceiverAccounts(transaction *transaction.Transaction) (state.UserAccountHandler, state.UserAccountHandler, error)
	IsInterfaceNil() bool
}

// DataFieldParser defines what a data field parser should be able to do
type DataFieldParser interface {
	Parse(dataField []byte, sender, receiver []byte, numOfShards uint32) *datafield.ResponseParseData
}
