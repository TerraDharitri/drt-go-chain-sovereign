package dcdt

import (
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	coreAPI "github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/stretchr/testify/require"

	sovereignChainSimulator "github.com/TerraDharitri/drt-go-chain/cmd/sovereignnode/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/config"
	chainSim "github.com/TerraDharitri/drt-go-chain/integrationTests/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/components/api"
)

var nftV2Roles = []string{
	core.DCDTRoleNFTCreate,
	core.DCDTRoleNFTBurn,
	core.DCDTRoleNFTUpdateAttributes,
	core.DCDTRoleNFTAddURI,
	core.DCDTRoleNFTRecreate,
	core.DCDTRoleModifyCreator,
	core.DCDTRoleModifyRoyalties,
	core.DCDTRoleSetNewURI,
	core.DCDTRoleNFTUpdate,
}

func TestSovereignChainSimulator_RegisterNftWithPrefix(t *testing.T) {
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

	nodeHandler := cs.GetNodeHandler(core.SovereignChainShardId)

	wallet, err := cs.GenerateAndMintWalletAddress(core.SovereignChainShardId, chainSim.InitialAmount)
	require.Nil(t, err)
	nonce := uint64(0)

	err = cs.GenerateBlocks(1)
	require.Nil(t, err)

	issueCost, _ := big.NewInt(0).SetString(issuePrice, 10)
	nftName := "NFTNAME"
	nftTicker := "NFTTICKER"
	nftIdentifier := chainSim.RegisterAndSetAllRoles(t, cs, wallet.Bytes, &nonce, issueCost, nftName, nftTicker, core.NonFungibleDCDTv2, 0)

	checkAllRoles(t, nodeHandler, wallet.Bech32, nftIdentifier, nftV2Roles)

	initialSupply := big.NewInt(1)
	createArgs := createNftArgs(nftIdentifier, initialSupply, "NFTNAME #1")
	chainSim.SendTransactionWithSuccess(t, cs, wallet.Bytes, &nonce, wallet.Bytes, chainSim.ZeroValue, createArgs, uint64(60000000))

	tokens, _, err := nodeHandler.GetFacadeHandler().GetAllDCDTTokens(wallet.Bech32, coreAPI.AccountQueryOptions{})
	require.Nil(t, err)
	require.NotNil(t, tokens)
	require.True(t, len(tokens) == 2)
	require.Equal(t, initialSupply, tokens[nftIdentifier+"-01"].Value)
}
