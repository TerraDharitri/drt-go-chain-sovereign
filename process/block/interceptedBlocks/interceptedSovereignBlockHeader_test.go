package interceptedBlocks_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/interceptedBlocks"
	"github.com/stretchr/testify/require"
)

func TestNewSovereignInterceptedBlockHeader(t *testing.T) {
	t.Parallel()

	t.Run("nil input in args", func(t *testing.T) {
		args := createDefaultShardArgument()
		args.Marshalizer = nil
		sovInterceptedBlock, err := interceptedBlocks.NewSovereignInterceptedBlockHeader(args)
		require.Nil(t, sovInterceptedBlock)
		require.Equal(t, process.ErrNilMarshalizer, err)
	})
	t.Run("should work", func(t *testing.T) {
		args := createDefaultShardArgument()
		sovInterceptedBlock, err := interceptedBlocks.NewSovereignInterceptedBlockHeader(args)
		require.Nil(t, err)
		require.False(t, sovInterceptedBlock.IsInterfaceNil())
	})
}

func TestInterceptedSovereignMiniBlock_checkMiniBlocksHeaders(t *testing.T) {
	t.Parallel()

	args := createDefaultShardArgument()
	sovInterceptedBlock, _ := interceptedBlocks.NewSovereignInterceptedBlockHeader(args)

	miniBlockHeaders := []data.MiniBlockHeaderHandler{
		&block.MiniBlockHeader{
			Hash:            make([]byte, 0),
			SenderShardID:   core.SovereignChainShardId,
			ReceiverShardID: core.SovereignChainShardId,
			TxCount:         0,
			Type:            0,
			Reserved:        []byte("rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr"),
		},
	}

	err := sovInterceptedBlock.CheckMiniBlocksHeaders(miniBlockHeaders, args.ShardCoordinator)
	require.Nil(t, err)

	err = miniBlockHeaders[0].SetReceiverShardID(core.MainChainShardId)
	require.Nil(t, err)
	err = sovInterceptedBlock.CheckMiniBlocksHeaders(miniBlockHeaders, args.ShardCoordinator)
	require.Nil(t, err)

	err = miniBlockHeaders[0].SetSenderShardID(core.MainChainShardId)
	require.Nil(t, err)
	err = sovInterceptedBlock.CheckMiniBlocksHeaders(miniBlockHeaders, args.ShardCoordinator)
	require.Nil(t, err)

	err = miniBlockHeaders[0].SetReceiverShardID(core.MetachainShardId)
	require.Nil(t, err)
	err = sovInterceptedBlock.CheckMiniBlocksHeaders(miniBlockHeaders, args.ShardCoordinator)
	require.Equal(t, process.ErrInvalidShardId, err)
}
