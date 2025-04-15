package block

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/process/block/bootstrapStorage"
)

type sovereignShardCrossNotarizer struct {
	*baseBlockNotarizer
}

func (scn *sovereignShardCrossNotarizer) getLastCrossNotarizedHeaders() []bootstrapStorage.BootstrapHeaderInfo {
	bootstrapHeaderInfo := scn.getLastCrossNotarizedHeadersForShard(core.MainChainShardId)
	if bootstrapHeaderInfo == nil {
		return nil
	}

	bootstrapHeaderInfo.ShardId = core.MainChainShardId

	lastCrossNotarizedHeaders := make([]bootstrapStorage.BootstrapHeaderInfo, 0, 1)
	lastCrossNotarizedHeaders = append(lastCrossNotarizedHeaders, *bootstrapHeaderInfo)
	return trimSliceBootstrapHeaderInfo(lastCrossNotarizedHeaders)
}
