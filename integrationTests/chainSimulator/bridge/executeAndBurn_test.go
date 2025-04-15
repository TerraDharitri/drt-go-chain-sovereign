package bridge

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	dataApi "github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/config"
	chainSim "github.com/TerraDharitri/drt-go-chain/integrationTests/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/components/api"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/process"
)

const (
	defaultPathToInitialConfig = "../../../cmd/node/config/"
)

var sovChainPrefix = "sov"

func TestChainSimulator_ExecuteWithMintAndBurnFungibleWithDeposit(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	token := sovChainPrefix + "-SOVTKN-1a2b3c"
	tokenNonce := uint64(0)

	bridgedInTokens := make([]chainSim.ArgsDepositToken, 0)
	bridgedInTokens = append(bridgedInTokens, chainSim.ArgsDepositToken{
		Identifier: token,
		Nonce:      tokenNonce,
		Amount:     big.NewInt(123),
		Type:       core.Fungible,
	})
	bridgedInTokens = append(bridgedInTokens, chainSim.ArgsDepositToken{
		Identifier: token,
		Nonce:      tokenNonce,
		Amount:     big.NewInt(100),
		Type:       core.Fungible,
	})

	bridgedOutTokens := make([]chainSim.ArgsDepositToken, 0)
	bridgedOutTokens = append(bridgedOutTokens, chainSim.ArgsDepositToken{
		Identifier: token,
		Nonce:      tokenNonce,
		Amount:     big.NewInt(12),
		Type:       core.Fungible,
	})
	bridgedOutTokens = append(bridgedOutTokens, chainSim.ArgsDepositToken{
		Identifier: token,
		Nonce:      tokenNonce,
		Amount:     big.NewInt(10),
		Type:       core.Fungible,
	})

	simulateExecutionAndDeposit(t, bridgedInTokens, bridgedOutTokens)
}

func TestChainSimulator_ExecuteWithMintMultipleDcdtsAndBurnNftWithDeposit(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	nftV2 := sovChainPrefix + "-NFTV2-a1b2c3"
	nftV2Nonce := uint64(10)
	token := sovChainPrefix + "-TKN-1d2e3a"

	// TODO DRT-15942 add dynamic NFT type for bridge transfer
	bridgedInTokens := make([]chainSim.ArgsDepositToken, 0)
	bridgedInTokens = append(bridgedInTokens, chainSim.ArgsDepositToken{
		Identifier: nftV2,
		Nonce:      nftV2Nonce,
		Amount:     big.NewInt(1),
		Type:       core.NonFungibleV2,
	})
	bridgedInTokens = append(bridgedInTokens, chainSim.ArgsDepositToken{
		Identifier: token,
		Nonce:      0,
		Amount:     big.NewInt(1),
		Type:       core.Fungible,
	})

	bridgedOutTokens := make([]chainSim.ArgsDepositToken, 0)
	bridgedOutTokens = append(bridgedOutTokens, chainSim.ArgsDepositToken{
		Identifier: nftV2,
		Nonce:      nftV2Nonce,
		Amount:     big.NewInt(1),
		Type:       core.NonFungibleV2,
	})

	simulateExecutionAndDeposit(t, bridgedInTokens, bridgedOutTokens)
}

func TestChainSimulator_ExecuteWithMintAndBurnSftWithDeposit(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	sft := sovChainPrefix + "-SOVSFT-654321"
	sftNonce := uint64(123)

	// TODO DRT-15942 add dynamic SFT type for bridge transfer
	bridgedInTokens := make([]chainSim.ArgsDepositToken, 0)
	bridgedInTokens = append(bridgedInTokens, chainSim.ArgsDepositToken{
		Identifier: sft,
		Nonce:      sftNonce,
		Amount:     big.NewInt(50),
		Type:       core.SemiFungible,
	})

	bridgedOutTokens := make([]chainSim.ArgsDepositToken, 0)
	bridgedOutTokens = append(bridgedOutTokens, chainSim.ArgsDepositToken{
		Identifier: sft,
		Nonce:      sftNonce,
		Amount:     big.NewInt(20),
		Type:       core.SemiFungible,
	})

	simulateExecutionAndDeposit(t, bridgedInTokens, bridgedOutTokens)
}

func simulateExecutionAndDeposit(
	t *testing.T,
	bridgedInTokens []chainSim.ArgsDepositToken,
	bridgedOutTokens []chainSim.ArgsDepositToken,
) {
	roundsPerEpoch := core.OptionalUint64{
		HasValue: true,
		Value:    20,
	}

	whiteListedAddress := "drt1qqqqqqqqqqqqqpgqmzzm05jeav6d5qvna0q2pmcllelkz8xddz3sew8p92"
	cs, err := chainSimulator.NewChainSimulator(chainSimulator.ArgsChainSimulator{
		BypassTxSignatureCheck:   true,
		TempDir:                  t.TempDir(),
		PathToInitialConfig:      defaultPathToInitialConfig,
		NumOfShards:              1,
		GenesisTimestamp:         time.Now().Unix(),
		RoundDurationInMillis:    uint64(6000),
		RoundsPerEpoch:           roundsPerEpoch,
		ApiInterface:             api.NewNoApiInterface(),
		MinNodesPerShard:         3,
		MetaChainMinNodes:        3,
		NumNodesWaitingListMeta:  0,
		NumNodesWaitingListShard: 0,
		AlterConfigsFunction: func(cfg *config.Configs) {
			cfg.GeneralConfig.VirtualMachine.Execution.TransferAndExecuteByUserAddresses = []string{whiteListedAddress}
			cfg.EpochConfig.EnableEpochs.DynamicDCDTEnableEpoch = 0
		},
	})
	require.Nil(t, err)
	require.NotNil(t, cs)

	defer cs.Close()

	err = cs.GenerateBlocksUntilEpochIsReached(4)
	require.Nil(t, err)

	nodeHandler := cs.GetNodeHandler(0)

	// Deploy bridge setup
	initialAddress := "drt1l6xt0rqlyzw56a3k8xwwshq2dcjwy3q9cppucvqsmdyw8r98dz3sq9c49p"
	argsDcdtSafe := ArgsDcdtSafe{
		ChainPrefix:       sovChainPrefix,
		IssuePaymentToken: "WREWA-bd4d79",
	}
	initOwnerAndSysAccState(t, cs, initialAddress, argsDcdtSafe)
	bridgeData := deployBridgeSetup(t, cs, initialAddress, argsDcdtSafe, enshrineDcdtSafeContract, enshrineDcdtSafeWasmPath)
	chainSim.RequireAccountHasToken(t, cs, argsDcdtSafe.IssuePaymentToken, initialAddress, big.NewInt(0))

	dcdtSafeEncoded, _ := nodeHandler.GetCoreComponents().AddressPubKeyConverter().Encode(bridgeData.DCDTSafeAddress)
	require.Equal(t, whiteListedAddress, dcdtSafeEncoded)

	wallet, err := cs.GenerateAndMintWalletAddress(0, chainSim.InitialAmount)
	require.Nil(t, err)
	nonce := uint64(0)
	paymentTokenAmount, _ := big.NewInt(0).SetString("1000000000000000000", 10)
	chainSim.SetDcdtInWallet(t, cs, wallet, argsDcdtSafe.IssuePaymentToken, 0, dcdt.DCDigitalToken{Value: paymentTokenAmount})

	// We need to register tokens originated from sovereign (to pay the issue cost)
	// Only the tokens with sovereign prefix need to be registered (these are the ones that will be minted), the rest will be taken from contract balance
	tokens := getUniquePrefixedTokens(bridgedInTokens, argsDcdtSafe.ChainPrefix)
	registerSovereignNewTokens(t, cs, wallet, &nonce, bridgeData.DCDTSafeAddress, argsDcdtSafe.IssuePaymentToken, tokens)

	// We will deposit an array of prefixed tokens from a sovereign chain to the main chain,
	// expecting these tokens to be minted by the whitelisted DCDT safe sc and transferred to our wallet address.
	txResult := executeOperation(t, cs, bridgeData.OwnerAccount.Wallet, wallet.Bytes, &bridgeData.OwnerAccount.Nonce, bridgeData.DCDTSafeAddress, bridgedInTokens, wallet.Bytes, nil)
	chainSim.RequireSuccessfulTransaction(t, txResult)
	for _, bridgedInToken := range groupTokens(bridgedInTokens) {
		chainSim.RequireAccountHasToken(t, cs, getTokenIdentifier(bridgedInToken), wallet.Bech32, bridgedInToken.Amount)
		checkMetaDataInAccounts(t, cs, bridgedInToken, wallet.Bech32, bridgedInToken.Amount)
	}

	// deposit an array of tokens from main chain to sovereign chain,
	// expecting these tokens to be burned by the whitelisted DCDT safe sc
	txResult = deposit(t, cs, wallet.Bytes, &nonce, bridgeData.DCDTSafeAddress, bridgedOutTokens, wallet.Bytes)
	chainSim.RequireSuccessfulTransaction(t, txResult)

	bridgedTokens := groupTokens(bridgedInTokens)
	for _, bridgedOutToken := range groupTokens(bridgedOutTokens) {
		bridgedValue, err := getBridgedValue(bridgedTokens, bridgedOutToken.Identifier)
		require.Nil(t, err)

		fullTokenIdentifier := getTokenIdentifier(bridgedOutToken)
		remainingAmount := big.NewInt(0).Sub(bridgedValue, bridgedOutToken.Amount)
		chainSim.RequireAccountHasToken(t, cs, fullTokenIdentifier, wallet.Bech32, remainingAmount)
		chainSim.RequireAccountHasToken(t, cs, fullTokenIdentifier, dcdtSafeEncoded, big.NewInt(0))
		checkMetaDataInAccounts(t, cs, bridgedOutToken, wallet.Bech32, remainingAmount)

		tokenSupply, err := nodeHandler.GetFacadeHandler().GetTokenSupply(fullTokenIdentifier)
		require.Nil(t, err)
		require.NotNil(t, tokenSupply)
		require.Equal(t, bridgedOutToken.Amount.String(), tokenSupply.Burned)
	}
}

func getUniquePrefixedTokens(bridgedTokens []chainSim.ArgsDepositToken, prefix string) []string {
	tokens := make([]string, 0)
	seen := make(map[string]bool)
	for _, token := range bridgedTokens {
		if strings.HasPrefix(token.Identifier, prefix+"-") && !seen[token.Identifier] {
			tokens = append(tokens, token.Identifier)
			seen[token.Identifier] = true
		}
	}
	return tokens
}

func getBridgedValue(bridgeInTokens []chainSim.ArgsDepositToken, token string) (*big.Int, error) {
	for _, tkn := range bridgeInTokens {
		if tkn.Identifier == token {
			return tkn.Amount, nil
		}
	}
	return nil, fmt.Errorf("token not found")
}

func checkMetaDataInAccounts(
	t *testing.T,
	cs chainSim.ChainSimulator,
	token chainSim.ArgsDepositToken,
	account string,
	expectedAmount *big.Int,
) {
	addressShardID := chainSim.GetShardForAddress(cs, account)
	nodeHandler := cs.GetNodeHandler(addressShardID)

	// get user account token data
	dcdtValue := getAccountTokenData(t, nodeHandler, account, token)

	if expectedAmount.Cmp(big.NewInt(0)) == 0 {
		require.Empty(t, dcdtValue)            // expect that key doesn't exist in account
		if metaDataOnUserAccount(token.Type) { // for other token types the key can exist because other wallets in shard can have the token
			requireNoTokenDataInSysAccount(t, nodeHandler, token) // no keys in system account
		}
		return
	}

	dcdtData := &dcdt.DCDigitalToken{}
	err := nodeHandler.GetCoreComponents().InternalMarshalizer().Unmarshal(dcdtData, dcdtValue)
	require.Nil(t, err)
	require.NotNil(t, dcdtData)
	require.Equal(t, expectedAmount, dcdtData.Value)
	require.Equal(t, uint32(token.Type), dcdtData.Type)

	if token.Type == core.Fungible {
		require.Nil(t, dcdtData.TokenMetaData)
		requireNoTokenDataInSysAccount(t, nodeHandler, token) // no keys in system account
	} else if token.Type == core.NonFungibleV2 || token.Type == core.DynamicNFT {
		require.NotNil(t, dcdtData.TokenMetaData)
		require.Equal(t, token.Nonce, dcdtData.TokenMetaData.Nonce)
		requireNoTokenDataInSysAccount(t, nodeHandler, token) // no keys in system account
	} else {
		require.Nil(t, dcdtData.TokenMetaData)

		// get system account token data
		dcdtValue = getAccountTokenData(t, nodeHandler, chainSim.DCDTSystemAccount, token)

		dcdtData = &dcdt.DCDigitalToken{}
		err = nodeHandler.GetCoreComponents().InternalMarshalizer().Unmarshal(dcdtData, dcdtValue)
		require.Nil(t, err)
		require.NotNil(t, dcdtData)
		require.GreaterOrEqual(t, dcdtData.Value.Uint64(), expectedAmount.Uint64()) // greater if other wallets have the token
		require.Equal(t, uint32(token.Type), dcdtData.Type)
		require.NotNil(t, dcdtData.TokenMetaData)
		require.Equal(t, token.Nonce, dcdtData.TokenMetaData.Nonce)
	}
}

func getAccountTokenData(
	t *testing.T,
	nodeHandler process.NodeHandler,
	account string,
	token chainSim.ArgsDepositToken,
) []byte {
	accountKeys, _, err := nodeHandler.GetFacadeHandler().GetKeyValuePairs(account, dataApi.AccountQueryOptions{})
	require.Nil(t, err)
	require.NotNil(t, accountKeys)

	dcdtValue, err := hex.DecodeString(accountKeys[getTokenKey(token)])
	require.Nil(t, err)

	return dcdtValue
}

func requireNoTokenDataInSysAccount(
	t *testing.T,
	nodeHandler process.NodeHandler,
	token chainSim.ArgsDepositToken,
) {
	accountKeys, _, err := nodeHandler.GetFacadeHandler().GetKeyValuePairs(chainSim.DCDTSystemAccount, dataApi.AccountQueryOptions{})
	require.Nil(t, err)
	require.NotNil(t, accountKeys)
	require.Empty(t, accountKeys[getTokenKey(token)])
}

func metaDataOnUserAccount(dcdtType core.DCDTType) bool {
	return dcdtType == core.Fungible ||
		dcdtType == core.NonFungibleV2 ||
		dcdtType == core.DynamicNFT
}

func getTokenKey(token chainSim.ArgsDepositToken) string {
	nonce := ""
	if token.Nonce > 0 {
		nonce = hex.EncodeToString(big.NewInt(0).SetUint64(token.Nonce).Bytes())
	}
	return hex.EncodeToString([]byte(core.ProtectedKeyPrefix+core.DCDTKeyIdentifier)) +
		hex.EncodeToString([]byte(token.Identifier)) +
		nonce
}
