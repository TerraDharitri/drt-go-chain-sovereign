package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/accounts"
	"github.com/TerraDharitri/drt-go-chain/state/parsers"
	"github.com/TerraDharitri/drt-go-chain/state/trackableDataTrie"
)

// ArgsSovereignAccountCreator defines needed arguments for a sovereign account creator
type ArgsSovereignAccountCreator struct {
	ArgsAccountCreator
	BaseTokenID string
}

// sovereignAccountCreator has method to create a new account
type sovereignAccountCreator struct {
	hasher              hashing.Hasher
	marshaller          marshal.Marshalizer
	dcdtAsBalance       state.DCDTAsBalanceHandler
	enableEpochsHandler common.EnableEpochsHandler
}

// NewSovereignAccountCreator creates a new instance of AccountCreator
func NewSovereignAccountCreator(args ArgsSovereignAccountCreator) (state.AccountFactory, error) {
	if check.IfNil(args.Hasher) {
		return nil, errors.ErrNilHasher
	}
	if check.IfNil(args.Marshaller) {
		return nil, errors.ErrNilMarshalizer
	}
	if check.IfNil(args.EnableEpochsHandler) {
		return nil, errors.ErrNilEnableEpochsHandler
	}

	dcdtAsBalance, err := accounts.NewDCDTAsBalance(args.BaseTokenID, args.Marshaller)
	if err != nil {
		return nil, err
	}

	return &sovereignAccountCreator{
		hasher:              args.Hasher,
		marshaller:          args.Marshaller,
		enableEpochsHandler: args.EnableEpochsHandler,
		dcdtAsBalance:       dcdtAsBalance,
	}, nil
}

// CreateAccount calls the new Account creator and returns the result
func (s *sovereignAccountCreator) CreateAccount(address []byte) (vmcommon.AccountHandler, error) {
	tdt, err := trackableDataTrie.NewTrackableDataTrie(address, s.hasher, s.marshaller, s.enableEpochsHandler)
	if err != nil {
		return nil, err
	}

	dataTrieLeafParser, err := parsers.NewDataTrieLeafParser(address, s.marshaller, s.enableEpochsHandler)
	if err != nil {
		return nil, err
	}

	return accounts.NewSovereignAccount(address, tdt, dataTrieLeafParser, s.dcdtAsBalance)
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *sovereignAccountCreator) IsInterfaceNil() bool {
	return s == nil
}
