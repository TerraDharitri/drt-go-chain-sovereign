package common

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"

	sovRunType "github.com/TerraDharitri/drt-go-chain/cmd/sovereignnode/runType"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/factory/runType"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/process"
)

// GetCurrentSovereignHeader returns current sovereign chain block handler from blockchain hook
func GetCurrentSovereignHeader(nodeHandler process.NodeHandler) data.SovereignChainHeaderHandler {
	return nodeHandler.GetChainHandler().GetCurrentBlockHeader().(data.SovereignChainHeaderHandler)
}

// CreateSovereignRunTypeComponents will create sovereign run type components
func CreateSovereignRunTypeComponents(args runType.ArgsRunTypeComponents, sovereignExtraConfig config.SovereignConfig) (factory.RunTypeComponentsHolder, error) {
	argsSovRunType, err := sovRunType.CreateSovereignArgsRunTypeComponents(args, sovereignExtraConfig)
	if err != nil {
		return nil, err
	}

	sovereignComponentsFactory, err := runType.NewSovereignRunTypeComponentsFactory(*argsSovRunType)
	if err != nil {
		return nil, err
	}

	managedRunTypeComponents, err := runType.NewManagedRunTypeComponents(sovereignComponentsFactory)
	if err != nil {
		return nil, err
	}
	err = managedRunTypeComponents.Create()
	if err != nil {
		return nil, err
	}

	return managedRunTypeComponents, nil
}
