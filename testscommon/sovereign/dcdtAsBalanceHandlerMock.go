package sovereign

import (
	"math/big"

	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

// DCDTAsBalanceHandlerMock -
type DCDTAsBalanceHandlerMock struct {
	GetBalanceCalled     func(accountDataHandler vmcommon.AccountDataHandler) *big.Int
	AddToBalanceCalled   func(accountDataHandler vmcommon.AccountDataHandler, value *big.Int) error
	SubFromBalanceCalled func(accountDataHandler vmcommon.AccountDataHandler, value *big.Int) error
}

// GetBalance -
func (mock *DCDTAsBalanceHandlerMock) GetBalance(accountDataHandler vmcommon.AccountDataHandler) *big.Int {
	if mock.GetBalanceCalled != nil {
		return mock.GetBalanceCalled(accountDataHandler)
	}
	return nil
}

// AddToBalance -
func (mock *DCDTAsBalanceHandlerMock) AddToBalance(accountDataHandler vmcommon.AccountDataHandler, value *big.Int) error {
	if mock.AddToBalanceCalled != nil {
		return mock.AddToBalanceCalled(accountDataHandler, value)
	}
	return nil
}

// SubFromBalance -
func (mock *DCDTAsBalanceHandlerMock) SubFromBalance(accountDataHandler vmcommon.AccountDataHandler, value *big.Int) error {
	if mock.SubFromBalanceCalled != nil {
		return mock.SubFromBalanceCalled(accountDataHandler, value)
	}
	return nil
}

// IsInterfaceNil -
func (mock *DCDTAsBalanceHandlerMock) IsInterfaceNil() bool {
	return mock == nil
}
