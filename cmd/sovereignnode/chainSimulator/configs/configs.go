package configs

import (
	"encoding/json"
	"math/big"
	"os"

	"github.com/TerraDharitri/drt-go-chain-core/core"

	"github.com/TerraDharitri/drt-go-chain/common/factory"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/genesis/data"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/chainSimulator"
	chainSimulatorConfigs "github.com/TerraDharitri/drt-go-chain/node/chainSimulator/configs"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/dtos"
)

var initialStakedRewaPerNode = big.NewInt(0).Mul(chainSimulator.OneREWA, big.NewInt(2500))
var initialRewaPerNode = big.NewInt(0).Mul(chainSimulator.OneREWA, big.NewInt(5000))
var initialSupply = big.NewInt(0).Mul(chainSimulator.OneREWA, big.NewInt(20000000)) // 20 million REWA

// GenerateSovereignGenesisFile will generate sovereign initial wallet keys
func GenerateSovereignGenesisFile(args chainSimulatorConfigs.ArgsChainSimulatorConfigs, configs *config.Configs) (*dtos.InitialWalletKeys, error) {
	addressConverter, err := factory.NewPubkeyConverter(configs.GeneralConfig.AddressPubkeyConverter)
	if err != nil {
		return nil, err
	}

	initialWalletKeys := &dtos.InitialWalletKeys{
		BalanceWallets: make(map[uint32]*dtos.WalletKey),
		StakeWallets:   make([]*dtos.WalletKey, 0),
	}
	addresses := make([]data.InitialAccount, 0)
	numOfNodes := int(args.NumNodesWaitingListShard + args.MinNodesPerShard)

	totalStakedValue := big.NewInt(0).Set(initialStakedRewaPerNode)
	totalStakedValue.Mul(totalStakedValue, big.NewInt(int64(numOfNodes)))

	initialBalance := big.NewInt(0).Set(initialSupply)
	initialBalance.Sub(initialBalance, totalStakedValue)

	for i := 0; i < numOfNodes; i++ {
		initialRewa := big.NewInt(0).Set(initialRewaPerNode)

		if i == numOfNodes-1 {
			remainingSupply := big.NewInt(0).Set(initialBalance)
			allBalances := big.NewInt(0).Set(initialRewa)
			allBalances.Mul(allBalances, big.NewInt(int64(numOfNodes)))
			remainingSupply.Sub(remainingSupply, allBalances)
			if remainingSupply.Cmp(big.NewInt(0)) > 0 {
				initialRewa.Add(initialRewa, remainingSupply)
			}
		}

		walletKey, errG := chainSimulatorConfigs.GenerateWalletKeyForShard(core.SovereignChainShardId, args.NumOfShards, addressConverter)
		if errG != nil {
			return nil, errG
		}

		supply := big.NewInt(0).Set(initialRewa)
		supply.Add(supply, big.NewInt(0).Set(initialStakedRewaPerNode))

		addresses = append(addresses, data.InitialAccount{
			Address:      walletKey.Address.Bech32,
			Balance:      initialRewa,
			Supply:       supply,
			StakingValue: big.NewInt(0).Set(initialStakedRewaPerNode),
		})

		initialWalletKeys.StakeWallets = append(initialWalletKeys.StakeWallets, walletKey)
		initialWalletKeys.BalanceWallets[uint32(i)] = walletKey
	}

	addressesBytes, errM := json.Marshal(addresses)
	if errM != nil {
		return nil, errM
	}

	err = os.WriteFile(configs.ConfigurationPathsHolder.Genesis, addressesBytes, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return initialWalletKeys, nil
}
