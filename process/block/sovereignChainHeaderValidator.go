package block

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/process"
)

type sovereignChainHeaderValidator struct {
	*headerValidator
}

// NewSovereignChainHeaderValidator creates a new sovereign chain header validator
func NewSovereignChainHeaderValidator(
	headerValidator *headerValidator,
) (*sovereignChainHeaderValidator, error) {
	if headerValidator == nil {
		return nil, process.ErrNilHeaderValidator
	}

	schv := &sovereignChainHeaderValidator{
		headerValidator: headerValidator,
	}

	schv.calculateHeaderHashFunc = schv.calculateHeaderHash
	return schv, nil
}

func (schv *sovereignChainHeaderValidator) calculateHeaderHash(headerHandler data.HeaderHandler) ([]byte, error) {
	shardHeaderExtended, isShardHeaderExtended := headerHandler.(*block.ShardHeaderExtended)
	if isShardHeaderExtended {
		if check.IfNil(shardHeaderExtended.Header) {
			return nil, process.ErrNilHeaderHandler
		}

		return core.CalculateHash(schv.marshalizer, schv.hasher, shardHeaderExtended.Header)
	}

	return core.CalculateHash(schv.marshalizer, schv.hasher, headerHandler)
}
