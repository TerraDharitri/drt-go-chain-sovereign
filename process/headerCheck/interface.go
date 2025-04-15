package headerCheck

import (
	"github.com/TerraDharitri/drt-go-chain-core/data"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// ExtraHeaderSigVerifierHolder manages extra header verifiers
type ExtraHeaderSigVerifierHolder interface {
	VerifyAggregatedSignature(header data.HeaderHandler, multiSigVerifier crypto.MultiSigner, pubKeysSigners [][]byte) error
	VerifyLeaderSignature(header data.HeaderHandler, leaderPubKey crypto.PublicKey) error
	RemoveLeaderSignature(header data.HeaderHandler) error
	RemoveAllSignatures(header data.HeaderHandler) error
	RegisterExtraHeaderSigVerifier(extraVerifier process.ExtraHeaderSigVerifierHandler) error
	IsInterfaceNil() bool
}
