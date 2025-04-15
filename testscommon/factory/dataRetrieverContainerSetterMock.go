package factory

import (
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
)

// DataRetrieverContainersSetterMock -
type DataRetrieverContainersSetterMock struct {
	SetEpochHandlerToMetaBlockContainersCalled func(
		epochStartTrigger epochStart.TriggerHandler,
		resolversContainer dataRetriever.ResolversContainer,
		requestersContainer dataRetriever.RequestersContainer,
	) error
}

// SetEpochHandlerToMetaBlockContainers -
func (mock *DataRetrieverContainersSetterMock) SetEpochHandlerToMetaBlockContainers(
	epochStartTrigger epochStart.TriggerHandler,
	resolversContainer dataRetriever.ResolversContainer,
	requestersContainer dataRetriever.RequestersContainer,
) error {
	if mock.SetEpochHandlerToMetaBlockContainersCalled != nil {
		return mock.SetEpochHandlerToMetaBlockContainersCalled(epochStartTrigger, resolversContainer, requestersContainer)
	}

	return nil
}

// IsInterfaceNil -
func (mock *DataRetrieverContainersSetterMock) IsInterfaceNil() bool {
	return mock == nil
}
