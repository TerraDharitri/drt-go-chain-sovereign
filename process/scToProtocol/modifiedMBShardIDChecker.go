package scToProtocol

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
)

type modifiedMBShardIDChecker struct {
}

func (c *modifiedMBShardIDChecker) isModifiedStateMBValid(miniBlock *block.MiniBlock) bool {
	return miniBlock.SenderShardID == core.MetachainShardId && miniBlock.ReceiverShardID == core.MetachainShardId
}
