package process

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain/errors"

	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/stretchr/testify/require"
)

func TestNewSovereignGenesisMetaBlockChecker(t *testing.T) {
	t.Parallel()

	checker := NewSovereignGenesisMetaBlockChecker()
	require.False(t, checker.IsInterfaceNil())
}

func TestSovereignGenesisMetaBlockChecker_CheckGenesisMetaBlock(t *testing.T) {
	t.Parallel()

	checker := NewSovereignGenesisMetaBlockChecker()
	err := checker.SetValidatorRootHashOnGenesisMetaBlock(nil, nil)
	require.Nil(t, err)

	err = checker.SetValidatorRootHashOnGenesisMetaBlock(&block.MetaBlock{}, nil)
	require.Equal(t, errors.ErrGenesisMetaBlockOnSovereign, err)
}
