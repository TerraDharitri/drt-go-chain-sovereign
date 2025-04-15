package consensus

import "github.com/TerraDharitri/drt-go-chain/process"

func (cc *consensusComponents) BootStrapper() process.Bootstrapper {
	return cc.bootstrapper
}
