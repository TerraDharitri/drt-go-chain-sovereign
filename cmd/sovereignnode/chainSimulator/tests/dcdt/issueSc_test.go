package dcdt

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	sovereignChainSimulator "github.com/TerraDharitri/drt-go-chain/cmd/sovereignnode/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/config"
	chainSim "github.com/TerraDharitri/drt-go-chain/integrationTests/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/components/api"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/process"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	coreAPI "github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/stretchr/testify/require"
)

const (
	defaultPathToInitialConfig = "../../../../node/config/"
	sovereignConfigPath        = "../../../config/"
	issueWasmPath              = "../testdata/issue.wasm"
	issuePrice                 = "5000000000000000000"
)

var fungibleRoles = []string{
	core.DCDTRoleLocalMint,
	core.DCDTRoleLocalBurn,
	core.DCDTRoleTransfer,
}

// The test will deploy issue.wasm contract.
// The contract contains 3 endpoints (issue, setRoles and mint) which are called in the test
func TestSovereignChainSimulator_SmartContract_IssueToken(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	cs, err := sovereignChainSimulator.NewSovereignChainSimulator(sovereignChainSimulator.ArgsSovereignChainSimulator{
		SovereignConfigPath: sovereignConfigPath,
		ArgsChainSimulator: &chainSimulator.ArgsChainSimulator{
			BypassTxSignatureCheck: true,
			TempDir:                t.TempDir(),
			PathToInitialConfig:    defaultPathToInitialConfig,
			GenesisTimestamp:       time.Now().Unix(),
			RoundDurationInMillis:  uint64(6000),
			RoundsPerEpoch: core.OptionalUint64{
				HasValue: true,
				Value:    25,
			},
			ApiInterface:     api.NewNoApiInterface(),
			MinNodesPerShard: 2,
			AlterConfigsFunction: func(cfg *config.Configs) {
				cfg.SystemSCConfig.DCDTSystemSCConfig.BaseIssuingCost = issuePrice
			},
		},
	})
	require.Nil(t, err)
	require.NotNil(t, cs)

	defer cs.Close()

	err = cs.GenerateBlocksUntilEpochIsReached(6)
	require.Nil(t, err)

	nodeHandler := cs.GetNodeHandler(core.SovereignChainShardId)
	systemScAddress := chainSim.GetSysAccBytesAddress(t, nodeHandler)

	wallet, err := cs.GenerateAndMintWalletAddress(core.SovereignChainShardId, big.NewInt(0).Mul(chainSim.OneREWA, big.NewInt(100)))
	require.Nil(t, err)
	nonce := uint64(0)

	err = cs.GenerateBlocks(1)
	require.Nil(t, err)

	deployedContractAddress := chainSim.DeployContract(t, cs, wallet.Bytes, &nonce, systemScAddress, "", issueWasmPath)
	deployedContractAddressBech32, err := nodeHandler.GetCoreComponents().AddressPubKeyConverter().Encode(deployedContractAddress)
	require.Nil(t, err)

	issueCost, _ := big.NewInt(0).SetString(issuePrice, 10)
	chainSim.SendTransactionWithSuccess(t, cs, wallet.Bytes, &nonce, deployedContractAddress, issueCost, "issue", uint64(60000000))

	account, _, err := nodeHandler.GetFacadeHandler().GetAccount(deployedContractAddressBech32, coreAPI.AccountQueryOptions{})
	require.Nil(t, err)
	require.True(t, account.Balance == "0")

	issuedDCDTs, err := nodeHandler.GetFacadeHandler().GetAllIssuedDCDTs(core.FungibleDCDT)
	require.Nil(t, err)
	require.NotNil(t, issuedDCDTs)
	require.True(t, len(issuedDCDTs) == 1)
	tokenIdentifier := issuedDCDTs[0]

	setRolesArgs := "setRoles@" + hex.EncodeToString([]byte(tokenIdentifier))
	chainSim.SendTransactionWithSuccess(t, cs, wallet.Bytes, &nonce, deployedContractAddress, chainSim.ZeroValue, setRolesArgs, uint64(60000000))

	checkAllRoles(t, nodeHandler, deployedContractAddressBech32, tokenIdentifier, fungibleRoles)

	expectedMintedAmount, _ := big.NewInt(0).SetString("123000000000000000000", 10)
	mintTxArgs := "mint" +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(expectedMintedAmount.Bytes())
	chainSim.SendTransactionWithSuccess(t, cs, wallet.Bytes, &nonce, deployedContractAddress, chainSim.ZeroValue, mintTxArgs, uint64(20000000))

	account, _, err = nodeHandler.GetFacadeHandler().GetAccount(deployedContractAddressBech32, coreAPI.AccountQueryOptions{})
	require.Nil(t, err)
	require.True(t, account.Balance == "0")

	dcdts, _, err := nodeHandler.GetFacadeHandler().GetAllDCDTTokens(deployedContractAddressBech32, coreAPI.AccountQueryOptions{})
	require.Nil(t, err)
	require.Equal(t, 1, len(dcdts))
	require.Equal(t, expectedMintedAmount, dcdts[tokenIdentifier].Value)
}

func checkAllRoles(t *testing.T, nodeHandler process.NodeHandler, address string, tokenIdentifier string, roles []string) {
	dcdtsRoles, _, err := nodeHandler.GetFacadeHandler().GetDCDTsRoles(address, coreAPI.AccountQueryOptions{})
	require.Nil(t, err)
	require.NotNil(t, dcdtsRoles)
	require.True(t, len(dcdtsRoles[tokenIdentifier]) == len(roles))
	for i, role := range roles {
		require.Equal(t, role, dcdtsRoles[tokenIdentifier][i])
	}
}
