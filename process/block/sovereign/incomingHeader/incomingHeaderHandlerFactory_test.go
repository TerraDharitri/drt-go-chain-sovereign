package incomingHeader

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain/config"
	retriever "github.com/TerraDharitri/drt-go-chain/dataRetriever"
	errorsDrt "github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/pool"
	"github.com/stretchr/testify/require"
)

func createWSCfg() config.WebSocketConfig {
	return config.WebSocketConfig{
		MarshallerType: "json",
		HasherType:     "keccak",
	}
}

func TestCreateIncomingHeaderProcessor(t *testing.T) {
	t.Parallel()

	runTypeComps := mock.NewRunTypeComponentsStub()
	headersPool := &dataRetriever.PoolsHolderStub{
		HeadersCalled: func() retriever.HeadersPool {
			return &pool.HeadersPoolStub{}
		},
	}

	t.Run("nil run type comps, should not work", func(t *testing.T) {
		headerProc, err := CreateIncomingHeaderProcessor(
			createWSCfg(),
			headersPool,
			11,
			nil,
		)
		require.Equal(t, errorsDrt.ErrNilRunTypeComponents, err)
		require.Nil(t, headerProc)
	})

	t.Run("invalid marshaller, should not work", func(t *testing.T) {
		cfg := createWSCfg()
		cfg.MarshallerType = ""

		headerProc, err := CreateIncomingHeaderProcessor(
			cfg,
			headersPool,
			11,
			runTypeComps,
		)
		require.NotNil(t, err)
		require.Nil(t, headerProc)
	})

	t.Run("invalid hasher, should not work", func(t *testing.T) {
		cfg := createWSCfg()
		cfg.HasherType = ""

		headerProc, err := CreateIncomingHeaderProcessor(
			cfg,
			headersPool,
			11,
			runTypeComps,
		)
		require.NotNil(t, err)
		require.Nil(t, headerProc)
	})

	t.Run("should work", func(t *testing.T) {
		headerProc, err := CreateIncomingHeaderProcessor(
			createWSCfg(),
			headersPool,
			11,
			runTypeComps,
		)
		require.Nil(t, err)
		require.False(t, headerProc.IsInterfaceNil())
	})
}
