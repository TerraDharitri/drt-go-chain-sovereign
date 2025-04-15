package runType

import (
	"github.com/TerraDharitri/drt-go-chain/config"
	mainFactory "github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/genesis"
	"github.com/TerraDharitri/drt-go-chain/genesis/data"

	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// CreateArgsRunTypeComponents creates the args for run type component
func CreateArgsRunTypeComponents(
	coreComponents mainFactory.CoreComponentsHandler,
	cryptoComponents mainFactory.CryptoComponentsHandler,
	configs config.Configs,
) (*ArgsRunTypeComponents, error) {
	initialAccounts := make([]*data.InitialAccount, 0)
	err := core.LoadJsonFile(&initialAccounts, configs.ConfigurationPathsHolder.Genesis)
	if err != nil {
		return nil, err
	}

	var accounts []genesis.InitialAccountHandler
	for _, ia := range initialAccounts {
		accounts = append(accounts, ia)
	}

	return &ArgsRunTypeComponents{
		CoreComponents:   coreComponents,
		CryptoComponents: cryptoComponents,
		Configs:          configs,
		InitialAccounts:  accounts,
	}, nil
}
