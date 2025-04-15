package accounts

import (
	"math/big"
	"testing"

	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	errorsDrt "github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/testscommon/sovereign"
	"github.com/TerraDharitri/drt-go-chain/testscommon/trie"
	"github.com/stretchr/testify/require"
)

func TestNewSovereignAccount(t *testing.T) {
	t.Parallel()

	t.Run("nil address, should return error", func(t *testing.T) {
		sovAcc, err := NewSovereignAccount(
			nil,
			&trie.DataTrieTrackerStub{},
			&trie.TrieLeafParserStub{},
			&sovereign.DCDTAsBalanceHandlerMock{},
		)
		require.Nil(t, sovAcc)
		require.Equal(t, errorsDrt.ErrNilAddress, err)
	})
	t.Run("nil data trie, should return error", func(t *testing.T) {
		sovAcc, err := NewSovereignAccount(
			[]byte("address"),
			nil,
			&trie.TrieLeafParserStub{},
			&sovereign.DCDTAsBalanceHandlerMock{},
		)
		require.Nil(t, sovAcc)
		require.Equal(t, errorsDrt.ErrNilTrackableDataTrie, err)
	})
	t.Run("nil leaf parser, should return error", func(t *testing.T) {
		sovAcc, err := NewSovereignAccount(
			[]byte("address"),
			&trie.DataTrieTrackerStub{},
			nil,
			&sovereign.DCDTAsBalanceHandlerMock{},
		)
		require.Nil(t, sovAcc)
		require.Equal(t, errorsDrt.ErrNilTrieLeafParser, err)
	})
	t.Run("nil dcdt balance, should return error", func(t *testing.T) {
		sovAcc, err := NewSovereignAccount(
			[]byte("address"),
			&trie.DataTrieTrackerStub{},
			&trie.TrieLeafParserStub{},
			nil,
		)
		require.Nil(t, sovAcc)
		require.Equal(t, errorsDrt.ErrNilDCDTAsBalanceHandler, err)
	})
	t.Run("should work", func(t *testing.T) {
		sovAcc, err := NewSovereignAccount(
			[]byte("address"),
			&trie.DataTrieTrackerStub{},
			&trie.TrieLeafParserStub{},
			&sovereign.DCDTAsBalanceHandlerMock{},
		)
		require.Nil(t, err)
		require.False(t, sovAcc.IsInterfaceNil())
	})
}

func TestSovereignAccount_AddToBalance_SubFromBalance_GetBalance(t *testing.T) {
	t.Parallel()

	wasGetBalanceCalled := false
	wasAddBalanceCalled := false
	wasSubBalanceCalled := false
	dcdtBalance := &sovereign.DCDTAsBalanceHandlerMock{
		GetBalanceCalled: func(accountDataHandler vmcommon.AccountDataHandler) *big.Int {
			wasGetBalanceCalled = true
			return nil
		},
		AddToBalanceCalled: func(accountDataHandler vmcommon.AccountDataHandler, value *big.Int) error {
			wasAddBalanceCalled = true
			return nil
		},
		SubFromBalanceCalled: func(accountDataHandler vmcommon.AccountDataHandler, value *big.Int) error {
			wasSubBalanceCalled = true
			return nil
		},
	}
	sovAcc, _ := NewSovereignAccount(
		[]byte("address"),
		&trie.DataTrieTrackerStub{},
		&trie.TrieLeafParserStub{},
		dcdtBalance,
	)

	require.Nil(t, sovAcc.GetBalance())
	require.Nil(t, sovAcc.AddToBalance(nil))
	require.Nil(t, sovAcc.SubFromBalance(nil))

	require.True(t, wasGetBalanceCalled)
	require.True(t, wasAddBalanceCalled)
	require.True(t, wasSubBalanceCalled)
}
