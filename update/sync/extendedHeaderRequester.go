package sync

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/process"
)

type extendedHeaderRequester struct {
	requestHandler ExtendedShardHeaderRequestHandler
}

// NewExtendedHeaderRequester creates an extended header requester wrapper
func NewExtendedHeaderRequester(requestHandler ExtendedShardHeaderRequestHandler) (*extendedHeaderRequester, error) {
	if check.IfNil(requestHandler) {
		return nil, process.ErrNilRequestHandler
	}

	return &extendedHeaderRequester{
		requestHandler: requestHandler,
	}, nil
}

// ShouldRequestHeader returns true if the shard id is main chain
func (ehr *extendedHeaderRequester) ShouldRequestHeader(shardId uint32) bool {
	return shardId == core.MainChainShardId
}

// RequestHeader requests extended shard header by hash
func (ehr *extendedHeaderRequester) RequestHeader(hash []byte) {
	ehr.requestHandler.RequestExtendedShardHeader(hash)
}

// IsInterfaceNil checks if underlying pointer is nil
func (ehr *extendedHeaderRequester) IsInterfaceNil() bool {
	return ehr == nil
}
