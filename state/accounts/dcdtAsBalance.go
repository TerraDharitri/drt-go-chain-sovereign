package accounts

import (
	"math/big"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"

	"github.com/TerraDharitri/drt-go-chain/errors"
)

const baseDCDTKeyPrefix = core.ProtectedKeyPrefix + core.DCDTKeyIdentifier

var log = logger.GetOrCreate("dcdt-as-balance")

type dcdtAsBalance struct {
	keyPrefix  []byte
	marshaller marshal.Marshalizer
}

// NewDCDTAsBalance creates the dcdtAsBalance component
func NewDCDTAsBalance(
	baseTokenID string,
	marshaller marshal.Marshalizer,
) (*dcdtAsBalance, error) {
	if check.IfNil(marshaller) {
		return nil, errors.ErrNilMarshalizer
	}
	err := validateBaseToken(baseTokenID)
	if err != nil {
		return nil, err
	}

	return &dcdtAsBalance{
		keyPrefix:  []byte(baseDCDTKeyPrefix + baseTokenID),
		marshaller: marshaller,
	}, nil
}

func validateBaseToken(baseTokenID string) error {
	if len(baseTokenID) == 0 {
		return errors.ErrEmptyBaseToken
	}

	if _, isValid := dcdt.IsValidPrefixedToken(baseTokenID); isValid {
		return nil
	}

	if vmcommon.ValidateToken([]byte(baseTokenID)) {
		return nil
	}

	return errors.ErrInvalidBaseToken
}

// GetBalance returns the native dcdt balance
func (e *dcdtAsBalance) GetBalance(accountDataHandler vmcommon.AccountDataHandler) *big.Int {
	dcdtData, err := e.getDCDTData(accountDataHandler)
	if err != nil {
		return big.NewInt(0)
	}

	return dcdtData.Value
}

// AddToBalance adds balance to the native dcdt balance
func (e *dcdtAsBalance) AddToBalance(accountDataHandler vmcommon.AccountDataHandler, value *big.Int) error {
	dcdtData, err := e.getDCDTData(accountDataHandler)
	if err != nil {
		return err
	}

	newBalance := big.NewInt(0).Add(dcdtData.Value, value)
	if newBalance.Cmp(zero) < 0 {
		return errors.ErrInsufficientFunds
	}

	dcdtData.Value.Set(newBalance)
	return e.saveDCDTData(accountDataHandler, dcdtData)
}

// SubFromBalance subtracts the value from the native dcdt balance
func (e *dcdtAsBalance) SubFromBalance(accountDataHandler vmcommon.AccountDataHandler, value *big.Int) error {
	dcdtData, err := e.getDCDTData(accountDataHandler)
	if err != nil {
		return err
	}

	newBalance := big.NewInt(0).Sub(dcdtData.Value, value)
	if newBalance.Cmp(zero) < 0 {
		return errors.ErrInsufficientFunds
	}

	dcdtData.Value.Set(newBalance)
	return e.saveDCDTData(accountDataHandler, dcdtData)
}

func (e *dcdtAsBalance) getDCDTData(accountDataHandler vmcommon.AccountDataHandler) (*dcdt.DCDigitalToken, error) {
	marshaledData, _, err := accountDataHandler.RetrieveValue(e.keyPrefix)
	if err != nil || len(marshaledData) == 0 {
		log.Trace("dcdtAsBalance.getDCDTData could not load account token", "error", err)
		return createEmptyDCDT(), nil
	}

	dcdtData := &dcdt.DCDigitalToken{}
	err = e.marshaller.Unmarshal(dcdtData, marshaledData)
	if err != nil {
		return nil, err
	}

	fillMandatoryFields(dcdtData)

	return dcdtData, nil
}

func createEmptyDCDT() *dcdt.DCDigitalToken {
	return &dcdt.DCDigitalToken{Value: big.NewInt(0), Type: uint32(core.Fungible)}
}

func fillMandatoryFields(dcdtData *dcdt.DCDigitalToken) {
	// make extra sure we have these fields set
	if dcdtData.Value == nil {
		dcdtData.Value = big.NewInt(0)
	}

	dcdtData.Type = uint32(core.Fungible)
}

func (e *dcdtAsBalance) saveDCDTData(accountDataHandler vmcommon.AccountDataHandler, dcdtData *dcdt.DCDigitalToken) error {
	marshaledData, err := e.marshaller.Marshal(dcdtData)
	if err != nil {
		return err
	}

	return accountDataHandler.SaveKeyValue(e.keyPrefix, marshaledData)
}

// IsInterfaceNil checks if the underlying pointer is nil
func (e *dcdtAsBalance) IsInterfaceNil() bool {
	return e == nil
}
