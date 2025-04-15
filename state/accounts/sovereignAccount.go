package accounts

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/state"
)

var _ state.UserAccountHandler = (*sovereignAccount)(nil)

type sovereignAccount struct {
	*userAccount
	dcdtBalance state.DCDTAsBalanceHandler
}

// NewSovereignAccount creates a user account for sovereign shards
func NewSovereignAccount(
	address []byte,
	trackableDataTrie state.DataTrieTracker,
	trieLeafParser common.TrieLeafParser,
	dcdtBalance state.DCDTAsBalanceHandler,
) (*sovereignAccount, error) {
	if len(address) == 0 {
		return nil, errors.ErrNilAddress
	}
	if check.IfNil(trackableDataTrie) {
		return nil, errors.ErrNilTrackableDataTrie
	}
	if check.IfNil(trieLeafParser) {
		return nil, errors.ErrNilTrieLeafParser
	}
	if check.IfNil(dcdtBalance) {
		return nil, errors.ErrNilDCDTAsBalanceHandler
	}

	userAcc := &userAccount{
		UserAccountData: UserAccountData{
			DeveloperReward: big.NewInt(0),
			Balance:         big.NewInt(0),
			Address:         address,
		},
		dataTrieInteractor: trackableDataTrie,
		dataTrieLeafParser: trieLeafParser,
	}

	return &sovereignAccount{
		userAccount: userAcc,
		dcdtBalance: dcdtBalance,
	}, nil
}

// AddToBalance adds new value to balance
func (s *sovereignAccount) AddToBalance(value *big.Int) error {
	return s.dcdtBalance.AddToBalance(s.dataTrieInteractor, value)
}

// SubFromBalance subtracts new value from balance
func (s *sovereignAccount) SubFromBalance(value *big.Int) error {
	return s.dcdtBalance.SubFromBalance(s.dataTrieInteractor, value)
}

// GetBalance returns the actual balance from the account
func (s *sovereignAccount) GetBalance() *big.Int {
	return s.dcdtBalance.GetBalance(s.dataTrieInteractor)
}

// IsInterfaceNil checks if the underlying pointer is nil
func (s *sovereignAccount) IsInterfaceNil() bool {
	return s == nil
}
