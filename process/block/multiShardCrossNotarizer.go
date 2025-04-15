package block

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/process/block/bootstrapStorage"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

type multiShardCrossNotarizer struct {
	*baseBlockNotarizer
	shardCoordinator sharding.Coordinator
}

func (mcn *multiShardCrossNotarizer) getLastCrossNotarizedHeaders() []bootstrapStorage.BootstrapHeaderInfo {
	lastCrossNotarizedHeaders := make([]bootstrapStorage.BootstrapHeaderInfo, 0, mcn.shardCoordinator.NumberOfShards()+1)

	for shardID := uint32(0); shardID < mcn.shardCoordinator.NumberOfShards(); shardID++ {
		bootstrapHeaderInfo := mcn.getLastCrossNotarizedHeadersForShard(shardID)
		if bootstrapHeaderInfo != nil {
			lastCrossNotarizedHeaders = append(lastCrossNotarizedHeaders, *bootstrapHeaderInfo)
		}
	}

	bootstrapHeaderInfo := mcn.getLastCrossNotarizedHeadersForShard(core.MetachainShardId)
	if bootstrapHeaderInfo != nil {
		lastCrossNotarizedHeaders = append(lastCrossNotarizedHeaders, *bootstrapHeaderInfo)
	}

	if len(lastCrossNotarizedHeaders) == 0 {
		return nil
	}

	return trimSliceBootstrapHeaderInfo(lastCrossNotarizedHeaders)
}
