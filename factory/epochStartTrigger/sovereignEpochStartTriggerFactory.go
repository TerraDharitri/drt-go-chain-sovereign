package epochStartTrigger

import (
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/epochStart/shardchain"
	"github.com/TerraDharitri/drt-go-chain/factory"
)

type sovereignEpochStartTriggerFactory struct {
}

// NewSovereignEpochStartTriggerFactory creates a sovereign epoch start trigger. This will be a metachain one, since
// nodes inside sovereign chain will not need to receive meta information, but they will actually execute the meta code.
func NewSovereignEpochStartTriggerFactory() *sovereignEpochStartTriggerFactory {
	return &sovereignEpochStartTriggerFactory{}
}

// CreateEpochStartTrigger creates a meta epoch start trigger for sovereign run type
func (f *sovereignEpochStartTriggerFactory) CreateEpochStartTrigger(args factory.ArgsEpochStartTrigger) (epochStart.TriggerHandler, error) {
	err := checkNilArgs(args)
	if err != nil {
		return nil, err
	}

	metaTriggerArgs, err := createMetaEpochStartTriggerArgs(args)
	if err != nil {
		return nil, err
	}

	argsPeerMiniBlockSyncer := shardchain.ArgPeerMiniBlockSyncer{
		MiniBlocksPool:     args.DataComps.Datapool().MiniBlocks(),
		ValidatorsInfoPool: args.DataComps.Datapool().ValidatorsInfo(),
		RequestHandler:     args.RequestHandler,
	}
	peerMiniBlockSyncer, err := shardchain.NewPeerMiniBlockSyncer(argsPeerMiniBlockSyncer)
	if err != nil {
		return nil, err
	}

	argsSovTrigger := metachain.ArgsSovereignTrigger{
		ArgsNewMetaEpochStartTrigger: metaTriggerArgs,
		ValidatorInfoSyncer:          peerMiniBlockSyncer,
	}

	return metachain.NewSovereignTrigger(argsSovTrigger)
}

// IsInterfaceNil checks if the underlying pointer is nil
func (f *sovereignEpochStartTriggerFactory) IsInterfaceNil() bool {
	return f == nil
}
