package accounts

import (
	"errors"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/stretchr/testify/require"

	errorsDrt "github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/testscommon/marshallerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/trie"
)

const baseTokenID = "WREWA-bd4d79"
const prefixedBaseTokenID = "sov-DCDT-a1b2c3"

func TestNewDCDTAsBalance(t *testing.T) {
	t.Parallel()

	t.Run("empty base token, should return error", func(t *testing.T) {
		dcdtBalance, err := NewDCDTAsBalance("", &marshallerMock.MarshalizerMock{})
		require.Equal(t, errorsDrt.ErrEmptyBaseToken, err)
		require.Nil(t, dcdtBalance)
	})
	t.Run("invalid prefixed base token, should return error", func(t *testing.T) {
		dcdtBalance, err := NewDCDTAsBalance("svn12-ABC-1a2f3f", &marshallerMock.MarshalizerMock{})
		require.Equal(t, errorsDrt.ErrInvalidBaseToken, err)
		require.Nil(t, dcdtBalance)
	})
	t.Run("invalid base token, should return error", func(t *testing.T) {
		dcdtBalance, err := NewDCDTAsBalance("MvX-1c4f2a", &marshallerMock.MarshalizerMock{})
		require.Equal(t, errorsDrt.ErrInvalidBaseToken, err)
		require.Nil(t, dcdtBalance)
	})
	t.Run("nil marshaller, should return error", func(t *testing.T) {
		dcdtBalance, err := NewDCDTAsBalance(baseTokenID, nil)
		require.Equal(t, errorsDrt.ErrNilMarshalizer, err)
		require.Nil(t, dcdtBalance)
	})
	t.Run("should work", func(t *testing.T) {
		dcdtBalance, err := NewDCDTAsBalance(baseTokenID, &marshallerMock.MarshalizerMock{})
		require.Nil(t, err)
		require.False(t, dcdtBalance.IsInterfaceNil())
	})
	t.Run("prefixed base token, should work", func(t *testing.T) {
		dcdtBalance, err := NewDCDTAsBalance(prefixedBaseTokenID, &marshallerMock.MarshalizerMock{})
		require.Nil(t, err)
		require.False(t, dcdtBalance.IsInterfaceNil())
	})
}

func TestDcdtAsBalance_getDCDTData(t *testing.T) {
	t.Parallel()

	testGetDcdtData(t, baseTokenID)
	testGetDcdtData(t, prefixedBaseTokenID)
}

func testGetDcdtData(t *testing.T, baseTokenID string) {
	t.Run("no value stored in account for dcdt, should return empty token", func(t *testing.T) {
		dcdtBalance, _ := NewDCDTAsBalance(baseTokenID, &marshallerMock.MarshalizerMock{})
		accHandler := &trie.DataTrieTrackerStub{
			RetrieveValueCalled: func(key []byte) ([]byte, uint32, error) {
				require.Equal(t, []byte(baseDCDTKeyPrefix+baseTokenID), key)
				return nil, 0, errors.New("error retrieving value")
			},
		}

		dcdtData, err := dcdtBalance.getDCDTData(accHandler)
		require.Nil(t, err)
		require.Equal(t, createEmptyDCDT(), dcdtData)
	})

	t.Run("empty buffer when retrieving data, should return empty token", func(t *testing.T) {
		dcdtBalance, _ := NewDCDTAsBalance(baseTokenID, &marshallerMock.MarshalizerMock{})
		accHandler := &trie.DataTrieTrackerStub{
			RetrieveValueCalled: func(key []byte) ([]byte, uint32, error) {
				return nil, 0, nil
			},
		}

		dcdtData, err := dcdtBalance.getDCDTData(accHandler)
		require.Nil(t, err)
		require.Equal(t, createEmptyDCDT(), dcdtData)
	})

	t.Run("cannot unmarshall, should return error", func(t *testing.T) {
		errUnmarshall := errors.New("cannot unmarshall")
		marshaller := &marshallerMock.MarshalizerStub{
			UnmarshalCalled: func(obj interface{}, buff []byte) error {
				return errUnmarshall
			},
		}
		dcdtBalance, _ := NewDCDTAsBalance(baseTokenID, marshaller)

		dcdtData, err := dcdtBalance.getDCDTData(&trie.DataTrieTrackerStub{})
		require.Nil(t, dcdtData)
		require.Equal(t, errUnmarshall, err)
	})

	t.Run("empty value when unmarshalling, should fill mandatory fields", func(t *testing.T) {
		marshaller := &marshallerMock.MarshalizerStub{
			UnmarshalCalled: func(obj interface{}, buff []byte) error {
				expectedObj := obj.(*dcdt.DCDigitalToken)
				expectedObj.Properties = []byte("properties")

				return nil
			},
		}
		dcdtBalance, _ := NewDCDTAsBalance(baseTokenID, marshaller)

		dcdtData, err := dcdtBalance.getDCDTData(&trie.DataTrieTrackerStub{})
		require.Equal(t, &dcdt.DCDigitalToken{
			Type:       uint32(core.Fungible),
			Value:      big.NewInt(0),
			Properties: []byte("properties"),
		}, dcdtData)
		require.Nil(t, err)
	})
}

func TestDcdtAsBalance_GetBalance(t *testing.T) {
	t.Parallel()

	testGetBalance(t, baseTokenID)
	testGetBalance(t, prefixedBaseTokenID)
}

func testGetBalance(t *testing.T, baseTokenID string) {
	t.Run("could not load balance, should return 0 value", func(t *testing.T) {
		marshaller := &marshallerMock.MarshalizerStub{
			UnmarshalCalled: func(obj interface{}, buff []byte) error {
				return errors.New("cannot unmarshall")
			},
		}
		dcdtBalance, _ := NewDCDTAsBalance(baseTokenID, marshaller)
		balance := dcdtBalance.GetBalance(&trie.DataTrieTrackerStub{})
		require.Equal(t, big.NewInt(0), balance)
	})

	t.Run("should work", func(t *testing.T) {
		expectedBalance := big.NewInt(444)
		marshaller := &marshallerMock.MarshalizerStub{
			UnmarshalCalled: func(obj interface{}, buff []byte) error {
				expectedObj := obj.(*dcdt.DCDigitalToken)
				expectedObj.Value = expectedBalance

				return nil
			},
		}
		dcdtBalance, _ := NewDCDTAsBalance(baseTokenID, marshaller)
		balance := dcdtBalance.GetBalance(&trie.DataTrieTrackerStub{})
		require.Equal(t, expectedBalance, balance)
	})
}

func TestDcdtAsBalance_AddToBalance(t *testing.T) {
	t.Parallel()

	testAddToBalance(t, baseTokenID)
	testAddToBalance(t, prefixedBaseTokenID)
}

func testAddToBalance(t *testing.T, baseTokenID string) {
	currentBalance := &dcdt.DCDigitalToken{
		Value: big.NewInt(123),
	}
	marshaller := &marshallerMock.MarshalizerMock{}
	dcdtBalance, _ := NewDCDTAsBalance(baseTokenID, marshaller)

	newValue := big.NewInt(321)
	expectedNewBalance := big.NewInt(0).Add(currentBalance.Value, newValue)
	wasBalanceSaved := false
	accHandler := &trie.DataTrieTrackerStub{
		RetrieveValueCalled: func(key []byte) ([]byte, uint32, error) {
			require.Equal(t, []byte(baseDCDTKeyPrefix+baseTokenID), key)

			storedValue, err := marshaller.Marshal(currentBalance)
			require.Nil(t, err)

			return storedValue, 0, nil
		},
		SaveKeyValueCalled: func(key []byte, value []byte) error {
			expectedBalance := &dcdt.DCDigitalToken{
				Value: expectedNewBalance,
			}
			expectedBalanceMarshalledData, err := marshaller.Marshal(expectedBalance)
			require.Nil(t, err)
			require.Equal(t, expectedBalanceMarshalledData, value)
			require.Equal(t, []byte(baseDCDTKeyPrefix+baseTokenID), key)

			wasBalanceSaved = true
			return nil
		},
	}

	err := dcdtBalance.AddToBalance(accHandler, newValue)
	require.Nil(t, err)
	require.True(t, wasBalanceSaved)

	err = dcdtBalance.AddToBalance(accHandler, big.NewInt(-4412))
	require.Equal(t, errorsDrt.ErrInsufficientFunds, err)
}

func TestDcdtAsBalance_SubFromBalance(t *testing.T) {
	t.Parallel()

	testSubFromBalance(t, baseTokenID)
	testSubFromBalance(t, prefixedBaseTokenID)
}

func testSubFromBalance(t *testing.T, baseTokenID string) {
	currentBalance := &dcdt.DCDigitalToken{
		Value: big.NewInt(121),
	}
	marshaller := &marshallerMock.MarshalizerMock{}
	dcdtBalance, _ := NewDCDTAsBalance(baseTokenID, marshaller)

	subBalance := big.NewInt(22)
	expectedNewBalance := big.NewInt(0).Sub(currentBalance.Value, subBalance)
	wasBalanceSaved := false
	accHandler := &trie.DataTrieTrackerStub{
		RetrieveValueCalled: func(key []byte) ([]byte, uint32, error) {
			require.Equal(t, []byte(baseDCDTKeyPrefix+baseTokenID), key)

			storedValue, err := marshaller.Marshal(currentBalance)
			require.Nil(t, err)

			return storedValue, 0, nil
		},
		SaveKeyValueCalled: func(key []byte, value []byte) error {
			expectedBalance := &dcdt.DCDigitalToken{
				Value: expectedNewBalance,
			}
			expectedBalanceMarshalledData, err := marshaller.Marshal(expectedBalance)
			require.Nil(t, err)
			require.Equal(t, expectedBalanceMarshalledData, value)
			require.Equal(t, []byte(baseDCDTKeyPrefix+baseTokenID), key)

			wasBalanceSaved = true
			return nil
		},
	}

	err := dcdtBalance.SubFromBalance(accHandler, subBalance)
	require.Nil(t, err)
	require.True(t, wasBalanceSaved)

	err = dcdtBalance.SubFromBalance(accHandler, big.NewInt(122))
	require.Equal(t, errorsDrt.ErrInsufficientFunds, err)
}
