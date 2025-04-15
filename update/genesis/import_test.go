package genesis

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon/storageManager"
	"github.com/TerraDharitri/drt-go-chain/update"
	"github.com/TerraDharitri/drt-go-chain/update/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//TODO increase code coverage

func createArgsNewStateImport() ArgsNewStateImport {
	trieStorageManagers := make(map[string]common.StorageManager)
	trieStorageManagers[dataRetriever.UserAccountsUnit.String()] = &storageManager.StorageManagerStub{}
	trieStorageManagers[dataRetriever.PeerAccountsUnit.String()] = &storageManager.StorageManagerStub{}
	return ArgsNewStateImport{
		HardforkStorer:      &mock.HardforkStorerStub{},
		Marshalizer:         &mock.MarshalizerMock{},
		Hasher:              &mock.HasherStub{},
		TrieStorageManagers: trieStorageManagers,
		AddressConverter:    &testscommon.PubkeyConverterMock{},
		EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		AccountCreator:      &stateMock.AccountsFactoryStub{},
	}
}

func TestNewStateImport(t *testing.T) {

	tests := []struct {
		name    string
		args    func() ArgsNewStateImport
		exError error
	}{
		{
			name: "NilHarforkStorer",
			args: func() ArgsNewStateImport {
				args := createArgsNewStateImport()
				args.HardforkStorer = nil
				return args
			},
			exError: update.ErrNilHardforkStorer,
		},
		{
			name: "NilMarshalizer",
			args: func() ArgsNewStateImport {
				args := createArgsNewStateImport()
				args.Marshalizer = nil
				return args
			},
			exError: update.ErrNilMarshalizer,
		},
		{
			name: "NilHasher",
			args: func() ArgsNewStateImport {
				args := createArgsNewStateImport()
				args.Hasher = nil
				return args
			},
			exError: update.ErrNilHasher,
		},
		{
			name: "NilAccountCreator",
			args: func() ArgsNewStateImport {
				args := createArgsNewStateImport()
				args.AccountCreator = nil
				return args
			},
			exError: state.ErrNilAccountFactory,
		},
		{
			name: "Ok",
			args: func() ArgsNewStateImport {
				return createArgsNewStateImport()
			},
			exError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStateImport(tt.args())
			require.Equal(t, tt.exError, err)
		})
	}
}

func TestImportAll(t *testing.T) {
	t.Parallel()

	args := createArgsNewStateImport()
	importState, _ := NewStateImport(args)
	require.False(t, check.IfNil(importState))

	err := importState.ImportAll()
	require.Nil(t, err)
}

func TestStateImport_ImportUnFinishedMetaBlocksShouldWork(t *testing.T) {
	t.Parallel()

	args := createArgsNewStateImport()
	metaBlock := &block.MetaBlock{
		Round:   1,
		ChainID: []byte("chainId"),
	}
	metaBlockHash, _ := core.CalculateHash(args.Marshalizer, args.Hasher, metaBlock)

	args.HardforkStorer = &mock.HardforkStorerStub{
		GetCalled: func(identifier string, key []byte) ([]byte, error) {
			return args.Marshalizer.Marshal(metaBlock)
		},
	}
	importState, _ := NewStateImport(args)
	require.False(t, check.IfNil(importState))

	key := fmt.Sprintf("meta@chainId@%s", hex.EncodeToString(metaBlockHash))
	err := importState.importUnFinishedMetaBlocks(UnFinishedMetaBlocksIdentifier, [][]byte{
		[]byte(key),
	})

	require.Nil(t, err)
	assert.Equal(t, importState.importedUnFinishedMetaBlocks[string(metaBlockHash)], metaBlock)
}
