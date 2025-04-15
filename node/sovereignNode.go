package node

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/api"

	"github.com/TerraDharitri/drt-go-chain/errors"
)

const (
	baseDCDTKeyPrefix = core.ProtectedKeyPrefix + core.DCDTKeyIdentifier
)

type sovereignNode struct {
	*Node
	nativeDCDT string
}

// NewSovereignNode creates a new sovereign node instance
func NewSovereignNode(node *Node, nativeDCDT string) (*sovereignNode, error) {
	if check.IfNil(node) {
		return nil, errors.ErrNilNode
	}
	if len(nativeDCDT) == 0 {
		return nil, ErrEmptyNativeDcdt
	}

	return &sovereignNode{
		Node:       node,
		nativeDCDT: nativeDCDT,
	}, nil
}

// GetAllIssuedDCDTs returns all the issued dcdt tokens, works only on metachain
func (sn *sovereignNode) GetAllIssuedDCDTs(tokenType string, ctx context.Context) ([]string, error) {
	tokens, err := sn.baseGetAllIssuedDCDTs(tokenType, ctx)
	if err != nil {
		return make([]string, 0), err
	}

	return sn.getTokensWithoutNativeDCDT(tokens), nil
}

func (sn *sovereignNode) getTokensWithoutNativeDCDT(tokens []string) []string {
	nativeDcdtWithBasePrefix := baseDCDTKeyPrefix + sn.nativeDCDT
	issuedTokens := make([]string, 0)
	for _, token := range tokens {
		if token != nativeDcdtWithBasePrefix {
			issuedTokens = append(issuedTokens, token)
		}
	}
	return issuedTokens
}

// GetNFTTokenIDsRegisteredByAddress returns all the token identifiers for semi or non fungible tokens registered by the address
func (sn *sovereignNode) GetNFTTokenIDsRegisteredByAddress(address string, options api.AccountQueryOptions, ctx context.Context) ([]string, api.BlockInfo, error) {
	addressBytes, err := sn.coreComponents.AddressPubKeyConverter().Decode(address)
	if err != nil {
		return nil, api.BlockInfo{}, err
	}

	f := &getRegisteredNftsFilter{
		addressBytes: addressBytes,
	}
	return sn.baseGetTokensIDsWithFilter(f, options, ctx)
}

// GetDCDTsWithRole returns all the tokens with the given role for the given address
func (sn *sovereignNode) GetDCDTsWithRole(address string, role string, options api.AccountQueryOptions, ctx context.Context) ([]string, api.BlockInfo, error) {
	if !core.IsValidDCDTRole(role) {
		return nil, api.BlockInfo{}, ErrInvalidDCDTRole
	}

	addressBytes, err := sn.coreComponents.AddressPubKeyConverter().Decode(address)
	if err != nil {
		return nil, api.BlockInfo{}, err
	}

	f := &getTokensWithRoleFilter{
		addressBytes: addressBytes,
		role:         role,
	}
	return sn.baseGetTokensIDsWithFilter(f, options, ctx)
}

// GetDCDTsRoles returns all the tokens identifiers and roles for the given address
func (sn *sovereignNode) GetDCDTsRoles(address string, options api.AccountQueryOptions, ctx context.Context) (map[string][]string, api.BlockInfo, error) {
	addressBytes, err := sn.coreComponents.AddressPubKeyConverter().Decode(address)
	if err != nil {
		return nil, api.BlockInfo{}, err
	}

	tokensRoles := make(map[string][]string)

	f := &getAllTokensRolesFilter{
		addressBytes: addressBytes,
		outputRoles:  tokensRoles,
	}
	_, blockInfo, err := sn.baseGetTokensIDsWithFilter(f, options, ctx)
	if err != nil {
		return nil, api.BlockInfo{}, err
	}

	return tokensRoles, blockInfo, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sn *sovereignNode) IsInterfaceNil() bool {
	return sn == nil
}
