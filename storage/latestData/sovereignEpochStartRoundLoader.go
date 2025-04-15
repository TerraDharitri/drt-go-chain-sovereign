package latestData

import (
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/storage"
)

type sovereignEpochStartRoundLoader struct {
	registryHandler MetaEpochStartTriggerRegistryHandler
}

func newSovereignEpochStartRoundLoader() *sovereignEpochStartRoundLoader {
	return &sovereignEpochStartRoundLoader{
		registryHandler: metachain.NewSovereignTriggerRegistryCreator(),
	}
}

// loadEpochStartRound will return the epoch start round from the bootstrap unit
func (l *sovereignEpochStartRoundLoader) loadEpochStartRound(
	_ uint32,
	key []byte,
	storer storage.Storer,
) (uint64, error) {
	trigInternalKey := append([]byte(common.TriggerRegistryKeyPrefix), key...)
	trigData, err := storer.Get(trigInternalKey)
	if err != nil {
		return 0, err
	}

	marshaller := &marshal.GogoProtoMarshalizer{}
	state, err := l.registryHandler.UnmarshalTrigger(marshaller, trigData)
	if err != nil {
		return 0, err
	}

	return state.GetCurrEpochStartRound(), nil
}
