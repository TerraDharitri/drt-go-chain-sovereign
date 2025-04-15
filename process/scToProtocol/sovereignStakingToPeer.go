package scToProtocol

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/process"
)

type sovereignStakingToPeer struct {
	*stakingToPeer
}

// NewSovereignStakingToPeer creates the staking to peer protocol for sovereign chain
func NewSovereignStakingToPeer(sp *stakingToPeer) (*sovereignStakingToPeer, error) {
	if check.IfNil(sp) {
		return nil, process.ErrNilSCToProtocol
	}

	sp.modifiedMBShardIDCheckerHandler = &sovModifiedMBShardIDChecker{}

	return &sovereignStakingToPeer{
		sp,
	}, nil
}
