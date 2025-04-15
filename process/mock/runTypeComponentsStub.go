package mock

import (
	sovereignBlock "github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool/sovereign"
	requesterscontainer "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/requestersContainer"
	storageRequestFactory "github.com/TerraDharitri/drt-go-chain/dataRetriever/factory/storageRequestersContainer/factory"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/requestHandlers"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/block/sovereign"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/state"
	syncerFactory "github.com/TerraDharitri/drt-go-chain/state/syncer/factory"
	"github.com/TerraDharitri/drt-go-chain/storage/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	sovereignMocks "github.com/TerraDharitri/drt-go-chain/testscommon/sovereign"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
)

// RunTypeComponentsStub -
type RunTypeComponentsStub struct {
	AdditionalStorageServiceFactory             process.AdditionalStorageServiceCreator
	ShardCoordinatorFactory                     sharding.ShardCoordinatorFactory
	NodesCoordinatorWithRaterFactory            nodesCoordinator.NodesCoordinatorWithRaterFactory
	RequestHandlerFactory                       requestHandlers.RequestHandlerCreator
	AccountCreator                              state.AccountFactory
	OutGoingOperationsPool                      sovereignBlock.OutGoingOperationsPool
	DataCodec                                   sovereign.DataCodecHandler
	TopicsChecker                               sovereign.TopicsCheckerHandler
	RequestersContainerFactoryCreatorField      requesterscontainer.RequesterContainerFactoryCreator
	ValidatorAccountsSyncerFactoryHandlerField  syncerFactory.ValidatorAccountsSyncerFactoryHandler
	ShardRequestersContainerCreatorHandlerField storageRequestFactory.ShardRequestersContainerCreatorHandler
}

// NewRunTypeComponentsStub -
func NewRunTypeComponentsStub() *RunTypeComponentsStub {
	return &RunTypeComponentsStub{
		AdditionalStorageServiceFactory:             &testscommon.AdditionalStorageServiceFactoryMock{},
		ShardCoordinatorFactory:                     sharding.NewMultiShardCoordinatorFactory(),
		NodesCoordinatorWithRaterFactory:            nodesCoordinator.NewIndexHashedNodesCoordinatorWithRaterFactory(),
		RequestHandlerFactory:                       requestHandlers.NewResolverRequestHandlerFactory(),
		AccountCreator:                              &stateMock.AccountsFactoryStub{},
		OutGoingOperationsPool:                      &sovereignMocks.OutGoingOperationsPoolMock{},
		DataCodec:                                   &sovereignMocks.DataCodecMock{},
		TopicsChecker:                               &sovereignMocks.TopicsCheckerMock{},
		RequestersContainerFactoryCreatorField:      requesterscontainer.NewShardRequestersContainerFactoryCreator(),
		ValidatorAccountsSyncerFactoryHandlerField:  syncerFactory.NewValidatorAccountsSyncerFactory(),
		ShardRequestersContainerCreatorHandlerField: storageRequestFactory.NewShardRequestersContainerCreator(),
	}
}

// NewSovereignRunTypeComponentsStub -
func NewSovereignRunTypeComponentsStub() *RunTypeComponentsStub {
	rt := NewRunTypeComponentsStub()
	requestHandlerFactory, _ := requestHandlers.NewSovereignResolverRequestHandlerFactory(rt.RequestHandlerFactory)

	return &RunTypeComponentsStub{
		AdditionalStorageServiceFactory:             factory.NewSovereignAdditionalStorageServiceFactory(),
		ShardCoordinatorFactory:                     sharding.NewSovereignShardCoordinatorFactory(),
		NodesCoordinatorWithRaterFactory:            &testscommon.NodesCoordinatorFactoryMock{},
		RequestHandlerFactory:                       requestHandlerFactory,
		AccountCreator:                              &stateMock.AccountsFactoryStub{},
		OutGoingOperationsPool:                      &sovereignMocks.OutGoingOperationsPoolMock{},
		DataCodec:                                   &sovereignMocks.DataCodecMock{},
		TopicsChecker:                               &sovereignMocks.TopicsCheckerMock{},
		RequestersContainerFactoryCreatorField:      requesterscontainer.NewSovereignShardRequestersContainerFactoryCreator(),
		ValidatorAccountsSyncerFactoryHandlerField:  syncerFactory.NewSovereignValidatorAccountsSyncerFactory(),
		ShardRequestersContainerCreatorHandlerField: storageRequestFactory.NewSovereignShardRequestersContainerCreator(),
	}
}

// AdditionalStorageServiceCreator -
func (r *RunTypeComponentsStub) AdditionalStorageServiceCreator() process.AdditionalStorageServiceCreator {
	return r.AdditionalStorageServiceFactory
}

// ShardCoordinatorCreator -
func (r *RunTypeComponentsStub) ShardCoordinatorCreator() sharding.ShardCoordinatorFactory {
	return r.ShardCoordinatorFactory
}

// NodesCoordinatorWithRaterCreator -
func (r *RunTypeComponentsStub) NodesCoordinatorWithRaterCreator() nodesCoordinator.NodesCoordinatorWithRaterFactory {
	return r.NodesCoordinatorWithRaterFactory
}

// RequestHandlerCreator -
func (r *RunTypeComponentsStub) RequestHandlerCreator() requestHandlers.RequestHandlerCreator {
	return r.RequestHandlerFactory
}

// AccountsCreator -
func (r *RunTypeComponentsStub) AccountsCreator() state.AccountFactory {
	return r.AccountCreator
}

// OutGoingOperationsPoolHandler -
func (r *RunTypeComponentsStub) OutGoingOperationsPoolHandler() sovereignBlock.OutGoingOperationsPool {
	return r.OutGoingOperationsPool
}

// DataCodecHandler -
func (r *RunTypeComponentsStub) DataCodecHandler() sovereign.DataCodecHandler {
	return r.DataCodec
}

// TopicsCheckerHandler -
func (r *RunTypeComponentsStub) TopicsCheckerHandler() sovereign.TopicsCheckerHandler {
	return r.TopicsChecker
}

// RequestersContainerFactoryCreator -
func (r *RunTypeComponentsStub) RequestersContainerFactoryCreator() requesterscontainer.RequesterContainerFactoryCreator {
	return r.RequestersContainerFactoryCreatorField
}

// ValidatorAccountsSyncerFactoryHandler -
func (r *RunTypeComponentsStub) ValidatorAccountsSyncerFactoryHandler() syncerFactory.ValidatorAccountsSyncerFactoryHandler {
	return r.ValidatorAccountsSyncerFactoryHandlerField
}

// ShardRequestersContainerCreatorHandler -
func (r *RunTypeComponentsStub) ShardRequestersContainerCreatorHandler() storageRequestFactory.ShardRequestersContainerCreatorHandler {
	return r.ShardRequestersContainerCreatorHandlerField
}

// IsInterfaceNil -
func (r *RunTypeComponentsStub) IsInterfaceNil() bool {
	return r == nil
}
