package scToProtocol

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
)

type sovModifiedMBShardIDChecker struct {
}

func (c *sovModifiedMBShardIDChecker) isModifiedStateMBValid(miniBlock *block.MiniBlock) bool {
	return miniBlock.SenderShardID == core.SovereignChainShardId && miniBlock.ReceiverShardID == core.SovereignChainShardId
}
