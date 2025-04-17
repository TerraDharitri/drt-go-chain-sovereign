package examples

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/blake2b"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/keccak"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/ed25519"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/ed25519/singlesig"
	"github.com/stretchr/testify/require"
)

var (
	addressEncoder, _  = pubkeyConverter.NewBech32PubkeyConverter(32, "drt")
	signingMarshalizer = &marshal.JsonMarshalizer{}
	txSignHasher       = keccak.NewKeccak()
	signer             = &singlesig.Ed25519Signer{}
	signingCryptoSuite = ed25519.NewEd25519()
	contentMarshalizer = &marshal.GogoProtoMarshalizer{}
	contentHasher      = blake2b.NewBlake2b()
)

const alicePrivateKeyHex = "413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9"

func TestConstructTransaction_NoDataNoValue(t *testing.T) {
	tx := &transaction.Transaction{
		Nonce:    89,
		Value:    big.NewInt(0),
		RcvAddr:  getPubkeyOfAddress(t, "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
		SndAddr:  getPubkeyOfAddress(t, "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
		GasPrice: 1000000000,
		GasLimit: 50000,
		ChainID:  []byte("local-testnet"),
		Version:  1,
	}

	tx.Signature = computeTransactionSignature(t, alicePrivateKeyHex, tx)
	require.Equal(t, "6d308fe0924019c84d0c5894507435d4eedea1d3f992df5506daed1f2a2ec27e0c8176067c7a71b1680b3fe661c3b726db58fab4c9be52e169d7d4e78fd42a02", hex.EncodeToString(tx.Signature))

	data, _ := contentMarshalizer.Marshal(tx)
	require.Equal(t, "0859120200001a208049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f82a200139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1388094ebdc0340d08603520d6c6f63616c2d746573746e6574580162406d308fe0924019c84d0c5894507435d4eedea1d3f992df5506daed1f2a2ec27e0c8176067c7a71b1680b3fe661c3b726db58fab4c9be52e169d7d4e78fd42a02", hex.EncodeToString(data))

	txHash := contentHasher.Compute(string(data))
	require.Equal(t, "f04499c49e52297fcc29e35d9f211c64b37a0dcec00d0c7d95f59c9fdf5eef1b", hex.EncodeToString(txHash))
}

func TestConstructTransaction_Usernames(t *testing.T) {
	tx := &transaction.Transaction{
		Nonce:       89,
		Value:       big.NewInt(0),
		RcvAddr:     getPubkeyOfAddress(t, "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
		SndAddr:     getPubkeyOfAddress(t, "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
		GasPrice:    1000000000,
		GasLimit:    50000,
		ChainID:     []byte("local-testnet"),
		Version:     1,
		SndUserName: []byte("alice"),
		RcvUserName: []byte("bob"),
	}

	tx.Signature = computeTransactionSignature(t, alicePrivateKeyHex, tx)
	require.Equal(t, "952250b344541c336667948f2562fe707e2f4b9d00bc3d406937d19145586e427fa29fc947094ca7185d93d7ae41272b9d9a0b723ce4411956e1f5ec29d5f50f", hex.EncodeToString(tx.Signature))

	data, _ := contentMarshalizer.Marshal(tx)
	require.Equal(t, "0859120200001a208049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f82203626f622a200139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e13205616c696365388094ebdc0340d08603520d6c6f63616c2d746573746e657458016240952250b344541c336667948f2562fe707e2f4b9d00bc3d406937d19145586e427fa29fc947094ca7185d93d7ae41272b9d9a0b723ce4411956e1f5ec29d5f50f", hex.EncodeToString(data))

	txHash := contentHasher.Compute(string(data))
	require.Equal(t, "b5a3cf39d1b9eff800214db52dfc9bda69c70a923dc3882f3a0d8e2fabae2573", hex.EncodeToString(txHash))
}

func TestConstructTransaction_WithDataNoValue(t *testing.T) {
	tx := &transaction.Transaction{
		Nonce:    90,
		Value:    big.NewInt(0),
		RcvAddr:  getPubkeyOfAddress(t, "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
		SndAddr:  getPubkeyOfAddress(t, "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
		GasPrice: 1000000000,
		GasLimit: 80000,
		Data:     []byte("hello"),
		ChainID:  []byte("local-testnet"),
		Version:  1,
	}

	tx.Signature = computeTransactionSignature(t, alicePrivateKeyHex, tx)
	require.Equal(t, "75966e81a2fd78ce488c578345f2cb24ce623eb4370ed3ff8330ef9234a322b2b3ad1cc1cfb73c87e7c614bfa687b07c514010584aaa9d0a60a0781f5bdce904", hex.EncodeToString(tx.Signature))

	data, _ := contentMarshalizer.Marshal(tx)
	require.Equal(t, "085a120200001a208049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f82a200139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1388094ebdc034080f1044a0568656c6c6f520d6c6f63616c2d746573746e65745801624075966e81a2fd78ce488c578345f2cb24ce623eb4370ed3ff8330ef9234a322b2b3ad1cc1cfb73c87e7c614bfa687b07c514010584aaa9d0a60a0781f5bdce904", hex.EncodeToString(data))

	txHash := contentHasher.Compute(string(data))
	require.Equal(t, "0096860c0e465a662659208ff69970e06dce73ac3b8d4373b64d13f60253e235", hex.EncodeToString(txHash))
}

func TestConstructTransaction_WithDataWithValue(t *testing.T) {
	tx := &transaction.Transaction{
		Nonce:    91,
		Value:    stringToBigInt("10000000000000000000"),
		RcvAddr:  getPubkeyOfAddress(t, "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
		SndAddr:  getPubkeyOfAddress(t, "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
		GasPrice: 1000000000,
		GasLimit: 100000,
		Data:     []byte("for the book"),
		ChainID:  []byte("local-testnet"),
		Version:  1,
	}

	tx.Signature = computeTransactionSignature(t, alicePrivateKeyHex, tx)
	require.Equal(t, "f0e7f02b7383fb790168d8037b1662b9f6e46f97e1f785f03c0a26718f3aa40886b6c6cadeeb4a036a3673b5569d8e6725a44684f33811187d2217fc114e1809", hex.EncodeToString(tx.Signature))

	data, _ := contentMarshalizer.Marshal(tx)
	require.Equal(t, "085b1209008ac7230489e800001a208049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f82a200139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1388094ebdc0340a08d064a0c666f722074686520626f6f6b520d6c6f63616c2d746573746e657458016240f0e7f02b7383fb790168d8037b1662b9f6e46f97e1f785f03c0a26718f3aa40886b6c6cadeeb4a036a3673b5569d8e6725a44684f33811187d2217fc114e1809", hex.EncodeToString(data))

	txHash := contentHasher.Compute(string(data))
	require.Equal(t, "c72cb860036d62f0d520c5330149342ce176458ca52878a2f1e522f2896c12c2", hex.EncodeToString(txHash))
}

func TestConstructTransaction_WithDataWithLargeValue(t *testing.T) {
	tx := &transaction.Transaction{
		Nonce:    92,
		Value:    stringToBigInt("123456789000000000000000000000"),
		RcvAddr:  getPubkeyOfAddress(t, "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
		SndAddr:  getPubkeyOfAddress(t, "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
		GasPrice: 1000000000,
		GasLimit: 100000,
		Data:     []byte("for the spaceship"),
		ChainID:  []byte("local-testnet"),
		Version:  1,
	}

	tx.Signature = computeTransactionSignature(t, alicePrivateKeyHex, tx)
	require.Equal(t, "40bcc1b98bd65ca833fa71bf2e94fd5307837ea56b4bba49bbed9bae6bc3846dfc2c244e315167b62249ad4d2799e2d8c46fa3e7913e6fb4922447354355ed04", hex.EncodeToString(tx.Signature))

	data, _ := contentMarshalizer.Marshal(tx)
	require.Equal(t, "085c120e00018ee90ff6181f3761632000001a208049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f82a200139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1388094ebdc0340a08d064a11666f722074686520737061636573686970520d6c6f63616c2d746573746e65745801624040bcc1b98bd65ca833fa71bf2e94fd5307837ea56b4bba49bbed9bae6bc3846dfc2c244e315167b62249ad4d2799e2d8c46fa3e7913e6fb4922447354355ed04", hex.EncodeToString(data))

	txHash := contentHasher.Compute(string(data))
	require.Equal(t, "475e2861208cd7683b1621a0cfe69a1af16b083c6787f1627ed4dc3960273f2e", hex.EncodeToString(txHash))
}

func TestConstructTransaction_WithGuardianFields(t *testing.T) {
	tx := &transaction.Transaction{
		Nonce:    92,
		Value:    stringToBigInt("123456789000000000000000000000"),
		RcvAddr:  getPubkeyOfAddress(t, "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
		SndAddr:  getPubkeyOfAddress(t, "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
		GasPrice: 1000000000,
		GasLimit: 150000,
		Data:     []byte("test data field"),
		ChainID:  []byte("local-testnet"),
		Version:  2,
		Options:  2,
	}

	tx.GuardianAddr = getPubkeyOfAddress(t, "drt1x23lzn8483xs2su4fak0r0dqx6w38enpmmqf2yrkylwq7mfnvyhsmueha6")
	tx.GuardianSignature = bytes.Repeat([]byte{0}, 64)

	tx.Signature = computeTransactionSignature(t, alicePrivateKeyHex, tx)
	require.Equal(t, "ac14f089dd19df4c3641bfe7796bb23717fc39eacf18eb915e514fd7fb31ba175c60b93a45d230b53c71b9763edb748ad3ab45972f7d09c69c212c258492c307", hex.EncodeToString(tx.Signature))

	data, _ := contentMarshalizer.Marshal(tx)
	require.Equal(t, "085c120e00018ee90ff6181f3761632000001a208049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f82a200139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1388094ebdc0340f093094a0f746573742064617461206669656c64520d6c6f63616c2d746573746e657458026240ac14f089dd19df4c3641bfe7796bb23717fc39eacf18eb915e514fd7fb31ba175c60b93a45d230b53c71b9763edb748ad3ab45972f7d09c69c212c258492c3076802722032a3f14cf53c4d0543954f6cf1bda0369d13e661dec095107627dc0f6d33612f7a4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", hex.EncodeToString(data))

	txHash := contentHasher.Compute(string(data))
	require.Equal(t, "a0427c60598931b7b3b12f7e546f5f73452a48f0136c3d1c51969a36733dbc3d", hex.EncodeToString(txHash))
}

func TestConstructTransaction_WithNonceZero(t *testing.T) {
	tx := &transaction.Transaction{
		Nonce:    0,
		Value:    big.NewInt(0),
		RcvAddr:  getPubkeyOfAddress(t, "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
		SndAddr:  getPubkeyOfAddress(t, "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"),
		GasPrice: 1000000000,
		GasLimit: 80000,
		Data:     []byte("hello"),
		ChainID:  []byte("local-testnet"),
		Version:  1,
	}

	tx.Signature = computeTransactionSignature(t, alicePrivateKeyHex, tx)
	require.Equal(t, "a46d0601db75691aafd16d14d44aaec73cdb3dcbf80aa72ebfaf8361a143714c851dbba72c3689a8a397f8f6ed6288f48efbd5c5bc6c7a74ae1482f38c4e8e03", hex.EncodeToString(tx.Signature))

	data, _ := contentMarshalizer.Marshal(tx)
	require.Equal(t, "120200001a208049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f82a200139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1388094ebdc034080f1044a0568656c6c6f520d6c6f63616c2d746573746e657458016240a46d0601db75691aafd16d14d44aaec73cdb3dcbf80aa72ebfaf8361a143714c851dbba72c3689a8a397f8f6ed6288f48efbd5c5bc6c7a74ae1482f38c4e8e03", hex.EncodeToString(data))

	txHash := contentHasher.Compute(string(data))
	require.Equal(t, "c77182e6c5c5a2344152d415c913e6e3a5b85e477252e5a4fd10180f2a18ffbe", hex.EncodeToString(txHash))
}

func stringToBigInt(input string) *big.Int {
	result := big.NewInt(0)
	_, _ = result.SetString(input, 10)
	return result
}

func getPubkeyOfAddress(t *testing.T, address string) []byte {
	pubkey, err := addressEncoder.Decode(address)
	require.Nil(t, err)
	return pubkey
}

func computeTransactionSignature(t *testing.T, senderSeedHex string, tx *transaction.Transaction) []byte {
	keyGenerator := signing.NewKeyGenerator(signingCryptoSuite)

	senderSeed, err := hex.DecodeString(senderSeedHex)
	require.Nil(t, err)

	privateKey, err := keyGenerator.PrivateKeyFromByteArray(senderSeed)
	require.Nil(t, err)

	dataToSign, err := tx.GetDataForSigning(addressEncoder, signingMarshalizer, txSignHasher)
	require.Nil(t, err)

	signature, err := signer.Sign(privateKey, dataToSign)
	require.Nil(t, err)
	require.Len(t, signature, 64)

	return signature
}

func TestConstructMiniBlockHeaderReserved_WithMaxValues(t *testing.T) {
	mbhr := &block.MiniBlockHeaderReserved{
		ExecutionType: block.ProcessingType(math.MaxInt32),
		State:         block.MiniBlockState(math.MaxInt32),
	}

	data, _ := contentMarshalizer.Marshal(mbhr)
	fmt.Printf("size: %d\n", len(data))
	require.True(t, len(data) < 16)
}
