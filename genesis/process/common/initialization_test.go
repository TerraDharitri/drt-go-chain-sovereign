package common

import (
	"testing"

	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain/testscommon/state"
	"github.com/TerraDharitri/drt-go-chain/vm"
	"github.com/stretchr/testify/require"
)

func TestUpdateSystemSCContractsCode(t *testing.T) {
	contractsToUpdate := map[string]struct{}{
		string(vm.StakingSCAddress):           {},
		string(vm.ValidatorSCAddress):         {},
		string(vm.GovernanceSCAddress):        {},
		string(vm.DCDTSCAddress):              {},
		string(vm.DelegationManagerSCAddress): {},
		string(vm.FirstDelegationSCAddress):   {},
	}

	expectedCodeMetaData := &vmcommon.CodeMetadata{
		Readable:    true,
		PayableBySC: true,
	}

	accountsDB := &state.AccountsStub{
		SaveAccountCalled: func(account vmcommon.AccountHandler) error {
			addrStr := string(account.AddressBytes())
			_, found := contractsToUpdate[addrStr]
			require.True(t, found)

			userAcc := account.(*state.AccountWrapMock)
			require.NotEmpty(t, userAcc.GetCode())
			require.NotEmpty(t, userAcc.GetOwnerAddress())
			require.Equal(t, expectedCodeMetaData.ToBytes(), userAcc.GetCodeMetadata())

			delete(contractsToUpdate, addrStr)

			return nil
		},
		LoadAccountCalled: func(container []byte) (vmcommon.AccountHandler, error) {
			return &state.AccountWrapMock{
				Address: container,
			}, nil
		},
	}

	err := UpdateSystemSCContractsCode(expectedCodeMetaData.ToBytes(), accountsDB)
	require.Nil(t, err)
	require.Empty(t, contractsToUpdate)
}
