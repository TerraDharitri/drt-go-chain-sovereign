package sovereign

import "github.com/TerraDharitri/drt-go-chain-core/data/sovereign"

// DataCodecMock -
type DataCodecMock struct {
	SerializeEventDataCalled   func(eventData sovereign.EventData) ([]byte, error)
	DeserializeEventDataCalled func(data []byte) (*sovereign.EventData, error)
	SerializeTokenDataCalled   func(tokenData sovereign.DcdtTokenData) ([]byte, error)
	DeserializeTokenDataCalled func(data []byte) (*sovereign.DcdtTokenData, error)
	GetTokenDataBytesCalled    func(tokenNonce []byte, tokenData []byte) ([]byte, error)
	SerializeOperationCalled   func(operation sovereign.Operation) ([]byte, error)
}

// SerializeEventData -
func (dcm *DataCodecMock) SerializeEventData(eventData sovereign.EventData) ([]byte, error) {
	if dcm.SerializeEventDataCalled != nil {
		return dcm.SerializeEventDataCalled(eventData)
	}

	return make([]byte, 0), nil
}

// DeserializeEventData -
func (dcm *DataCodecMock) DeserializeEventData(data []byte) (*sovereign.EventData, error) {
	if dcm.DeserializeEventDataCalled != nil {
		return dcm.DeserializeEventDataCalled(data)
	}

	return &sovereign.EventData{}, nil
}

// SerializeTokenData -
func (dcm *DataCodecMock) SerializeTokenData(tokenData sovereign.DcdtTokenData) ([]byte, error) {
	if dcm.SerializeTokenDataCalled != nil {
		return dcm.SerializeTokenDataCalled(tokenData)
	}

	return make([]byte, 0), nil
}

// DeserializeTokenData -
func (dcm *DataCodecMock) DeserializeTokenData(data []byte) (*sovereign.DcdtTokenData, error) {
	if dcm.DeserializeTokenDataCalled != nil {
		return dcm.DeserializeTokenDataCalled(data)
	}

	return &sovereign.DcdtTokenData{}, nil
}

// GetTokenDataBytes -
func (dcm *DataCodecMock) GetTokenDataBytes(tokenNonce []byte, tokenData []byte) ([]byte, error) {
	if dcm.GetTokenDataBytesCalled != nil {
		return dcm.GetTokenDataBytesCalled(tokenNonce, tokenData)
	}

	return make([]byte, 0), nil
}

// SerializeOperation -
func (dcm *DataCodecMock) SerializeOperation(operation sovereign.Operation) ([]byte, error) {
	if dcm.SerializeOperationCalled != nil {
		return dcm.SerializeOperationCalled(operation)
	}

	return make([]byte, 0), nil
}

// IsInterfaceNil -
func (dcm *DataCodecMock) IsInterfaceNil() bool {
	return dcm == nil
}
