package dataCodec

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
)

// AbiSerializer is the interface to work with abi codec
type AbiSerializer interface {
	Serialize(inputValues []any) (string, error)
	Deserialize(data string, outputValues []any) error
}

// EventDataEncoder is the interface for serializing/deserializing event data
type EventDataEncoder interface {
	SerializeEventData(eventData sovereign.EventData) ([]byte, error)
	DeserializeEventData(data []byte) (*sovereign.EventData, error)
}

// TokenDataEncoder is the interface for serializing/deserializing token data
type TokenDataEncoder interface {
	SerializeTokenData(tokenData sovereign.DcdtTokenData) ([]byte, error)
	DeserializeTokenData(data []byte) (*sovereign.DcdtTokenData, error)
}

// OperationDataEncoder is the interface for serializing operations
type OperationDataEncoder interface {
	SerializeOperation(operation sovereign.Operation) ([]byte, error)
}
