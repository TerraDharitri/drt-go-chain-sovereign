package metachain

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
)

type sovereignTriggerRegistryCreator struct {
}

// NewSovereignTriggerRegistryCreator creates a new sovereign trigger registry creator/retriever
func NewSovereignTriggerRegistryCreator() *sovereignTriggerRegistryCreator {
	return &sovereignTriggerRegistryCreator{}
}

func (st *sovereignTriggerRegistryCreator) createRegistry(header data.HeaderHandler, t *trigger) (data.MetaTriggerRegistryHandler, error) {
	sovChainHdr, ok := header.(*block.SovereignChainHeader)
	if !ok {
		return nil, epochStart.ErrWrongTypeAssertion
	}

	registry := &block.SovereignShardTriggerRegistry{}
	registry.CurrentRound = t.currentRound
	registry.EpochFinalityAttestingRound = t.epochFinalityAttestingRound
	registry.CurrEpochStartRound = t.currEpochStartRound
	registry.PrevEpochStartRound = t.prevEpochStartRound
	registry.Epoch = t.epoch
	registry.EpochStartMetaHash = t.epochStartMetaHash
	registry.SovereignChainHeader = sovChainHdr

	return registry, nil
}

// UnmarshalTrigger returns the sovereign shard trigger registry from provided data
func (st *sovereignTriggerRegistryCreator) UnmarshalTrigger(marshaller marshal.Marshalizer, data []byte) (data.MetaTriggerRegistryHandler, error) {
	state := &block.SovereignShardTriggerRegistry{
		SovereignChainHeader: &block.SovereignChainHeader{},
	}

	err := marshaller.Unmarshal(state, data)
	if err != nil {
		return nil, err
	}

	return state, nil
}
