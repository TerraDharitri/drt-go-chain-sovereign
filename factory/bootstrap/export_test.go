package bootstrap

import "github.com/TerraDharitri/drt-go-chain/factory"

func (bc *bootstrapComponents) EpochStartBootstrapper() factory.EpochStartBootstrapper {
	return bc.epochStartBootstrapper
}
