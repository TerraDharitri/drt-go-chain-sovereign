package testscommon

import (
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-chain/testscommon/shardingMocks"
)

// NodesCoordinatorFactoryMock -
type NodesCoordinatorFactoryMock struct {
	CreateNodesCoordinatorWithRaterCalled func(args *nodesCoordinator.NodesCoordinatorWithRaterArgs) (nodesCoordinator.NodesCoordinator, error)
}

// CreateNodesCoordinatorWithRater -
func (n *NodesCoordinatorFactoryMock) CreateNodesCoordinatorWithRater(args *nodesCoordinator.NodesCoordinatorWithRaterArgs) (nodesCoordinator.NodesCoordinator, error) {
	if n.CreateNodesCoordinatorWithRaterCalled != nil {
		return n.CreateNodesCoordinatorWithRaterCalled(args)
	}
	return &shardingMocks.NodesCoordinatorMock{}, nil
}

// IsInterfaceNil -
func (n *NodesCoordinatorFactoryMock) IsInterfaceNil() bool {
	return n == nil
}
