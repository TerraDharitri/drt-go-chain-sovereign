package factory

import (
	"testing"

	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/epochStart/mock"
	processMock "github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/state"
)

func createArgs() ArgsEpochStartInterceptorContainer {
	pubKeyConv := &testscommon.PubkeyConverterMock{}
	return ArgsEpochStartInterceptorContainer{
		CoreComponents: &mock.CoreComponentsMock{
			AddrPubKeyConv: pubKeyConv,
		},
		CryptoComponents:        &mock.CryptoComponentsMock{},
		Config:                  config.Config{},
		ShardCoordinator:        &mock.ShardCoordinatorStub{},
		MainMessenger:           &processMock.TopicHandlerStub{},
		FullArchiveMessenger:    &processMock.TopicHandlerStub{},
		DataPool:                &dataRetriever.PoolsHolderStub{},
		WhiteListHandler:        &testscommon.WhiteListHandlerStub{},
		WhiteListerVerifiedTxs:  &testscommon.WhiteListHandlerStub{},
		AddressPubkeyConv:       pubKeyConv,
		NonceConverter:          &testscommon.Uint64ByteSliceConverterStub{},
		ChainID:                 []byte{1},
		ArgumentsParser:         &testscommon.ArgumentParserMock{},
		HeaderIntegrityVerifier: &processMock.HeaderIntegrityVerifierStub{},
		RequestHandler:          &testscommon.RequestHandlerStub{},
		SignaturesHandler:       &processMock.SignaturesHandlerStub{},
		NodeOperationMode:       "normal",
		AccountFactory:          &state.AccountsFactoryStub{},
	}
}

func TestCreateEpochStartContainerFactoryArgs(t *testing.T) {
	args := createArgs()

	accMock := &state.AccountWrapMock{
		Address: []byte("addr"),
	}
	accFactory := &state.AccountsFactoryStub{
		CreateAccountCalled: func(address []byte) (vmcommon.AccountHandler, error) {
			return accMock, nil
		},
	}

	args.AccountFactory = accFactory
	createdArgs, err := CreateEpochStartContainerFactoryArgs(args)
	require.Nil(t, err)
	require.NotNil(t, createdArgs)

	acc, err := createdArgs.Accounts.LoadAccount([]byte("addr"))
	require.Nil(t, err)
	require.Equal(t, accMock, acc)
}
