package transaction

import "github.com/TerraDharitri/drt-go-chain/process"

// sovereignTransactionProcessor implements the transaction processor for sovereign shards
type sovereignTransactionProcessor struct {
	process.TransactionProcessor
}

// ArgsNewSovereignTxProcessor -
type ArgsNewSovereignTxProcessor struct {
}

// NewSovereignTransactionProcessor -
func NewSovereignTransactionProcessor(args ArgsNewSovereignTxProcessor) (*sovereignTransactionProcessor, error) {
	return nil, nil
}
