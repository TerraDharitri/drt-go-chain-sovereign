package process

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	outportcore "github.com/TerraDharitri/drt-go-chain-core/data/outport"

	"github.com/TerraDharitri/drt-go-chain/process"
)

type sovereignOutportDataProvider struct {
	*outportDataProvider
}

// NewSovereignOutportDataProvider creates an outport data provider for sovereign chain
func NewSovereignOutportDataProvider(odp *outportDataProvider) (*sovereignOutportDataProvider, error) {
	if check.IfNil(odp) {
		return nil, process.ErrNilOutportDataProvider
	}

	sodp := &sovereignOutportDataProvider{
		odp,
	}
	sodp.shardRewardsCreator = sodp
	return sodp, nil
}

func (sodp *sovereignOutportDataProvider) getRewardsForShard(rewardsTxs map[string]data.TransactionHandler) (map[string]*outportcore.RewardInfo, error) {
	return getRewards(rewardsTxs)
}

// IsInterfaceNil returns true if there is no value under the interface
func (sodp *sovereignOutportDataProvider) IsInterfaceNil() bool {
	return sodp == nil
}
