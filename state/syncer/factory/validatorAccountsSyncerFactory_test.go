package factory

import (
	"fmt"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain/state/syncer"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/marshallerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/statusHandler"
	"github.com/TerraDharitri/drt-go-chain/testscommon/storageManager"
	"github.com/stretchr/testify/require"
)

func getArgs() syncer.ArgsNewValidatorAccountsSyncer {
	return syncer.ArgsNewValidatorAccountsSyncer{
		ArgsNewBaseAccountsSyncer: syncer.ArgsNewBaseAccountsSyncer{
			Hasher:                            &hashingMocks.HasherMock{},
			Marshalizer:                       marshallerMock.MarshalizerMock{},
			TrieStorageManager:                &storageManager.StorageManagerStub{},
			RequestHandler:                    &testscommon.RequestHandlerStub{},
			Timeout:                           time.Second,
			Cacher:                            testscommon.NewCacherMock(),
			UserAccountsSyncStatisticsHandler: &testscommon.SizeSyncStatisticsHandlerStub{},
			AppStatusHandler:                  &statusHandler.AppStatusHandlerStub{},
			EnableEpochsHandler:               &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
			MaxTrieLevelInMemory:              5,
			MaxHardCapForMissingNodes:         100,
			TrieSyncerVersion:                 3,
			CheckNodesOnDisk:                  false,
		},
	}

}

func TestValidatorAccountsSyncerFactory_CreateValidatorAccountsSyncer(t *testing.T) {
	t.Parallel()

	args := getArgs()
	factory := NewValidatorAccountsSyncerFactory()
	require.False(t, factory.IsInterfaceNil())

	valSyncer, err := factory.CreateValidatorAccountsSyncer(args)
	require.Nil(t, err)
	require.Equal(t, "*syncer.validatorAccountsSyncer", fmt.Sprintf("%T", valSyncer))
}
