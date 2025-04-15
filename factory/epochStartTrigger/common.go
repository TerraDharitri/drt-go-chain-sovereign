package epochStartTrigger

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/check"

	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/factory"
	"github.com/TerraDharitri/drt-go-chain/process"
)

func checkNilArgs(args factory.ArgsEpochStartTrigger) error {
	if check.IfNil(args.DataComps) {
		return process.ErrNilDataComponentsHolder
	}
	if check.IfNil(args.DataComps.Datapool()) {
		return process.ErrNilDataPoolHolder
	}
	if check.IfNil(args.DataComps.Blockchain()) {
		return process.ErrNilBlockChain
	}
	if check.IfNil(args.DataComps.Datapool().MiniBlocks()) {
		return dataRetriever.ErrNilMiniblocksPool
	}
	if check.IfNil(args.DataComps.Datapool().ValidatorsInfo()) {
		return process.ErrNilValidatorInfoPool
	}
	if check.IfNil(args.BootstrapComponents) {
		return process.ErrNilBootstrapComponentsHolder
	}
	if check.IfNil(args.BootstrapComponents.ShardCoordinator()) {
		return process.ErrNilShardCoordinator
	}
	if check.IfNil(args.RequestHandler) {
		return process.ErrNilRequestHandler
	}

	return nil
}
