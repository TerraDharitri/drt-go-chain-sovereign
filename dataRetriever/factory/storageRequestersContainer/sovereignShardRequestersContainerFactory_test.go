package storagerequesterscontainer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	storagerequesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/storageRequestersContainer"
	errorsDrt "github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/stretchr/testify/require"
)

func TestNewSovereignShardRequestersContainerFactory(t *testing.T) {
	t.Parallel()

	t.Run("nil input, should fail", func(t *testing.T) {
		sovContainer, err := storagerequesterscontainer.NewSovereignShardRequestersContainerFactory(nil)
		require.Nil(t, sovContainer)
		require.Equal(t, errorsDrt.ErrNilShardRequesterContainerFactory, err)
	})
	t.Run("should work", func(t *testing.T) {
		args := getArgumentsShard()
		shardContainer, _ := storagerequesterscontainer.NewShardRequestersContainerFactory(args)
		sovContainer, err := storagerequesterscontainer.NewSovereignShardRequestersContainerFactory(shardContainer)
		require.NotNil(t, sovContainer)
		require.Nil(t, err)
		require.False(t, sovContainer.IsInterfaceNil())
	})
}

func TestSovereignShardRequestersContainerFactory_Create(t *testing.T) {
	t.Parallel()

	args := getArgumentsShard()
	shardContainer, _ := storagerequesterscontainer.NewShardRequestersContainerFactory(args)
	sovContainer, _ := storagerequesterscontainer.NewSovereignShardRequestersContainerFactory(shardContainer)

	container, err := sovContainer.Create()
	require.Nil(t, err)

	numRequesterSCRs := 1
	numRequesterTxs := 1
	numRequesterRewardTxs := 0
	numRequesterHeaders := 1
	numRequesterMiniBlocks := 1
	numRequesterMetaBlockHeaders := 0
	numRequesterPeerAuth := 1
	numRequesterValidatorInfo := 1
	numRequesterExtendedHeader := 1
	numRequesters := numRequesterTxs + numRequesterHeaders + numRequesterMiniBlocks + numRequesterMetaBlockHeaders +
		numRequesterSCRs + numRequesterRewardTxs + numRequesterPeerAuth + numRequesterValidatorInfo + numRequesterExtendedHeader

	require.Equal(t, numRequesters, container.Len())

	sovShardIDStr := fmt.Sprintf("_%d", core.SovereignChainShardId)
	allKeys := map[string]struct{}{
		factory.TransactionTopic + sovShardIDStr:         {},
		factory.UnsignedTransactionTopic + sovShardIDStr: {},
		factory.ShardBlocksTopic + sovShardIDStr:         {},
		factory.MiniBlocksTopic + sovShardIDStr:          {},
		common.PeerAuthenticationTopic:                   {},
		common.ValidatorInfoTopic + sovShardIDStr:        {},
		factory.ExtendedHeaderProofTopic + sovShardIDStr: {},
	}
	iterateFunc := func(key string, requester dataRetriever.Requester) bool {
		require.False(t, strings.Contains(strings.ToLower(key), "meta"))
		delete(allKeys, key)
		return true
	}

	container.Iterate(iterateFunc)
	require.Empty(t, allKeys)
}
