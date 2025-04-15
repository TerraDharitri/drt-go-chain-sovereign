package interceptedBlocks

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// IsMetaHeaderOutOfRange -
func (imh *InterceptedMetaHeader) IsMetaHeaderOutOfRange() bool {
	return imh.isMetaHeaderEpochOutOfRange()
}

// CheckMiniBlocksHeaders -
func (isbh *interceptedSovereignBlockHeader) CheckMiniBlocksHeaders(mbHeaders []data.MiniBlockHeaderHandler, coordinator sharding.Coordinator) error {
	return isbh.checkMiniBlocksHeaders(mbHeaders, coordinator)
}
