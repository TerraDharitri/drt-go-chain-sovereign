package chainSimulator

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	dataApi "github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/wasm"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/configs"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/dtos"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/process"
	"github.com/TerraDharitri/drt-go-chain/vm"
)

const (
	vmTypeHex                               = "0500"
	codeMetadata                            = "0500"
	minGasPrice                             = 1000000000
	txVersion                               = 1
	mockTxSignature                         = "sig"
	maxNumOfBlocksToGenerateWhenExecutingTx = 10
	signalError                             = "signalError"
	internalVMError                         = "internalVMErrors"

	// OkReturnCode the const for the ok return code
	OkReturnCode = "ok"
	// DCDTSystemAccount the bech32 address for dcdt system account
	DCDTSystemAccount = "drt1llllllllllllllllllllllllllllllllllllllllllllllllllls9258a4"
)

var (
	// ZeroValue the variable for the zero big int
	ZeroValue = big.NewInt(0)
	// OneREWA the variable for one rewa value
	OneREWA = big.NewInt(1000000000000000000)
	// MinimumStakeValue the variable for the minimum stake value
	MinimumStakeValue = big.NewInt(0).Mul(OneREWA, big.NewInt(2500))
	// InitialAmount the variable for initial minting amount in account
	InitialAmount = big.NewInt(0).Mul(OneREWA, big.NewInt(100))
)

// ArgsDepositToken holds the arguments for a token
type ArgsDepositToken struct {
	Identifier string
	Nonce      uint64
	Amount     *big.Int
	Type       core.DCDTType
}

// Account holds the arguments for a user account
type Account struct {
	Wallet dtos.WalletAddress
	Nonce  uint64
}

// GetSysAccBytesAddress will return the system account bytes address
func GetSysAccBytesAddress(t *testing.T, nodeHandler process.NodeHandler) []byte {
	addressBytes, err := nodeHandler.GetCoreComponents().AddressPubKeyConverter().Decode("drt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq85hk5z")
	require.Nil(t, err)

	return addressBytes
}

// GetSysContactDeployAddressBytes will return the system contract deploy address
func GetSysContactDeployAddressBytes(t *testing.T, nodeHandler process.NodeHandler) []byte {
	addressBytes, err := nodeHandler.GetCoreComponents().AddressPubKeyConverter().Decode("drt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq85hk5z")
	require.Nil(t, err)

	return addressBytes
}

// GetShardForAddress will return the shard of the address
func GetShardForAddress(cs ChainSimulator, address string) uint32 {
	nodeHandler := cs.GetNodeHandler(0)
	pubKey, _ := nodeHandler.GetCoreComponents().AddressPubKeyConverter().Decode(address)
	return nodeHandler.GetShardCoordinator().ComputeId(pubKey)
}

// DeployContract will deploy a smart contract and return its address
func DeployContract(
	t *testing.T,
	cs ChainSimulator,
	sender []byte,
	nonce *uint64,
	receiver []byte,
	data string,
	wasmPath string,
) []byte {
	data = wasm.GetSCCode(wasmPath) + "@" + vmTypeHex + "@" + codeMetadata + data

	tx := GenerateTransaction(sender, *nonce, receiver, ZeroValue, data, uint64(200000000))
	txResult, err := cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlocksToGenerateWhenExecutingTx)
	*nonce++

	require.Nil(t, err)
	RequireSuccessfulTransaction(t, txResult)

	address := txResult.Logs.Events[0].Topics[0]
	require.NotNil(t, address)
	return address
}

// GenerateTransaction will generate a transaction object
func GenerateTransaction(sender []byte, nonce uint64, receiver []byte, value *big.Int, data string, gasLimit uint64) *transaction.Transaction {
	return &transaction.Transaction{
		Nonce:     nonce,
		Value:     value,
		SndAddr:   sender,
		RcvAddr:   receiver,
		Data:      []byte(data),
		GasLimit:  gasLimit,
		GasPrice:  minGasPrice,
		ChainID:   []byte(configs.ChainID),
		Version:   txVersion,
		Signature: []byte(mockTxSignature),
	}
}

// SendTransactionWithSuccess will send a transaction and expect successful execution and return the result
func SendTransactionWithSuccess(
	t *testing.T,
	cs ChainSimulator,
	sender []byte,
	nonce *uint64,
	receiver []byte,
	value *big.Int,
	data string,
	gasLimit uint64,
) *transaction.ApiTransactionResult {
	txResult := SendTransaction(t, cs, sender, nonce, receiver, value, data, gasLimit)
	RequireSuccessfulTransaction(t, txResult)
	return txResult
}

// SendTransaction will send a transaction and return the result
func SendTransaction(
	t *testing.T,
	cs ChainSimulator,
	sender []byte,
	nonce *uint64,
	receiver []byte,
	value *big.Int,
	data string,
	gasLimit uint64,
) *transaction.ApiTransactionResult {
	tx := GenerateTransaction(sender, *nonce, receiver, value, data, gasLimit)
	txResult, err := cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlocksToGenerateWhenExecutingTx)
	*nonce++
	require.Nil(t, err)

	return txResult
}

// RequireSuccessfulTransaction require that the transaction doesn't have signal error event
func RequireSuccessfulTransaction(t *testing.T, txResult *transaction.ApiTransactionResult) {
	require.NotNil(t, txResult)
	event := getEvent(txResult.Logs, signalError)
	if event != nil {
		require.Fail(t, string(event.Topics[1]))
	}
	require.Equal(t, transaction.TxStatusSuccess, txResult.Status)
}

// RequireSignalError require that the transaction has specific signal error
func RequireSignalError(t *testing.T, txResult *transaction.ApiTransactionResult, error string) {
	require.NotNil(t, txResult)
	event := getEvent(txResult.Logs, signalError)
	if event == nil {
		require.Fail(t, "%s event not found", signalError)
		return
	}
	require.Equal(t, error, string(event.Topics[1]))
	require.Equal(t, transaction.TxStatusSuccess, txResult.Status)
}

// RequireInternalVMError require that the transaction has specific invernal vm error
func RequireInternalVMError(t *testing.T, txResult *transaction.ApiTransactionResult, error string) {
	require.NotNil(t, txResult)
	event := getEvent(txResult.Logs, internalVMError)
	if event == nil {
		require.Fail(t, "%s event not found", internalVMError)
		return
	}
	require.Contains(t, string(event.Data), error)
	require.Equal(t, transaction.TxStatusSuccess, txResult.Status)
}

func getEvent(logs *transaction.ApiLogs, eventID string) *transaction.Events {
	if logs == nil || len(logs.Events) == 0 {
		return nil
	}

	for _, event := range logs.Events {
		if event.Identifier == eventID {
			return event
		}
	}
	return nil
}

// RequireAccountHasToken checks if the account has the amount of tokens (can also be zero)
func RequireAccountHasToken(
	t *testing.T,
	cs ChainSimulator,
	token string,
	address string,
	value *big.Int,
) {
	addressShardID := GetShardForAddress(cs, address)
	tokens, _, err := cs.GetNodeHandler(addressShardID).GetFacadeHandler().GetAllDCDTTokens(address, dataApi.AccountQueryOptions{})
	require.Nil(t, err)

	tokenData, found := tokens[token]

	if value.Cmp(big.NewInt(0)) == 0 {
		require.False(t, found)
		return
	}
	require.True(t, found, fmt.Sprintf("%s token not found", token))
	require.Equal(t, tokenData.Value, value)
}

// TransferDCDT will transfer the amount of dcdt token to an address
func TransferDCDT(
	t *testing.T,
	cs ChainSimulator,
	sender, receiver []byte,
	nonce *uint64,
	token string,
	amount *big.Int,
	args ...[]byte,
) {
	dcdtTransferArgs := core.BuiltInFunctionDCDTTransfer +
		"@" + hex.EncodeToString([]byte(token)) +
		"@" + hex.EncodeToString(amount.Bytes())
	for _, arg := range args {
		dcdtTransferArgs = dcdtTransferArgs +
			"@" + hex.EncodeToString(arg)
	}
	txResult := SendTransaction(t, cs, sender, nonce, receiver, ZeroValue, dcdtTransferArgs, uint64(5000000))
	RequireSuccessfulTransaction(t, txResult)
}

// TransferDCDTNFT will transfer the amount of NFT/SFT token to an address
func TransferDCDTNFT(
	t *testing.T,
	cs ChainSimulator,
	sender, receiver []byte,
	nonce *uint64,
	token string,
	tokenNonce uint64,
	amount *big.Int,
	args ...[]byte,
) {
	dcdtNftTransferArgs :=
		core.BuiltInFunctionDCDTNFTTransfer +
			"@" + hex.EncodeToString([]byte(token)) +
			"@" + hex.EncodeToString(big.NewInt(int64(tokenNonce)).Bytes()) +
			"@" + hex.EncodeToString(amount.Bytes()) +
			"@" + hex.EncodeToString(receiver)
	for _, arg := range args {
		dcdtNftTransferArgs = dcdtNftTransferArgs +
			"@" + hex.EncodeToString(arg)
	}
	txResult := SendTransaction(t, cs, sender, nonce, sender, ZeroValue, dcdtNftTransferArgs, uint64(5000000))
	RequireSuccessfulTransaction(t, txResult)
}

// IssueFungible will issue a fungible token
func IssueFungible(
	t *testing.T,
	cs ChainSimulator,
	sender []byte,
	nonce *uint64,
	issueCost *big.Int,
	tokenName string,
	tokenTicker string,
	numDecimals int,
	supply *big.Int,
) string {
	issueArgs := "issue" +
		"@" + hex.EncodeToString([]byte(tokenName)) +
		"@" + hex.EncodeToString([]byte(tokenTicker)) +
		"@" + hex.EncodeToString(supply.Bytes()) +
		"@" + fmt.Sprintf("%X", numDecimals) +
		"@" + hex.EncodeToString([]byte("canAddSpecialRoles")) +
		"@" + hex.EncodeToString([]byte("true"))
	txResult := SendTransaction(t, cs, sender, nonce, vm.DCDTSCAddress, issueCost, issueArgs, uint64(60000000))
	RequireSuccessfulTransaction(t, txResult)

	return GetIssuedDcdtIdentifier(t, cs, tokenTicker, core.FungibleDCDT)
}

// GetIssuedDcdtIdentifier will return the token identifier for and issued token
func GetIssuedDcdtIdentifier(t *testing.T, cs ChainSimulator, ticker string, tokenType string) string {
	dcdtScAddressShardId := cs.GetNodeHandler(0).GetShardCoordinator().ComputeId(vm.DCDTSCAddress)
	issuedTokens, err := cs.GetNodeHandler(dcdtScAddressShardId).GetFacadeHandler().GetAllIssuedDCDTs(tokenType)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(issuedTokens), 1, "no issued tokens found of type %s", tokenType)

	for _, issuedToken := range issuedTokens {
		if strings.Contains(issuedToken, ticker) {
			return issuedToken
		}
	}

	require.Fail(t, "could not find the issued token")
	return ""
}

// InitAddressesAndSysAccState will initialize system account state and other addresses if provided
func InitAddressesAndSysAccState(
	t *testing.T,
	cs ChainSimulator,
	initialAddresses ...string,
) {
	addressesState := []*dtos.AddressState{
		{
			Address: DCDTSystemAccount,
		},
	}
	for _, address := range initialAddresses {
		addressesState = append(addressesState,
			&dtos.AddressState{
				Address: address,
				Balance: "10000000000000000000000",
			},
		)
	}
	err := cs.SetStateMultiple(addressesState)
	require.Nil(t, err)

	err = cs.GenerateBlocks(1)
	require.Nil(t, err)
}

// SetDcdtInWallet will add token key in wallet storage without adding key in system account
func SetDcdtInWallet(
	t *testing.T,
	cs ChainSimulator,
	wallet dtos.WalletAddress,
	token string,
	tokenNonce uint64,
	tokenData dcdt.DCDigitalToken,
) {
	marshalledTokenData, err := cs.GetNodeHandler(0).GetCoreComponents().InternalMarshalizer().Marshal(&tokenData)
	require.NoError(t, err)

	nonce := ""
	if tokenNonce != 0 {
		nonce = hex.EncodeToString(big.NewInt(0).SetUint64(tokenNonce).Bytes())
	}
	tokenKey := hex.EncodeToString([]byte(core.ProtectedKeyPrefix+core.DCDTKeyIdentifier+token)) + nonce
	tokenValue := hex.EncodeToString(marshalledTokenData)
	keyValueMap := map[string]string{
		tokenKey: tokenValue,
	}
	err = cs.SetKeyValueForAddress(wallet.Bech32, keyValueMap)
	require.NoError(t, err)

	err = cs.GenerateBlocks(1)
	require.Nil(t, err)
}

// IssueSemiFungible will issue a semi fungible token
func IssueSemiFungible(
	t *testing.T,
	cs ChainSimulator,
	sender []byte,
	nonce *uint64,
	issueCost *big.Int,
	sftName string,
	sftTicker string,
) string {
	issueArgs := "issueSemiFungible" +
		"@" + hex.EncodeToString([]byte(sftName)) +
		"@" + hex.EncodeToString([]byte(sftTicker))
	SendTransaction(t, cs, sender, nonce, vm.DCDTSCAddress, issueCost, issueArgs, uint64(60000000))

	return GetIssuedDcdtIdentifier(t, cs, sftTicker, core.SemiFungibleDCDT)
}

// RegisterAndSetAllRoles will issue an dcdt collection with all roles enabled
func RegisterAndSetAllRoles(
	t *testing.T,
	cs ChainSimulator,
	sender []byte,
	nonce *uint64,
	issueCost *big.Int,
	dcdtName string,
	dcdtTicker string,
	tokenType string,
	numDecimals int,
) string {
	dcdtType := getTokenRegisterType(tokenType)
	registerArgs := "registerAndSetAllRoles" +
		"@" + hex.EncodeToString([]byte(dcdtName)) +
		"@" + hex.EncodeToString([]byte(dcdtTicker)) +
		"@" + hex.EncodeToString([]byte(dcdtType)) +
		"@" + fmt.Sprintf("%02X", numDecimals)
	SendTransaction(t, cs, sender, nonce, vm.DCDTSCAddress, issueCost, registerArgs, uint64(60000000))

	return GetIssuedDcdtIdentifier(t, cs, dcdtTicker, tokenType)
}

// RegisterAndSetAllRolesDynamic will issue a dynamic dcdt collection with all roles enabled
func RegisterAndSetAllRolesDynamic(
	t *testing.T,
	cs ChainSimulator,
	sender []byte,
	nonce *uint64,
	issueCost *big.Int,
	dcdtName string,
	dcdtTicker string,
	tokenType string,
	numDecimals int,
) string {
	dcdtType := getTokenRegisterType(tokenType)
	registerArgs := "registerAndSetAllRolesDynamic" +
		"@" + hex.EncodeToString([]byte(dcdtName)) +
		"@" + hex.EncodeToString([]byte(dcdtTicker)) +
		"@" + hex.EncodeToString([]byte(dcdtType))
	if numDecimals != 0 {
		registerArgs += "@" + fmt.Sprintf("%02X", numDecimals)
	}
	SendTransaction(t, cs, sender, nonce, vm.DCDTSCAddress, issueCost, registerArgs, uint64(60000000))

	return GetIssuedDcdtIdentifier(t, cs, dcdtTicker, tokenType)
}

func getTokenRegisterType(tokenType string) string {
	switch tokenType {
	case core.FungibleDCDT:
		return "FNG"
	case core.NonFungibleDCDT, core.NonFungibleDCDTv2, core.DynamicNFTDCDT:
		return "NFT"
	case core.SemiFungibleDCDT, core.DynamicSFTDCDT:
		return "SFT"
	case core.MetaDCDT, core.DynamicMetaDCDT:
		return "META"
	}
	return ""
}
