package chainSimulator

import (
	"github.com/TerraDharitri/drt-go-chain/process"

	"github.com/TerraDharitri/drt-go-chain-core/data"
)

// BlockProcessorMock -
type BlockProcessorMock struct {
	ProcessBlockCalled func(blockProcessor process.BlockProcessor, header data.HeaderHandler) (data.HeaderHandler, data.BodyHandler, error)
}

// ProcessBlock -
func (mock *BlockProcessorMock) ProcessBlock(blockProcessor process.BlockProcessor, header data.HeaderHandler) (data.HeaderHandler, data.BodyHandler, error) {
	if mock.ProcessBlockCalled != nil {
		return mock.ProcessBlockCalled(blockProcessor, header)
	}
	return nil, nil, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (mock *BlockProcessorMock) IsInterfaceNil() bool {
	return mock == nil
}
