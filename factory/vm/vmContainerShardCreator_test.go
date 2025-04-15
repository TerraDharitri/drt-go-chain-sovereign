package vm_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	vmcommonBuiltInFunctions "github.com/TerraDharitri/drt-go-chain-vm-common/builtInFunctions"
	"github.com/TerraDharitri/drt-go-chain-vm-common/parsers"
	wasmConfig "github.com/TerraDharitri/drt-go-chain-vm/config"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/factory/vm"
	"github.com/TerraDharitri/drt-go-chain/process/factory/shard"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/epochNotifier"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
)

func makeVMConfig() config.VirtualMachineConfig {
	return config.VirtualMachineConfig{
		WasmVMVersions: []config.WasmVMVersionByEpoch{
			{StartEpoch: 0, Version: "v1.4"},
			{StartEpoch: 10, Version: "v1.5"},
		},
		TransferAndExecuteByUserAddresses: []string{"3132333435363738393031323334353637383930313233343536373839303234"},
	}
}

func createMockVMAccountsArguments() shard.ArgVMContainerFactory {
	dcdtTransferParser, _ := parsers.NewDCDTTransferParser(&mock.MarshalizerMock{})
	return shard.ArgVMContainerFactory{
		Config:              makeVMConfig(),
		BlockGasLimit:       10000,
		GasSchedule:         testscommon.NewGasScheduleNotifierMock(wasmConfig.MakeGasMapForTests()),
		EpochNotifier:       &epochNotifier.EpochNotifierStub{},
		EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		WasmVMChangeLocker:  &sync.RWMutex{},
		DCDTTransferParser:  dcdtTransferParser,
		BuiltInFunctions:    vmcommonBuiltInFunctions.NewBuiltInFunctionContainer(),
		BlockChainHook:      &testscommon.BlockChainHookStub{},
		Hasher:              &hashingMocks.HasherMock{},
		PubKeyConverter:     &testscommon.PubkeyConverterMock{},
	}
}

func TestNewVmContainerShardCreatorFactory(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		vmContainerShardFactory := vm.NewVmContainerShardFactory()
		require.False(t, vmContainerShardFactory.IsInterfaceNil())
	})
}

func TestNewVmContainerShardFactory_CreateVmContainerFactoryShard(t *testing.T) {
	t.Parallel()
	if runtime.GOARCH == "arm64" {
		t.Skip("skipping test on arm64")
	}

	vmContainerShardFactory := vm.NewVmContainerShardFactory()
	require.False(t, vmContainerShardFactory.IsInterfaceNil())

	argsShard := createMockVMAccountsArguments()
	args := vm.ArgsVmContainerFactory{
		BlockChainHook:      argsShard.BlockChainHook,
		Config:              argsShard.Config,
		BlockGasLimit:       argsShard.BlockGasLimit,
		GasSchedule:         argsShard.GasSchedule,
		EpochNotifier:       argsShard.EpochNotifier,
		EnableEpochsHandler: argsShard.EnableEpochsHandler,
		WasmVMChangeLocker:  argsShard.WasmVMChangeLocker,
		DCDTTransferParser:  argsShard.DCDTTransferParser,
		BuiltInFunctions:    argsShard.BuiltInFunctions,
		Hasher:              argsShard.Hasher,
		PubkeyConv:          argsShard.PubKeyConverter,
	}

	vmContainer, vmFactory, err := vmContainerShardFactory.CreateVmContainerFactory(args)
	require.Nil(t, err)
	require.Equal(t, "*containers.virtualMachinesContainer", fmt.Sprintf("%T", vmContainer))
	require.Equal(t, "*shard.vmContainerFactory", fmt.Sprintf("%T", vmFactory))
}
