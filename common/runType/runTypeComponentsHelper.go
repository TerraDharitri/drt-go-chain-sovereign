package runType

import (
	"github.com/TerraDharitri/drt-go-chain/genesis"
	"github.com/TerraDharitri/drt-go-chain/genesis/data"

	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// ReadInitialAccounts returns the genesis accounts from a file
func ReadInitialAccounts(filePath string) ([]genesis.InitialAccountHandler, error) {
	initialAccounts := make([]*data.InitialAccount, 0)
	err := core.LoadJsonFile(&initialAccounts, filePath)
	if err != nil {
		return nil, err
	}

	var accounts []genesis.InitialAccountHandler
	for _, ia := range initialAccounts {
		accounts = append(accounts, ia)
	}

	return accounts, nil
}
