package genesisMocks

import (
	nodeFactory "github.com/TerraDharitri/drt-go-chain/cmd/node/factory"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/factory/block"
	factoryVm "github.com/TerraDharitri/drt-go-chain/factory/vm"
	"github.com/TerraDharitri/drt-go-chain/genesis"
	"github.com/TerraDharitri/drt-go-chain/process/block/preprocess"
	"github.com/TerraDharitri/drt-go-chain/process/coordinator"
	"github.com/TerraDharitri/drt-go-chain/process/factory/shard"
	"github.com/TerraDharitri/drt-go-chain/process/factory/shard/data"
	"github.com/TerraDharitri/drt-go-chain/process/factory/sovereign"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/processProxy"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/processorV2"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/scrCommon"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/factory"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/marshallerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/vmContext"
	"github.com/TerraDharitri/drt-go-chain/vm/systemSmartContracts"
)

// RunTypeComponentsStub -
type RunTypeComponentsStub struct {
	BlockChainHookHandlerFactory              hooks.BlockChainHookHandlerCreator
	TransactionCoordinatorFactory             coordinator.TransactionCoordinatorCreator
	SCResultsPreProcessorFactory              preprocess.SmartContractResultPreProcessorCreator
	SCProcessorFactory                        scrCommon.SCProcessorCreator
	AccountParser                             genesis.AccountsParser
	AccountCreator                            state.AccountFactory
	VMContextCreatorHandler                   systemSmartContracts.VMContextCreatorHandler
	ShardCoordinatorFactory                   sharding.ShardCoordinatorFactory
	TxPreProcessorFactory                     preprocess.TxPreProcessorCreator
	VmContainerShardFactory                   factoryVm.VmContainerCreator
	VmContainerMetaFactory                    factoryVm.VmContainerCreator
	PreProcessorsContainerFactoryCreatorField data.PreProcessorsContainerFactoryCreator
	VersionedHeaderFactoryField               genesis.VersionedHeaderFactory
}

// NewRunTypeComponentsStub -
func NewRunTypeComponentsStub() *RunTypeComponentsStub {
	accountsCreator, _ := factory.NewAccountCreator(factory.ArgsAccountCreator{
		Hasher:              &hashingMocks.HasherMock{},
		Marshaller:          &marshallerMock.MarshalizerMock{},
		EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	})
	vmContainerMeta, _ := factoryVm.NewVmContainerMetaFactory(systemSmartContracts.NewVMContextCreator())
	hdrFactory, _ := block.NewShardHeaderFactory(createHeaderVersionHandler("*"))

	return &RunTypeComponentsStub{
		BlockChainHookHandlerFactory:              hooks.NewBlockChainHookFactory(),
		TransactionCoordinatorFactory:             coordinator.NewShardTransactionCoordinatorFactory(),
		SCResultsPreProcessorFactory:              preprocess.NewSmartContractResultPreProcessorFactory(),
		SCProcessorFactory:                        processProxy.NewSCProcessProxyFactory(),
		AccountParser:                             &AccountsParserStub{},
		AccountCreator:                            accountsCreator,
		VMContextCreatorHandler:                   systemSmartContracts.NewVMContextCreator(),
		ShardCoordinatorFactory:                   sharding.NewMultiShardCoordinatorFactory(),
		TxPreProcessorFactory:                     preprocess.NewTxPreProcessorCreator(),
		VmContainerShardFactory:                   factoryVm.NewVmContainerShardFactory(),
		VmContainerMetaFactory:                    vmContainerMeta,
		PreProcessorsContainerFactoryCreatorField: shard.NewPreProcessorContainerFactoryCreator(),
		VersionedHeaderFactoryField:               hdrFactory,
	}
}

// NewSovereignRunTypeComponentsStub -
func NewSovereignRunTypeComponentsStub() *RunTypeComponentsStub {
	runTypeComponents := NewRunTypeComponentsStub()
	transactionCoordinatorFactory, _ := coordinator.NewSovereignTransactionCoordinatorFactory(runTypeComponents.TransactionCoordinatorFactory)
	scResultsPreProcessorCreator, _ := preprocess.NewSovereignSmartContractResultPreProcessorFactory(runTypeComponents.SCResultsPreProcessorFactory)
	scProcessorCreator, _ := processorV2.NewSovereignSCProcessFactory(runTypeComponents.SCProcessorFactory)
	accountsCreator, _ := factory.NewSovereignAccountCreator(factory.ArgsSovereignAccountCreator{
		ArgsAccountCreator: factory.ArgsAccountCreator{
			Hasher:              &hashingMocks.HasherMock{},
			Marshaller:          &marshallerMock.MarshalizerMock{},
			EnableEpochsHandler: &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
		},
		BaseTokenID: "WREWA-bd4d79",
	})

	oneShardVM := systemSmartContracts.NewOneShardSystemVMEEICreator()
	vmMetaFactory, _ := factoryVm.NewVmContainerMetaFactory(oneShardVM)
	sovVMContainerShardFactory, _ := factoryVm.NewSovereignVmContainerShardFactory(vmMetaFactory, runTypeComponents.VmContainerShardFactory)
	sovVMContainerMeta, _ := factoryVm.NewVmContainerMetaFactory(oneShardVM)
	sovHdrFactory, _ := block.NewSovereignShardHeaderFactory(createHeaderVersionHandler("S1"))

	return &RunTypeComponentsStub{
		BlockChainHookHandlerFactory:              hooks.NewSovereignBlockChainHookFactory(),
		TransactionCoordinatorFactory:             transactionCoordinatorFactory,
		SCResultsPreProcessorFactory:              scResultsPreProcessorCreator,
		SCProcessorFactory:                        scProcessorCreator,
		AccountParser:                             &AccountsParserStub{},
		AccountCreator:                            accountsCreator,
		VMContextCreatorHandler:                   &vmContext.VMContextCreatorStub{},
		ShardCoordinatorFactory:                   sharding.NewSovereignShardCoordinatorFactory(),
		TxPreProcessorFactory:                     preprocess.NewSovereignTxPreProcessorCreator(),
		VmContainerShardFactory:                   sovVMContainerShardFactory,
		VmContainerMetaFactory:                    sovVMContainerMeta,
		PreProcessorsContainerFactoryCreatorField: sovereign.NewSovereignPreProcessorContainerFactoryCreator(),
		VersionedHeaderFactoryField:               sovHdrFactory,
	}
}

func createHeaderVersionHandler(version string) nodeFactory.HeaderVersionHandler {
	hdrVersionHandler, _ := block.NewHeaderVersionHandler(
		[]config.VersionByEpochs{
			{
				Version: version,
			},
		},
		version,
		&testscommon.CacherStub{},
	)
	return hdrVersionHandler
}

// BlockChainHookHandlerCreator -
func (r *RunTypeComponentsStub) BlockChainHookHandlerCreator() hooks.BlockChainHookHandlerCreator {
	return r.BlockChainHookHandlerFactory
}

// TransactionCoordinatorCreator -
func (r *RunTypeComponentsStub) TransactionCoordinatorCreator() coordinator.TransactionCoordinatorCreator {
	return r.TransactionCoordinatorFactory
}

// SCProcessorCreator -
func (r *RunTypeComponentsStub) SCProcessorCreator() scrCommon.SCProcessorCreator {
	return r.SCProcessorFactory
}

// SCResultsPreProcessorCreator -
func (r *RunTypeComponentsStub) SCResultsPreProcessorCreator() preprocess.SmartContractResultPreProcessorCreator {
	return r.SCResultsPreProcessorFactory
}

// AccountsParser -
func (r *RunTypeComponentsStub) AccountsParser() genesis.AccountsParser {
	return r.AccountParser
}

// AccountsCreator -
func (r *RunTypeComponentsStub) AccountsCreator() state.AccountFactory {
	return r.AccountCreator
}

// VMContextCreator -
func (r *RunTypeComponentsStub) VMContextCreator() systemSmartContracts.VMContextCreatorHandler {
	return r.VMContextCreatorHandler
}

// ShardCoordinatorCreator -
func (r *RunTypeComponentsStub) ShardCoordinatorCreator() sharding.ShardCoordinatorFactory {
	return r.ShardCoordinatorFactory
}

// TxPreProcessorCreator -
func (r *RunTypeComponentsStub) TxPreProcessorCreator() preprocess.TxPreProcessorCreator {
	return r.TxPreProcessorFactory
}

// VmContainerShardFactoryCreator -
func (r *RunTypeComponentsStub) VmContainerShardFactoryCreator() factoryVm.VmContainerCreator {
	return r.VmContainerShardFactory
}

// VmContainerMetaFactoryCreator -
func (r *RunTypeComponentsStub) VmContainerMetaFactoryCreator() factoryVm.VmContainerCreator {
	return r.VmContainerMetaFactory
}

// PreProcessorsContainerFactoryCreator -
func (r *RunTypeComponentsStub) PreProcessorsContainerFactoryCreator() data.PreProcessorsContainerFactoryCreator {
	return r.PreProcessorsContainerFactoryCreatorField
}

// VersionedHeaderFactory  -
func (r *RunTypeComponentsStub) VersionedHeaderFactory() genesis.VersionedHeaderFactory {
	return r.VersionedHeaderFactoryField
}

// IsInterfaceNil -
func (r *RunTypeComponentsStub) IsInterfaceNil() bool {
	return r == nil
}
