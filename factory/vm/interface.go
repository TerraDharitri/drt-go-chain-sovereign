package vm

import (
	"github.com/TerraDharitri/drt-go-chain/process"
)

// VmContainerCreator can create vm container and factory for metachain/shard/sovereign
type VmContainerCreator interface {
	CreateVmContainerFactory(args ArgsVmContainerFactory) (process.VirtualMachinesContainer, process.VirtualMachinesContainerFactory, error)
	IsInterfaceNil() bool
}
