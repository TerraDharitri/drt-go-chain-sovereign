package dataRetriever

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	retriever "github.com/TerraDharitri/drt-go-chain/dataRetriever"
	mockRetriever "github.com/TerraDharitri/drt-go-chain/dataRetriever/mock"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/stretchr/testify/require"
)

func TestDataRetrieverContainersSetter_SetEpochHandlerToMetaBlockContainers(t *testing.T) {
	t.Parallel()

	sovSetter := NewSovereignDataRetrieverContainerSetter()
	require.False(t, sovSetter.IsInterfaceNil())

	expectedTopic := fmt.Sprintf("%s_%d", factory.ShardBlocksTopic, core.SovereignChainShardId)
	triggerStub := &testscommon.EpochStartTriggerStub{}

	wasEpochHandlerSetInResolver := false
	resolver := &mockRetriever.HeaderResolverStub{
		SetEpochHandlerCalled: func(epochHandler retriever.EpochHandler) error {
			require.Equal(t, triggerStub, epochHandler)
			wasEpochHandlerSetInResolver = true
			return nil
		},
	}
	resolversContainer := &dataRetriever.ResolversContainerStub{
		GetCalled: func(key string) (retriever.Resolver, error) {
			require.Equal(t, expectedTopic, key)
			return resolver, nil
		},
	}

	wasEpochHandlerSetInRequester := false
	requester := &mockRetriever.HeaderRequesterStub{
		SetEpochHandlerCalled: func(epochHandler retriever.EpochHandler) error {
			require.Equal(t, triggerStub, epochHandler)
			wasEpochHandlerSetInRequester = true
			return nil
		},
	}
	requestersContainer := &dataRetriever.RequestersContainerStub{
		GetCalled: func(key string) (retriever.Requester, error) {
			require.Equal(t, expectedTopic, key)
			return requester, nil
		},
	}
	err := sovSetter.SetEpochHandlerToMetaBlockContainers(triggerStub, resolversContainer, requestersContainer)
	require.Nil(t, err)
	require.True(t, wasEpochHandlerSetInResolver)
	require.True(t, wasEpochHandlerSetInRequester)
}
