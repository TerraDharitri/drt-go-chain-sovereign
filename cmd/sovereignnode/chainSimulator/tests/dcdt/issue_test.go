package dcdt

import (
	"math/big"
	"testing"
	"time"

	sovereignChainSimulator "github.com/TerraDharitri/drt-go-chain/cmd/sovereignnode/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/config"
	chainSim "github.com/TerraDharitri/drt-go-chain/integrationTests/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/components/api"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/vm"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	coreAPI "github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/stretchr/testify/require"
)

const (
	getTokenPropertiesFunc = "getTokenProperties"
)

func TestSovereignChainSimulator_IssueFungible(t *testing.T) {
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
			RoundsPerEpoch:         core.OptionalUint64{},
			ApiInterface:           api.NewNoApiInterface(),
			MinNodesPerShard:       2,
			AlterConfigsFunction: func(cfg *config.Configs) {
				cfg.SystemSCConfig.DCDTSystemSCConfig.BaseIssuingCost = issuePrice
			},
		},
	})
	require.Nil(t, err)
	require.NotNil(t, cs)

	defer cs.Close()

	time.Sleep(time.Second) // wait for VM to be ready for processing queries

	nodeHandler := cs.GetNodeHandler(core.SovereignChainShardId)

	err = cs.GenerateBlocks(1)
	require.Nil(t, err)

	issuedDCDTs, err := nodeHandler.GetFacadeHandler().GetAllIssuedDCDTs("")
	require.Nil(t, err)
	require.NotNil(t, issuedDCDTs)
	require.Equal(t, 0, len(issuedDCDTs))

	nonce := uint64(0)
	wallet, err := cs.GenerateAndMintWalletAddress(core.SovereignChainShardId, chainSim.InitialAmount)
	require.Nil(t, err)

	err = cs.GenerateBlocks(1)
	require.Nil(t, err)

	issueCost, _ := big.NewInt(0).SetString(issuePrice, 10)
	supply, _ := big.NewInt(0).SetString("123000000000000000000", 10)
	tokenName := "SovereignTkn"
	tokenTicker := "SVN"
	numDecimals := 18
	tokenIdentifier := chainSim.IssueFungible(t, cs, wallet.Bytes, &nonce, issueCost, tokenName, tokenTicker, numDecimals, supply)

	issuedDCDTs, err = nodeHandler.GetFacadeHandler().GetAllIssuedDCDTs("")
	require.Nil(t, err)
	require.NotNil(t, issuedDCDTs)
	require.Equal(t, 1, len(issuedDCDTs))

	tokens, _, err := nodeHandler.GetFacadeHandler().GetAllDCDTTokens(wallet.Bech32, coreAPI.AccountQueryOptions{})
	require.Nil(t, err)
	require.NotNil(t, tokens)
	require.True(t, len(tokens) == 2)
	require.Equal(t, supply, tokens[tokenIdentifier].Value)

	res, _, err := nodeHandler.GetFacadeHandler().ExecuteSCQuery(&process.SCQuery{
		ScAddress: vm.DCDTSCAddress,
		FuncName:  getTokenPropertiesFunc,
		Arguments: [][]byte{[]byte(tokenIdentifier)},
	})
	require.Nil(t, err)
	require.Equal(t, chainSim.OkReturnCode, res.ReturnCode)
	dcdtSupply, _ := big.NewInt(0).SetString(string(res.ReturnData[3]), 10)
	require.Equal(t, supply, dcdtSupply)

	setRolesArgs := setSpecialRole(tokenIdentifier, wallet.Bytes, fungibleRoles)
	chainSim.SendTransactionWithSuccess(t, cs, wallet.Bytes, &nonce, vm.DCDTSCAddress, chainSim.ZeroValue, setRolesArgs, uint64(60000000))

	checkAllRoles(t, nodeHandler, wallet.Bech32, tokenIdentifier, fungibleRoles)
}
