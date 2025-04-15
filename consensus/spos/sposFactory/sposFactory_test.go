package sposFactory_test

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/consensus"
	"github.com/TerraDharitri/drt-go-chain/consensus/mock"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/bls"
	"github.com/TerraDharitri/drt-go-chain/consensus/spos/sposFactory"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/outport"
	"github.com/TerraDharitri/drt-go-chain/testscommon/p2pmocks"
	statusHandlerMock "github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
	"github.com/TerraDharitri/drt-go-chain/testscommon/subRoundsHolder"
	"github.com/stretchr/testify/assert"
)

var currentPid = core.PeerID("pid")

func TestGetConsensusCoreFactory_InvalidTypeShouldErr(t *testing.T) {
	t.Parallel()

	csf, err := sposFactory.GetConsensusCoreFactory("invalid")

	assert.Nil(t, csf)
	assert.Equal(t, sposFactory.ErrInvalidConsensusType, err)
}

func TestGetConsensusCoreFactory_BlsShouldWork(t *testing.T) {
	t.Parallel()

	csf, err := sposFactory.GetConsensusCoreFactory(consensus.BlsConsensusType)

	assert.Nil(t, err)
	assert.False(t, check.IfNil(csf))
}

func TestGetSubroundsFactory_BlsNilConsensusCoreShouldErr(t *testing.T) {
	t.Parallel()

	worker := &mock.SposWorkerMock{}
	consensusType := consensus.BlsConsensusType
	statusHandler := statusHandlerMock.NewAppStatusHandlerMock()
	chainID := []byte("chain-id")
	indexer := &outport.OutportStub{}
	sf, err := sposFactory.GetSubroundsFactory(
		nil,
		&spos.ConsensusState{},
		worker,
		consensusType,
		statusHandler,
		indexer,
		&testscommon.SentSignatureTrackerStub{},
		chainID,
		currentPid,
		consensus.ConsensusModelV1,
		&enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		&subRoundsHolder.ExtraSignersHolderMock{},
		bls.NewSubRoundEndV2Creator(),
	)

	assert.Nil(t, sf)
	assert.Equal(t, spos.ErrNilConsensusCore, err)
}

func TestGetSubroundsFactory_BlsNilStatusHandlerShouldErr(t *testing.T) {
	t.Parallel()

	consensusCore := mock.InitConsensusCore()
	worker := &mock.SposWorkerMock{}
	consensusType := consensus.BlsConsensusType
	chainID := []byte("chain-id")
	indexer := &outport.OutportStub{}
	sf, err := sposFactory.GetSubroundsFactory(
		consensusCore,
		&spos.ConsensusState{},
		worker,
		consensusType,
		nil,
		indexer,
		&testscommon.SentSignatureTrackerStub{},
		chainID,
		currentPid,
		consensus.ConsensusModelV1,
		&enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		&subRoundsHolder.ExtraSignersHolderMock{},
		bls.NewSubRoundEndV2Creator(),
	)

	assert.Nil(t, sf)
	assert.Equal(t, spos.ErrNilAppStatusHandler, err)
}

func TestGetSubroundsFactory_BlsShouldWork(t *testing.T) {
	t.Parallel()

	consensusCore := mock.InitConsensusCore()
	worker := &mock.SposWorkerMock{}
	consensusType := consensus.BlsConsensusType
	statusHandler := statusHandlerMock.NewAppStatusHandlerMock()
	chainID := []byte("chain-id")
	indexer := &outport.OutportStub{}
	sf, err := sposFactory.GetSubroundsFactory(
		consensusCore,
		&spos.ConsensusState{},
		worker,
		consensusType,
		statusHandler,
		indexer,
		&testscommon.SentSignatureTrackerStub{},
		chainID,
		currentPid,
		consensus.ConsensusModelV1,
		&enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		&subRoundsHolder.ExtraSignersHolderMock{},
		bls.NewSubRoundEndV2Creator(),
	)
	assert.Nil(t, err)
	assert.False(t, check.IfNil(sf))
}

func TestGetSubroundsFactory_InvalidConsensusTypeShouldErr(t *testing.T) {
	t.Parallel()

	consensusType := "invalid"
	sf, err := sposFactory.GetSubroundsFactory(
		nil,
		nil,
		nil,
		consensusType,
		nil,
		nil,
		nil,
		nil,
		currentPid,
		consensus.ConsensusModelV1,
		&enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		&subRoundsHolder.ExtraSignersHolderMock{},
		bls.NewSubRoundEndV2Creator(),
	)

	assert.Nil(t, sf)
	assert.Equal(t, sposFactory.ErrInvalidConsensusType, err)
}

func TestGetBroadcastMessenger_ShardShouldWork(t *testing.T) {
	t.Parallel()

	marshalizer := &mock.MarshalizerMock{}
	hasher := &hashingMocks.HasherMock{}
	messenger := &p2pmocks.MessengerStub{}
	shardCoord := mock.NewMultiShardsCoordinatorMock(3)
	shardCoord.SelfIDCalled = func() uint32 {
		return 0
	}
	peerSigHandler := &mock.PeerSignatureHandler{}
	headersSubscriber := &testscommon.HeadersCacherStub{}
	interceptosContainer := &testscommon.InterceptorsContainerStub{}
	alarmSchedulerStub := &mock.AlarmSchedulerStub{}

	bm, err := sposFactory.GetBroadcastMessenger(
		marshalizer,
		hasher,
		messenger,
		shardCoord,
		peerSigHandler,
		headersSubscriber,
		interceptosContainer,
		alarmSchedulerStub,
		&testscommon.KeysHandlerStub{},
		&factory.ShardChainMessengerFactoryMock{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, bm)
}

func TestGetBroadcastMessenger_MetachainShouldWork(t *testing.T) {
	t.Parallel()

	marshalizer := &mock.MarshalizerMock{}
	hasher := &hashingMocks.HasherMock{}
	messenger := &p2pmocks.MessengerStub{}
	shardCoord := mock.NewMultiShardsCoordinatorMock(3)
	shardCoord.SelfIDCalled = func() uint32 {
		return core.MetachainShardId
	}
	peerSigHandler := &mock.PeerSignatureHandler{}
	headersSubscriber := &testscommon.HeadersCacherStub{}
	interceptosContainer := &testscommon.InterceptorsContainerStub{}
	alarmSchedulerStub := &mock.AlarmSchedulerStub{}

	bm, err := sposFactory.GetBroadcastMessenger(
		marshalizer,
		hasher,
		messenger,
		shardCoord,
		peerSigHandler,
		headersSubscriber,
		interceptosContainer,
		alarmSchedulerStub,
		&testscommon.KeysHandlerStub{},
		&factory.ShardChainMessengerFactoryMock{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, bm)
}

func TestGetBroadcastMessenger_NilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	headersSubscriber := &testscommon.HeadersCacherStub{}
	interceptosContainer := &testscommon.InterceptorsContainerStub{}
	alarmSchedulerStub := &mock.AlarmSchedulerStub{}

	bm, err := sposFactory.GetBroadcastMessenger(
		nil,
		nil,
		nil,
		nil,
		nil,
		headersSubscriber,
		interceptosContainer,
		alarmSchedulerStub,
		&testscommon.KeysHandlerStub{},
		&factory.ShardChainMessengerFactoryMock{},
	)

	assert.Nil(t, bm)
	assert.Equal(t, spos.ErrNilShardCoordinator, err)
}

func TestGetBroadcastMessenger_InvalidShardIdShouldErr(t *testing.T) {
	t.Parallel()

	shardCoord := mock.NewMultiShardsCoordinatorMock(3)
	shardCoord.SelfIDCalled = func() uint32 {
		return 37
	}
	headersSubscriber := &testscommon.HeadersCacherStub{}
	interceptosContainer := &testscommon.InterceptorsContainerStub{}
	alarmSchedulerStub := &mock.AlarmSchedulerStub{}

	bm, err := sposFactory.GetBroadcastMessenger(
		nil,
		nil,
		nil,
		shardCoord,
		nil,
		headersSubscriber,
		interceptosContainer,
		alarmSchedulerStub,
		&testscommon.KeysHandlerStub{},
		&factory.ShardChainMessengerFactoryMock{},
	)

	assert.Nil(t, bm)
	assert.Equal(t, sposFactory.ErrInvalidShardId, err)
}
