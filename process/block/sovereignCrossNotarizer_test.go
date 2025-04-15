package block

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/process/block/bootstrapStorage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/stretchr/testify/require"
)

func TestSovereignShardCrossNotarizer_getLastCrossNotarizedHeaders(t *testing.T) {
	hash := []byte("hash")
	header := &block.SovereignChainHeader{
		Header: &block.Header{
			ShardID: core.SovereignChainShardId,
			Nonce:   4,
		},
	}
	sovereignNotarzier := &sovereignShardCrossNotarizer{
		&baseBlockNotarizer{
			blockTracker: &testscommon.BlockTrackerStub{
				GetLastCrossNotarizedHeaderCalled: func(shardID uint32) (data.HeaderHandler, []byte, error) {
					require.Equal(t, core.MainChainShardId, shardID)
					return header, hash, nil
				},
			},
		},
	}

	headers := sovereignNotarzier.getLastCrossNotarizedHeaders()
	expectedHeaders := []bootstrapStorage.BootstrapHeaderInfo{
		{
			ShardId: core.MainChainShardId,
			Nonce:   header.GetNonce(),
			Hash:    hash,
		},
	}
	require.Equal(t, expectedHeaders, headers)
}
