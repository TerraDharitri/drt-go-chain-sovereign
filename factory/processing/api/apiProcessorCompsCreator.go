package api

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	metachainEpochStart "github.com/TerraDharitri/drt-go-chain/epochStart/metachain"
	"github.com/TerraDharitri/drt-go-chain/epochStart/notifier"
	factoryDisabled "github.com/TerraDharitri/drt-go-chain/factory/disabled"

	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/peer"
	"github.com/TerraDharitri/drt-go-chain/sharding"
)

// APIProcessComps is a struct placeholder for created api process comps
type APIProcessComps struct {
	StakingDataProviderAPI peer.StakingDataProviderAPI
	AuctionListSelector    epochStart.AuctionListSelector
}

// ArgsCreateAPIProcessComps is a struct placeholder for args needed to create api process comps
type ArgsCreateAPIProcessComps struct {
	ArgsStakingDataProvider metachainEpochStart.StakingDataProviderArgs
	ShardCoordinator        sharding.Coordinator
	EpochNotifier           process.EpochNotifier
	SoftAuctionConfig       config.SoftAuctionConfig
	EnableEpochs            config.EnableEpochs
	Denomination            int
}

type apiProcessorCompsCreator struct {
}

// NewAPIProcessorCompsCreator creates a new api processor components creator for regular chain(meta+shard)
func NewAPIProcessorCompsCreator() *apiProcessorCompsCreator {
	return &apiProcessorCompsCreator{}
}

// CreateAPIComps creates api process comps for metachain, otherwise returns disabled comps
func (c *apiProcessorCompsCreator) CreateAPIComps(args ArgsCreateAPIProcessComps) (*APIProcessComps, error) {
	if args.ShardCoordinator.SelfId() == core.MetachainShardId {
		return createAPIComps(args)
	}

	return &APIProcessComps{
		StakingDataProviderAPI: factoryDisabled.NewDisabledStakingDataProvider(),
		AuctionListSelector:    factoryDisabled.NewDisabledAuctionListSelector(),
	}, nil
}

func createAPIComps(args ArgsCreateAPIProcessComps) (*APIProcessComps, error) {
	stakingDataProviderAPI, err := metachainEpochStart.NewStakingDataProvider(args.ArgsStakingDataProvider)
	if err != nil {
		return nil, err
	}

	maxNodesChangeConfigProviderAPI, err := notifier.NewNodesConfigProviderAPI(args.EpochNotifier, args.EnableEpochs)
	if err != nil {
		return nil, err
	}

	extendedShardCoordinator, castOk := args.ShardCoordinator.(metachainEpochStart.ExtendedShardCoordinatorHandler)
	if !castOk {
		return nil, fmt.Errorf("%w when trying to cast shard coordinator to extended shard coordinator", process.ErrWrongTypeAssertion)
	}

	argsAuctionListSelectorAPI := metachainEpochStart.AuctionListSelectorArgs{
		ShardCoordinator:             extendedShardCoordinator,
		StakingDataProvider:          stakingDataProviderAPI,
		MaxNodesChangeConfigProvider: maxNodesChangeConfigProviderAPI,
		SoftAuctionConfig:            args.SoftAuctionConfig,
		Denomination:                 args.Denomination,
		AuctionListDisplayHandler:    factoryDisabled.NewDisabledAuctionListDisplayer(),
	}

	auctionListSelectorAPI, err := metachainEpochStart.NewAuctionListSelector(argsAuctionListSelectorAPI)
	if err != nil {
		return nil, err
	}

	return &APIProcessComps{
		StakingDataProviderAPI: stakingDataProviderAPI,
		AuctionListSelector:    auctionListSelectorAPI,
	}, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (c *apiProcessorCompsCreator) IsInterfaceNil() bool {
	return c == nil
}
