package incomingHeader

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"

	"github.com/TerraDharitri/drt-go-chain/process"
)

type extendedHeaderProcessor struct {
	headersPool HeadersPool
	txPool      TransactionPool
	marshaller  marshal.Marshalizer
	hasher      hashing.Hasher
}

func createExtendedHeader(incomingHeader sovereign.IncomingHeaderHandler, scrs []*SCRInfo) (*block.ShardHeaderExtended, error) {
	headerV2, castOk := incomingHeader.GetHeaderHandler().(*block.HeaderV2)
	if !castOk {
		return nil, errInvalidHeaderType
	}
	events, err := getEvents(incomingHeader.GetIncomingEventHandlers())
	if err != nil {
		return nil, err
	}

	return &block.ShardHeaderExtended{
		Header:             headerV2,
		IncomingMiniBlocks: createIncomingMb(scrs),
		IncomingEvents:     events,
	}, nil
}

func getEvents(events []data.EventHandler) ([]*transaction.Event, error) {
	ret := make([]*transaction.Event, len(events))

	for idx, eventHandler := range events {
		event, castOk := eventHandler.(*transaction.Event)
		if !castOk {
			return nil, errInvalidEventType
		}

		ret[idx] = event
	}

	return ret, nil
}

func createIncomingMb(scrs []*SCRInfo) []*block.MiniBlock {
	if len(scrs) == 0 {
		return make([]*block.MiniBlock, 0)
	}

	scrHashes := make([][]byte, len(scrs))
	for idx, scrData := range scrs {
		scrHashes[idx] = scrData.Hash
	}

	return []*block.MiniBlock{
		{
			TxHashes:        scrHashes,
			ReceiverShardID: core.SovereignChainShardId,
			SenderShardID:   core.MainChainShardId,
			Type:            block.SmartContractResultBlock,
		},
	}
}

func (ehp *extendedHeaderProcessor) addPreGenesisExtendedHeaderToPool(incomingHeader sovereign.IncomingHeaderHandler) error {
	headerV2, castOk := incomingHeader.GetHeaderHandler().(*block.HeaderV2)
	if !castOk {
		return errInvalidHeaderType
	}

	extendedHeader := &block.ShardHeaderExtended{
		Header:             headerV2,
		IncomingMiniBlocks: []*block.MiniBlock{},
		IncomingEvents:     []*transaction.Event{},
	}

	return ehp.addExtendedHeaderAndSCRsToPool(extendedHeader, make([]*SCRInfo, 0))
}

func (ehp *extendedHeaderProcessor) addExtendedHeaderAndSCRsToPool(extendedHeader data.ShardHeaderExtendedHandler, scrs []*SCRInfo) error {
	extendedHeaderHash, err := core.CalculateHash(ehp.marshaller, ehp.hasher, extendedHeader)
	if err != nil {
		return err
	}

	ehp.addSCRsToPool(scrs)
	ehp.headersPool.AddHeaderInShard(extendedHeaderHash, extendedHeader, core.MainChainShardId)
	return nil
}

func (ehp *extendedHeaderProcessor) addSCRsToPool(scrs []*SCRInfo) {
	cacheID := process.ShardCacherIdentifier(core.MainChainShardId, core.SovereignChainShardId)

	for _, scrData := range scrs {
		ehp.txPool.AddData(scrData.Hash, scrData.SCR, scrData.SCR.Size(), cacheID)
	}
}
