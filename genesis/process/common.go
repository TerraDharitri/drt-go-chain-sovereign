package process

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/genesis"
	"github.com/TerraDharitri/drt-go-chain/update"
)

type headerCreatorArgs struct {
	mapArgsGenesisBlockCreator map[uint32]ArgsGenesisBlockCreator
	mapHardForkBlockProcessor  map[uint32]update.HardForkBlockProcessor
	mapBodies                  map[uint32]*block.Body
	shardIDs                   []uint32
	nodesListSplitter          genesis.NodesListSplitter
}
