package process

import (
	"github.com/TerraDharitri/drt-go-chain/errors"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data"
)

type sovereignGenesisMetaBlockChecker struct {
}

// NewSovereignGenesisMetaBlockChecker creates a sovereign meta block genesis checker
func NewSovereignGenesisMetaBlockChecker() *sovereignGenesisMetaBlockChecker {
	return &sovereignGenesisMetaBlockChecker{}
}

// SetValidatorRootHashOnGenesisMetaBlock checks that the genesis meta block does not exist
func (gmbc *sovereignGenesisMetaBlockChecker) SetValidatorRootHashOnGenesisMetaBlock(genesisMetaBlock data.HeaderHandler, _ []byte) error {
	if !check.IfNil(genesisMetaBlock) {
		return errors.ErrGenesisMetaBlockOnSovereign
	}

	return nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (gmbc *sovereignGenesisMetaBlockChecker) IsInterfaceNil() bool {
	return gmbc == nil
}
