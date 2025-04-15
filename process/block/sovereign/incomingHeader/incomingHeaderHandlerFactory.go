package incomingHeader

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	hasherFactory "github.com/TerraDharitri/drt-go-chain-core/hashing/factory"
	marshallerFactory "github.com/TerraDharitri/drt-go-chain-core/marshal/factory"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	errorsDrt "github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// CreateIncomingHeaderProcessor creates the incoming header processor
func CreateIncomingHeaderProcessor(
	config config.WebSocketConfig,
	dataPool dataRetriever.PoolsHolder,
	mainChainNotarizationStartRound uint64,
	runTypeComponents RunTypeComponentsHolder,
) (process.IncomingHeaderSubscriber, error) {
	if check.IfNil(runTypeComponents) {
		return nil, errorsDrt.ErrNilRunTypeComponents
	}
	marshaller, err := marshallerFactory.NewMarshalizer(config.MarshallerType)
	if err != nil {
		return nil, err
	}
	hasher, err := hasherFactory.NewHasher(config.HasherType)
	if err != nil {
		return nil, err
	}

	argsIncomingHeaderHandler := ArgsIncomingHeaderProcessor{
		HeadersPool:                     dataPool.Headers(),
		TxPool:                          dataPool.UnsignedTransactions(),
		Marshaller:                      marshaller,
		Hasher:                          hasher,
		MainChainNotarizationStartRound: mainChainNotarizationStartRound,
		OutGoingOperationsPool:          runTypeComponents.OutGoingOperationsPoolHandler(),
		DataCodec:                       runTypeComponents.DataCodecHandler(),
		TopicsChecker:                   runTypeComponents.TopicsCheckerHandler(),
	}

	return NewIncomingHeaderProcessor(argsIncomingHeaderHandler)
}
