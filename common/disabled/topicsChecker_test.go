package disabled

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/stretchr/testify/require"
)

func TestTopicsChecker_MethodsShouldNotPanic(t *testing.T) {
	t.Parallel()

	tc := NewDisabledTopicsChecker()
	require.False(t, check.IfNil(tc))

	require.NotPanics(t, func() {
		err := tc.CheckValidity([][]byte{[]byte("topic")})
		require.NoError(t, err)
	})
}
