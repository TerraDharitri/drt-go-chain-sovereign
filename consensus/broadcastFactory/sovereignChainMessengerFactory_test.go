package broadcastFactory

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain/consensus/broadcast"
	"github.com/TerraDharitri/drt-go-chain/consensus/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	"github.com/stretchr/testify/require"
)

func createDefaultShardChainArgs() broadcast.ShardChainMessengerArgs {
	return broadcast.ShardChainMessengerArgs{
		CommonMessengerArgs: broadcast.CommonMessengerArgs{
			Marshalizer:                &mock.MarshalizerMock{},
			Hasher:                     &hashingMocks.HasherMock{},
			Messenger:                  &p2pmocks.MessengerStub{},
			ShardCoordinator:           &mock.ShardCoordinatorMock{},
			PeerSignatureHandler:       &mock.PeerSignatureHandler{},
			HeadersSubscriber:          &testscommon.HeadersCacherStub{},
			InterceptorsContainer:      &testscommon.InterceptorsContainerStub{},
			MaxDelayCacheSize:          1,
			MaxValidatorDelayCacheSize: 1,
			AlarmScheduler:             &mock.AlarmSchedulerStub{},
			KeysHandler:                &testscommon.KeysHandlerStub{},
		},
	}
}

func TestSovereignChainMessengerFactory_CreateShardChainMessenger(t *testing.T) {
	t.Parallel()

	f := NewSovereignShardChainMessengerFactory()
	require.False(t, f.IsInterfaceNil())

	args := createDefaultShardChainArgs()
	msg, err := f.CreateShardChainMessenger(args)
	require.Nil(t, err)
	require.NotNil(t, msg)
	require.Equal(t, "*broadcast.sovereignChainMessenger", fmt.Sprintf("%T", msg))
}
