package wallet

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mintingValue = "100000000"

func TestInterceptedTxWithoutDataField(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	value := big.NewInt(0)
	value.SetString("999", 10)

	testInterceptedTxFromFrontendGeneratedParams(
		t,
		0,
		value,
		"drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c",
		"drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf",
		"a7f77c28e235ca5eb89ec6c8b0e3935648e2f907878d1f75df582c0db86a121399392c7d56f1ffefdd62bf30c0a796bca5e6c913034d183cb8e777e6e84b9f02",
		10,
		100000,
		[]byte(""),
		integrationTests.ChainID,
		integrationTests.MinTransactionVersion,
		0,
	)
}

func TestInterceptedTxWithDataField(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	value := big.NewInt(0)
	value.SetString("999", 10)

	testInterceptedTxFromFrontendGeneratedParams(
		t,
		0,
		value,
		"drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c",
		"drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf",
		"456f8de5d3ccc0701734eec7f8246d6321efd0b5dd76f69ea2e1ae4385f697924121c956f1531b308ce33c813a39e962e995d10f05b5e0cddb49a2f960da7c0a",
		10,
		100000,
		[]byte("data@~`!@#$^&*()_=[]{};'<>?,./|<>><!!!!!"),
		integrationTests.ChainID,
		integrationTests.MinTransactionVersion,
		0,
	)
}

func TestInterceptedTxWithSigningOverTxHash(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	value := big.NewInt(0)
	value.SetString("1000000000000000000", 10)

	testInterceptedTxFromFrontendGeneratedParams(
		t,
		1,
		value,
		"drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c",
		"drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf",
		"e3ab51ebcf02892433500aeb4e91476b77204913a4d91e14207ec9c0827b5a9a18f8d83452ee7d2cfa2112503752dd0bffcd76fa36987e75ec882db918ad050b",
		1000000000,
		56000,
		[]byte("test"),
		integrationTests.ChainID,
		2,
		0,
	)
}

// testInterceptedTxFromFrontendGeneratedParams tests that a frontend generated tx will pass through an interceptor
// and ends up in the datapool, concluding the tx is correctly signed and follows our protocol
func testInterceptedTxFromFrontendGeneratedParams(
	t *testing.T,
	frontendNonce uint64,
	frontendValue *big.Int,
	frontendReceiver string,
	frontendSender string,
	frontendSignatureHex string,
	frontendGasPrice uint64,
	frontendGasLimit uint64,
	frontendData []byte,
	chainID []byte,
	version uint32,
	options uint32,
) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	chDone := make(chan struct{})

	maxShards := uint32(1)
	nodeShardId := uint32(0)
	txSignPrivKeyShardId := uint32(0)
	valMinting, _ := big.NewInt(0).SetString(mintingValue, 10)
	valMinting.Mul(valMinting, big.NewInt(5000000))

	node := integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
		MaxShards:            maxShards,
		NodeShardId:          nodeShardId,
		TxSignPrivKeyShardId: txSignPrivKeyShardId,
	})

	node.EconomicsData.SetMinGasPrice(10)
	txHexHash := ""

	err := node.SetAccountNonce(uint64(0))
	assert.Nil(t, err)

	node.DataPool.Transactions().RegisterOnAdded(func(key []byte, value interface{}) {
		assert.Equal(t, txHexHash, hex.EncodeToString(key))

		dataRecovered, _ := node.DataPool.Transactions().SearchFirstData(key)
		assert.NotNil(t, dataRecovered)

		txRecovered, ok := dataRecovered.(*transaction.Transaction)
		assert.True(t, ok)

		assert.Equal(t, frontendNonce, txRecovered.Nonce)
		assert.Equal(t, frontendValue, txRecovered.Value)

		sender, _ := integrationTests.TestAddressPubkeyConverter.Decode(frontendSender)
		assert.Equal(t, sender, txRecovered.SndAddr)

		receiver, _ := integrationTests.TestAddressPubkeyConverter.Decode(frontendReceiver)
		assert.Equal(t, receiver, txRecovered.RcvAddr)

		sig, _ := hex.DecodeString(frontendSignatureHex)
		assert.Equal(t, sig, txRecovered.Signature)
		assert.Equal(t, len(frontendData), len(txRecovered.Data))

		chDone <- struct{}{}
	})

	rcvAddrBytes, _ := integrationTests.TestAddressPubkeyConverter.Decode(frontendReceiver)
	sndAddrBytes, _ := integrationTests.TestAddressPubkeyConverter.Decode(frontendSender)
	signatureBytes, _ := hex.DecodeString(frontendSignatureHex)

	integrationTests.MintAddress(node.AccntState, sndAddrBytes, valMinting)

	tx := &transaction.Transaction{
		Nonce:     frontendNonce,
		RcvAddr:   rcvAddrBytes,
		SndAddr:   sndAddrBytes,
		GasPrice:  frontendGasPrice,
		GasLimit:  frontendGasLimit,
		Data:      frontendData,
		Signature: signatureBytes,
		ChainID:   chainID,
		Version:   version,
		Options:   options,
	}
	tx.Value = big.NewInt(0).Set(frontendValue)
	txHexHash, err = node.SendTransaction(tx)
	require.Nil(t, err)

	select {
	case <-chDone:
	case <-time.After(time.Second * 3):
		assert.Fail(t, "timeout getting transaction")
	}
}
