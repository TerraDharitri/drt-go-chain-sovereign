package factory

import (
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
)

// BlockChainHookHandlerFactoryMock -
type BlockChainHookHandlerFactoryMock struct {
	CreateBlockChainHookHandlerCalled func(args hooks.ArgBlockChainHook) (process.BlockChainHookWithAccountsAdapter, error)
}

// CreateBlockChainHookHandler -
func (b *BlockChainHookHandlerFactoryMock) CreateBlockChainHookHandler(args hooks.ArgBlockChainHook) (process.BlockChainHookWithAccountsAdapter, error) {
	if b.CreateBlockChainHookHandlerCalled != nil {
		return b.CreateBlockChainHookHandlerCalled(args)
	}
	return &testscommon.BlockChainHookStub{}, nil
}

// IsInterfaceNil -
func (b *BlockChainHookHandlerFactoryMock) IsInterfaceNil() bool {
	return b == nil
}
