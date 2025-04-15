package sovereign

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
)

// OutgoingOperationsFormatter collects relevant outgoing events for bridge from the logs and creates outgoing data
// that needs to be signed by validators to bridge tokens
type OutgoingOperationsFormatter interface {
	CreateOutgoingTxsData(logs []*data.LogData) ([][]byte, error)
	IsInterfaceNil() bool
}

// DataCodecHandler is the interface for serializing/deserializing data
type DataCodecHandler interface {
	SerializeEventData(eventData sovereign.EventData) ([]byte, error)
	DeserializeEventData(data []byte) (*sovereign.EventData, error)
	SerializeTokenData(tokenData sovereign.DcdtTokenData) ([]byte, error)
	DeserializeTokenData(data []byte) (*sovereign.DcdtTokenData, error)
	SerializeOperation(operation sovereign.Operation) ([]byte, error)
	IsInterfaceNil() bool
}

// TopicsCheckerHandler should be able to check the topics validity
type TopicsCheckerHandler interface {
	CheckValidity(topics [][]byte) error
	IsInterfaceNil() bool
}
