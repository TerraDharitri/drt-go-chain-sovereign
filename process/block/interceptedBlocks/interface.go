package interceptedBlocks

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

type mbHeadersChecker interface {
	checkMiniBlocksHeaders(mbHeaders []data.MiniBlockHeaderHandler, coordinator sharding.Coordinator) error
}
