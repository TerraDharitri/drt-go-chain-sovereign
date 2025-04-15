package sovereign

import "github.com/TerraDharitri/drt-go-chain-core/data"

// OutgoingOperationsFormatterMock -
type OutgoingOperationsFormatterMock struct {
	CreateOutgoingTxDataCalled func(logs []*data.LogData) ([][]byte, error)
}

// CreateOutgoingTxsData -
func (stub *OutgoingOperationsFormatterMock) CreateOutgoingTxsData(logs []*data.LogData) ([][]byte, error) {
	if stub.CreateOutgoingTxDataCalled != nil {
		return stub.CreateOutgoingTxDataCalled(logs)
	}

	return make([][]byte, 0), nil
}

// IsInterfaceNil -
func (stub *OutgoingOperationsFormatterMock) IsInterfaceNil() bool {
	return stub == nil
}
