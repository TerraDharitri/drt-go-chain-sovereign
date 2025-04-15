package latestData

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/epochStart/shardchain"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

type epochStartRoundLoader struct {
	registryHandler MetaEpochStartTriggerRegistryHandler
}

func newEpochStartRoundLoader() *epochStartRoundLoader {
	return &epochStartRoundLoader{
		registryHandler: metachain.NewMetaTriggerRegistryCreator(),
	}
}

// loadEpochStartRound will return the epoch start round from the bootstrap unit
func (l *epochStartRoundLoader) loadEpochStartRound(
	shardID uint32,
	key []byte,
	storer storage.Storer,
) (uint64, error) {
	trigInternalKey := append([]byte(common.TriggerRegistryKeyPrefix), key...)
	trigData, err := storer.Get(trigInternalKey)
	if err != nil {
		return 0, err
	}

	marshaller := &marshal.GogoProtoMarshalizer{}
	if shardID == core.MetachainShardId {
		state, err := l.registryHandler.UnmarshalTrigger(marshaller, trigData)
		if err != nil {
			return 0, err
		}

		return state.GetCurrEpochStartRound(), nil
	}

	var trigHandler data.TriggerRegistryHandler
	trigHandler, err = shardchain.UnmarshalTrigger(marshaller, trigData)
	if err != nil {
		return 0, err
	}

	return trigHandler.GetEpochStartRound(), nil
}
