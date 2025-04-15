package bridge

import (
	"encoding/hex"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	coreAPI "github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk-abi/abi"
	"github.com/stretchr/testify/require"

	chainSim "github.com/TerraDharitri/drt-go-chain/integrationTests/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/dtos"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/process"
)

const (
	deposit = "deposit"
)

var serializer, _ = abi.NewSerializer(abi.ArgsNewSerializer{
	PartsSeparator: "@",
})

// ArgsBridgeSetup holds the arguments for bridge setup
type ArgsBridgeSetup struct {
	DCDTSafeAddress  []byte
	FeeMarketAddress []byte
}

// deploySovereignBridgeSetup will deploy all bridge contracts
// This function will:
// - deploy dcdt-safe contract
// - deploy fee-market contract
// - set the fee-market address inside dcdt-safe contract
// - disable fee in fee-market contract
// - unpause dcdt-safe contract so deposit operations can start
func deploySovereignBridgeSetup(
	t *testing.T,
	cs chainSim.ChainSimulator,
	wallet dtos.WalletAddress,
	dcdtSafeWasmPath string,
	feeMarketWasmPath string,
) ArgsBridgeSetup {
	nodeHandler := cs.GetNodeHandler(core.SovereignChainShardId)
	systemScAddress := chainSim.GetSysAccBytesAddress(t, nodeHandler)
	nonce := GetNonce(t, nodeHandler, wallet.Bech32)

	dcdtSafeArgs := "@000000000000000005002412c4ab184562d62a3eddaa7227730e9f17c53268a3" // pre-computed fee_market_address
	dcdtSafeAddress := chainSim.DeployContract(t, cs, wallet.Bytes, &nonce, systemScAddress, dcdtSafeArgs, dcdtSafeWasmPath)

	feeMarketArgs := "@" + hex.EncodeToString(dcdtSafeAddress) + // dcdt_safe_address
		"@00" // no fee
	feeMarketAddress := chainSim.DeployContract(t, cs, wallet.Bytes, &nonce, systemScAddress, feeMarketArgs, feeMarketWasmPath)

	chainSim.SendTransactionWithSuccess(t, cs, wallet.Bytes, &nonce, dcdtSafeAddress, chainSim.ZeroValue, "unpause", uint64(10000000))

	return ArgsBridgeSetup{
		DCDTSafeAddress:  dcdtSafeAddress,
		FeeMarketAddress: feeMarketAddress,
	}
}

// GetNonce returns account's nonce
func GetNonce(t *testing.T, nodeHandler process.NodeHandler, address string) uint64 {
	acc, _, err := nodeHandler.GetFacadeHandler().GetAccount(address, coreAPI.AccountQueryOptions{})
	require.Nil(t, err)

	return acc.Nonce
}

// Deposit will deposit tokens in the bridge sc safe contract
func Deposit(
	t *testing.T,
	cs chainSim.ChainSimulator,
	sender []byte,
	nonce *uint64,
	contract []byte,
	tokens []chainSim.ArgsDepositToken,
	receiver []byte,
	transferData *sovereign.TransferData,
) *transaction.ApiTransactionResult {
	require.True(t, len(tokens) > 0)

	args := make([]any, 0)
	args = append(args, &abi.AddressValue{Value: contract})
	args = append(args, &abi.U32Value{Value: uint32(len(tokens))})
	for _, token := range tokens {
		args = append(args, &abi.StringValue{Value: token.Identifier})
		args = append(args, &abi.U64Value{Value: token.Nonce})
		args = append(args, &abi.BigUIntValue{Value: token.Amount})
	}
	args = append(args, &abi.StringValue{Value: deposit})
	args = append(args, &abi.AddressValue{Value: receiver})
	args = append(args, &abi.OptionalValue{Value: getTransferDataValue(transferData)})

	multiTransferArg, err := serializer.Serialize(args)
	require.Nil(t, err)
	depositArgs := core.BuiltInFunctionMultiDCDTNFTTransfer +
		"@" + multiTransferArg

	return chainSim.SendTransaction(t, cs, sender, nonce, sender, chainSim.ZeroValue, depositArgs, uint64(20000000))
}

func getTransferDataValue(transferData *sovereign.TransferData) any {
	if transferData == nil {
		return nil
	}

	arguments := make([]abi.SingleValue, len(transferData.Args))
	for i, arg := range transferData.Args {
		arguments[i] = &abi.BytesValue{Value: arg}
	}
	return &abi.MultiValue{
		Items: []any{
			&abi.U64Value{Value: transferData.GasLimit},
			&abi.BytesValue{Value: transferData.Function},
			&abi.ListValue{Items: arguments},
		},
	}
}
