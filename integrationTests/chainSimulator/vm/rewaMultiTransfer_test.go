package vm

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/txsFee"
	"github.com/TerraDharitri/drt-go-chain/integrationTests/vm/txsFee/utils"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/components/api"
	"github.com/TerraDharitri/drt-go-chain/node/chainSimulator/configs"
	"github.com/TerraDharitri/drt-go-chain/vm"
	"github.com/stretchr/testify/require"
)

func TestChainSimulator_REWA_MultiTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	startTime := time.Now().Unix()
	roundDurationInMillis := uint64(6000)
	roundsPerEpoch := core.OptionalUint64{
		HasValue: true,
		Value:    20,
	}

	activationEpoch := uint32(4)

	baseIssuingCost := "1000"

	numOfShards := uint32(3)
	cs, err := chainSimulator.NewChainSimulator(chainSimulator.ArgsChainSimulator{
		BypassTxSignatureCheck:   true,
		TempDir:                  t.TempDir(),
		PathToInitialConfig:      defaultPathToInitialConfig,
		NumOfShards:              numOfShards,
		GenesisTimestamp:         startTime,
		RoundDurationInMillis:    roundDurationInMillis,
		RoundsPerEpoch:           roundsPerEpoch,
		ApiInterface:             api.NewNoApiInterface(),
		MinNodesPerShard:         3,
		MetaChainMinNodes:        3,
		NumNodesWaitingListMeta:  0,
		NumNodesWaitingListShard: 0,
		AlterConfigsFunction: func(cfg *config.Configs) {
			cfg.EpochConfig.EnableEpochs.REWAInMultiTransferEnableEpoch = activationEpoch
			cfg.SystemSCConfig.DCDTSystemSCConfig.BaseIssuingCost = baseIssuingCost
		},
	})
	require.Nil(t, err)
	require.NotNil(t, cs)

	defer cs.Close()

	addrs := createAddresses(t, cs, false)

	err = cs.GenerateBlocksUntilEpochIsReached(int32(activationEpoch))
	require.Nil(t, err)

	// issue metaDCDT
	metaDCDTTicker := []byte("METATICKER")
	nonce := uint64(0)
	tx := issueMetaDCDTTx(nonce, addrs[0].Bytes, metaDCDTTicker, baseIssuingCost)
	nonce++

	txResult, err := cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)
	require.Equal(t, "success", txResult.Status.String())

	metaDCDTTokenID := txResult.Logs.Events[0].Topics[0]

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleTransfer),
	}
	setAddressDcdtRoles(t, cs, nonce, addrs[0], metaDCDTTokenID, roles)
	nonce++

	log.Info("Issued metaDCDT token id", "tokenID", string(metaDCDTTokenID))

	// issue NFT
	nftTicker := []byte("NFTTICKER")
	tx = issueNonFungibleTx(nonce, addrs[0].Bytes, nftTicker, baseIssuingCost)
	nonce++

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)
	require.Equal(t, "success", txResult.Status.String())

	nftTokenID := txResult.Logs.Events[0].Topics[0]
	setAddressDcdtRoles(t, cs, nonce, addrs[0], nftTokenID, roles)
	nonce++

	log.Info("Issued NFT token id", "tokenID", string(nftTokenID))

	// issue SFT
	sftTicker := []byte("SFTTICKER")
	tx = issueSemiFungibleTx(nonce, addrs[0].Bytes, sftTicker, baseIssuingCost)
	nonce++

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)
	require.Equal(t, "success", txResult.Status.String())

	sftTokenID := txResult.Logs.Events[0].Topics[0]
	setAddressDcdtRoles(t, cs, nonce, addrs[0], sftTokenID, roles)
	nonce++

	log.Info("Issued SFT token id", "tokenID", string(sftTokenID))

	nftMetaData := txsFee.GetDefaultMetaData()
	nftMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))

	sftMetaData := txsFee.GetDefaultMetaData()
	sftMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))

	dcdtMetaData := txsFee.GetDefaultMetaData()
	dcdtMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))

	tokenIDs := [][]byte{
		nftTokenID,
		sftTokenID,
		metaDCDTTokenID,
	}

	tokensMetadata := []*txsFee.MetaData{
		nftMetaData,
		sftMetaData,
		dcdtMetaData,
	}

	for i := range tokenIDs {
		tx = dcdtNftCreateTx(nonce, addrs[0].Bytes, tokenIDs[i], tokensMetadata[i], 1)

		txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
		require.Nil(t, err)
		require.NotNil(t, txResult)

		fmt.Println(txResult)
		fmt.Println(string(txResult.Logs.Events[0].Topics[0]))
		fmt.Println(string(txResult.Logs.Events[0].Topics[1]))

		require.Equal(t, "success", txResult.Status.String())

		nonce++
	}

	err = cs.GenerateBlocks(10)
	require.Nil(t, err)

	account0, err := cs.GetAccount(addrs[0])
	require.Nil(t, err)

	beforeBalanceStr0 := account0.Balance

	account1, err := cs.GetAccount(addrs[1])
	require.Nil(t, err)

	beforeBalanceStr1 := account1.Balance

	rewaValue := oneREWA.Mul(oneREWA, big.NewInt(3))
	tx = multiDCDTNFTTransferWithREWATx(nonce, addrs[0].Bytes, addrs[1].Bytes, tokenIDs, rewaValue)

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)

	require.Equal(t, "success", txResult.Status.String())

	rewaLog := string(txResult.Logs.Events[0].Topics[0])
	require.Equal(t, "REWA-000000", rewaLog)

	err = cs.GenerateBlocks(10)
	require.Nil(t, err)

	// check accounts balance
	account0, err = cs.GetAccount(addrs[0])
	require.Nil(t, err)

	beforeBalance0, _ := big.NewInt(0).SetString(beforeBalanceStr0, 10)

	expectedBalance0 := big.NewInt(0).Sub(beforeBalance0, rewaValue)
	txsFee, _ := big.NewInt(0).SetString(txResult.Fee, 10)
	expectedBalanceWithFee0 := big.NewInt(0).Sub(expectedBalance0, txsFee)

	require.Equal(t, expectedBalanceWithFee0.String(), account0.Balance)

	account1, err = cs.GetAccount(addrs[1])
	require.Nil(t, err)

	beforeBalance1, _ := big.NewInt(0).SetString(beforeBalanceStr1, 10)
	expectedBalance1 := big.NewInt(0).Add(beforeBalance1, rewaValue)

	require.Equal(t, expectedBalance1.String(), account1.Balance)
}

func TestChainSimulator_REWA_MultiTransfer_Insufficient_Funds(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	startTime := time.Now().Unix()
	roundDurationInMillis := uint64(6000)
	roundsPerEpoch := core.OptionalUint64{
		HasValue: true,
		Value:    20,
	}

	activationEpoch := uint32(4)

	baseIssuingCost := "1000"

	numOfShards := uint32(3)
	cs, err := chainSimulator.NewChainSimulator(chainSimulator.ArgsChainSimulator{
		BypassTxSignatureCheck:   true,
		TempDir:                  t.TempDir(),
		PathToInitialConfig:      defaultPathToInitialConfig,
		NumOfShards:              numOfShards,
		GenesisTimestamp:         startTime,
		RoundDurationInMillis:    roundDurationInMillis,
		RoundsPerEpoch:           roundsPerEpoch,
		ApiInterface:             api.NewNoApiInterface(),
		MinNodesPerShard:         3,
		MetaChainMinNodes:        3,
		NumNodesWaitingListMeta:  0,
		NumNodesWaitingListShard: 0,
		AlterConfigsFunction: func(cfg *config.Configs) {
			cfg.EpochConfig.EnableEpochs.REWAInMultiTransferEnableEpoch = activationEpoch
			cfg.SystemSCConfig.DCDTSystemSCConfig.BaseIssuingCost = baseIssuingCost
		},
	})
	require.Nil(t, err)
	require.NotNil(t, cs)

	defer cs.Close()

	addrs := createAddresses(t, cs, false)

	err = cs.GenerateBlocksUntilEpochIsReached(int32(activationEpoch))
	require.Nil(t, err)

	// issue NFT
	nftTicker := []byte("NFTTICKER")
	nonce := uint64(0)
	tx := issueNonFungibleTx(nonce, addrs[0].Bytes, nftTicker, baseIssuingCost)
	nonce++

	txResult, err := cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)
	require.Equal(t, "success", txResult.Status.String())

	nftTokenID := txResult.Logs.Events[0].Topics[0]

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleTransfer),
	}
	setAddressDcdtRoles(t, cs, nonce, addrs[0], nftTokenID, roles)
	nonce++

	log.Info("Issued NFT token id", "tokenID", string(nftTokenID))

	nftMetaData := txsFee.GetDefaultMetaData()
	nftMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))

	tx = dcdtNftCreateTx(nonce, addrs[0].Bytes, nftTokenID, nftMetaData, 1)
	nonce++

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)

	require.Equal(t, "success", txResult.Status.String())

	err = cs.GenerateBlocks(10)
	require.Nil(t, err)

	account0, err := cs.GetAccount(addrs[0])
	require.Nil(t, err)

	beforeBalanceStr0 := account0.Balance

	account1, err := cs.GetAccount(addrs[1])
	require.Nil(t, err)

	beforeBalanceStr1 := account1.Balance

	rewaValue, _ := big.NewInt(0).SetString(beforeBalanceStr0, 10)
	rewaValue = rewaValue.Add(rewaValue, big.NewInt(13))
	tx = multiDCDTNFTTransferWithREWATx(nonce, addrs[0].Bytes, addrs[1].Bytes, [][]byte{nftTokenID}, rewaValue)

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)

	require.NotEqual(t, "success", txResult.Status.String())

	eventLog := string(txResult.Logs.Events[0].Topics[1])
	require.Equal(t, "insufficient funds for token REWA-000000", eventLog)

	// check accounts balance
	account0, err = cs.GetAccount(addrs[0])
	require.Nil(t, err)

	beforeBalance0, _ := big.NewInt(0).SetString(beforeBalanceStr0, 10)

	txsFee, _ := big.NewInt(0).SetString(txResult.Fee, 10)
	expectedBalanceWithFee0 := big.NewInt(0).Sub(beforeBalance0, txsFee)

	require.Equal(t, expectedBalanceWithFee0.String(), account0.Balance)

	account1, err = cs.GetAccount(addrs[1])
	require.Nil(t, err)

	require.Equal(t, beforeBalanceStr1, account1.Balance)
}

func TestChainSimulator_REWA_MultiTransfer_Invalid_Value(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	startTime := time.Now().Unix()
	roundDurationInMillis := uint64(6000)
	roundsPerEpoch := core.OptionalUint64{
		HasValue: true,
		Value:    20,
	}

	activationEpoch := uint32(4)

	baseIssuingCost := "1000"

	numOfShards := uint32(3)
	cs, err := chainSimulator.NewChainSimulator(chainSimulator.ArgsChainSimulator{
		BypassTxSignatureCheck:   true,
		TempDir:                  t.TempDir(),
		PathToInitialConfig:      defaultPathToInitialConfig,
		NumOfShards:              numOfShards,
		GenesisTimestamp:         startTime,
		RoundDurationInMillis:    roundDurationInMillis,
		RoundsPerEpoch:           roundsPerEpoch,
		ApiInterface:             api.NewNoApiInterface(),
		MinNodesPerShard:         3,
		MetaChainMinNodes:        3,
		NumNodesWaitingListMeta:  0,
		NumNodesWaitingListShard: 0,
		AlterConfigsFunction: func(cfg *config.Configs) {
			cfg.EpochConfig.EnableEpochs.REWAInMultiTransferEnableEpoch = activationEpoch
			cfg.SystemSCConfig.DCDTSystemSCConfig.BaseIssuingCost = baseIssuingCost
		},
	})
	require.Nil(t, err)
	require.NotNil(t, cs)

	defer cs.Close()

	addrs := createAddresses(t, cs, false)

	err = cs.GenerateBlocksUntilEpochIsReached(int32(activationEpoch))
	require.Nil(t, err)

	// issue NFT
	nftTicker := []byte("NFTTICKER")
	nonce := uint64(0)
	tx := issueNonFungibleTx(nonce, addrs[0].Bytes, nftTicker, baseIssuingCost)
	nonce++

	txResult, err := cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)
	require.Equal(t, "success", txResult.Status.String())

	nftTokenID := txResult.Logs.Events[0].Topics[0]

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleTransfer),
	}
	setAddressDcdtRoles(t, cs, nonce, addrs[0], nftTokenID, roles)
	nonce++

	log.Info("Issued NFT token id", "tokenID", string(nftTokenID))

	nftMetaData := txsFee.GetDefaultMetaData()
	nftMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))

	tx = dcdtNftCreateTx(nonce, addrs[0].Bytes, nftTokenID, nftMetaData, 1)
	nonce++

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)

	require.Equal(t, "success", txResult.Status.String())

	err = cs.GenerateBlocks(10)
	require.Nil(t, err)

	account0, err := cs.GetAccount(addrs[0])
	require.Nil(t, err)

	beforeBalanceStr0 := account0.Balance

	account1, err := cs.GetAccount(addrs[1])
	require.Nil(t, err)

	beforeBalanceStr1 := account1.Balance

	rewaValue := oneREWA.Mul(oneREWA, big.NewInt(3))
	tx = multiDCDTNFTTransferWithREWATx(nonce, addrs[0].Bytes, addrs[1].Bytes, [][]byte{nftTokenID}, rewaValue)
	tx.Value = rewaValue // invalid value field

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)

	require.NotEqual(t, "success", txResult.Status.String())

	eventLog := string(txResult.Logs.Events[0].Topics[1])
	require.Equal(t, "built in function called with tx value is not allowed", eventLog)

	// check accounts balance
	account0, err = cs.GetAccount(addrs[0])
	require.Nil(t, err)

	beforeBalance0, _ := big.NewInt(0).SetString(beforeBalanceStr0, 10)

	txsFee, _ := big.NewInt(0).SetString(txResult.Fee, 10)
	expectedBalanceWithFee0 := big.NewInt(0).Sub(beforeBalance0, txsFee)

	require.Equal(t, expectedBalanceWithFee0.String(), account0.Balance)

	account1, err = cs.GetAccount(addrs[1])
	require.Nil(t, err)

	require.Equal(t, beforeBalanceStr1, account1.Balance)
}

func TestChainSimulator_Multiple_REWA_Transfers(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	startTime := time.Now().Unix()
	roundDurationInMillis := uint64(6000)
	roundsPerEpoch := core.OptionalUint64{
		HasValue: true,
		Value:    20,
	}

	activationEpoch := uint32(4)

	baseIssuingCost := "1000"

	numOfShards := uint32(3)
	cs, err := chainSimulator.NewChainSimulator(chainSimulator.ArgsChainSimulator{
		BypassTxSignatureCheck:   true,
		TempDir:                  t.TempDir(),
		PathToInitialConfig:      defaultPathToInitialConfig,
		NumOfShards:              numOfShards,
		GenesisTimestamp:         startTime,
		RoundDurationInMillis:    roundDurationInMillis,
		RoundsPerEpoch:           roundsPerEpoch,
		ApiInterface:             api.NewNoApiInterface(),
		MinNodesPerShard:         3,
		MetaChainMinNodes:        3,
		NumNodesWaitingListMeta:  0,
		NumNodesWaitingListShard: 0,
		AlterConfigsFunction: func(cfg *config.Configs) {
			cfg.EpochConfig.EnableEpochs.REWAInMultiTransferEnableEpoch = activationEpoch
			cfg.SystemSCConfig.DCDTSystemSCConfig.BaseIssuingCost = baseIssuingCost
		},
	})
	require.Nil(t, err)
	require.NotNil(t, cs)

	defer cs.Close()

	addrs := createAddresses(t, cs, false)

	err = cs.GenerateBlocksUntilEpochIsReached(int32(activationEpoch))
	require.Nil(t, err)

	// issue NFT
	nftTicker := []byte("NFTTICKER")
	nonce := uint64(0)
	tx := issueNonFungibleTx(nonce, addrs[0].Bytes, nftTicker, baseIssuingCost)
	nonce++

	txResult, err := cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)
	require.Equal(t, "success", txResult.Status.String())

	nftTokenID := txResult.Logs.Events[0].Topics[0]

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleTransfer),
	}
	setAddressDcdtRoles(t, cs, nonce, addrs[0], nftTokenID, roles)
	nonce++

	log.Info("Issued NFT token id", "tokenID", string(nftTokenID))

	nftMetaData := txsFee.GetDefaultMetaData()
	nftMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))

	tx = dcdtNftCreateTx(nonce, addrs[0].Bytes, nftTokenID, nftMetaData, 1)
	nonce++

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)

	require.Equal(t, "success", txResult.Status.String())

	err = cs.GenerateBlocks(10)
	require.Nil(t, err)

	account0, err := cs.GetAccount(addrs[0])
	require.Nil(t, err)

	beforeBalanceStr0 := account0.Balance

	account1, err := cs.GetAccount(addrs[1])
	require.Nil(t, err)

	beforeBalanceStr1 := account1.Balance

	// multi nft transfer with multiple REWA-000000 tokens
	numTransfers := 3
	encodedReceiver := hex.EncodeToString(addrs[1].Bytes)
	rewaValue := oneREWA.Mul(oneREWA, big.NewInt(3))

	txDataField := []byte(strings.Join(
		[]string{
			core.BuiltInFunctionMultiDCDTNFTTransfer,
			encodedReceiver,
			hex.EncodeToString(big.NewInt(int64(numTransfers)).Bytes()),
			hex.EncodeToString([]byte("REWA-000000")),
			"00",
			hex.EncodeToString(rewaValue.Bytes()),
			hex.EncodeToString(nftTokenID),
			hex.EncodeToString(big.NewInt(1).Bytes()),
			hex.EncodeToString(big.NewInt(int64(1)).Bytes()),
			hex.EncodeToString([]byte("REWA-000000")),
			"00",
			hex.EncodeToString(rewaValue.Bytes()),
		}, "@"),
	)

	tx = &transaction.Transaction{
		Nonce:     nonce,
		SndAddr:   addrs[0].Bytes,
		RcvAddr:   addrs[0].Bytes,
		GasLimit:  10_000_000,
		GasPrice:  minGasPrice,
		Data:      txDataField,
		Value:     big.NewInt(0),
		Version:   1,
		Signature: []byte("dummySig"),
		ChainID:   []byte(configs.ChainID),
	}

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)

	require.Equal(t, "success", txResult.Status.String())

	// check accounts balance
	account0, err = cs.GetAccount(addrs[0])
	require.Nil(t, err)

	beforeBalance0, _ := big.NewInt(0).SetString(beforeBalanceStr0, 10)

	expectedBalance0 := big.NewInt(0).Sub(beforeBalance0, rewaValue)
	expectedBalance0 = big.NewInt(0).Sub(expectedBalance0, rewaValue)
	txsFee, _ := big.NewInt(0).SetString(txResult.Fee, 10)
	expectedBalanceWithFee0 := big.NewInt(0).Sub(expectedBalance0, txsFee)

	require.Equal(t, expectedBalanceWithFee0.String(), account0.Balance)

	account1, err = cs.GetAccount(addrs[1])
	require.Nil(t, err)

	beforeBalance1, _ := big.NewInt(0).SetString(beforeBalanceStr1, 10)
	expectedBalance1 := big.NewInt(0).Add(beforeBalance1, rewaValue)
	expectedBalance1 = big.NewInt(0).Add(expectedBalance1, rewaValue)

	require.Equal(t, expectedBalance1.String(), account1.Balance)
}

func multiDCDTNFTTransferWithREWATx(nonce uint64, sndAdr, rcvAddr []byte, tokens [][]byte, rewaValue *big.Int) *transaction.Transaction {
	transferData := make([]*utils.TransferDCDTData, 0)

	for _, tokenID := range tokens {
		transferData = append(transferData, &utils.TransferDCDTData{
			Token: tokenID,
			Nonce: 1,
			Value: big.NewInt(1),
		})
	}

	numTransfers := len(tokens)
	encodedReceiver := hex.EncodeToString(rcvAddr)
	hexEncodedNumTransfers := hex.EncodeToString(big.NewInt(int64(numTransfers)).Bytes())
	hexEncodedREWA := hex.EncodeToString([]byte("REWA-000000"))
	hexEncodedREWANonce := "00"

	txDataField := []byte(strings.Join(
		[]string{
			core.BuiltInFunctionMultiDCDTNFTTransfer,
			encodedReceiver,
			hexEncodedNumTransfers,
			hexEncodedREWA,
			hexEncodedREWANonce,
			hex.EncodeToString(rewaValue.Bytes()),
		}, "@"),
	)

	for _, td := range transferData {
		hexEncodedToken := hex.EncodeToString(td.Token)
		dcdtValueEncoded := hex.EncodeToString(td.Value.Bytes())
		hexEncodedNonce := "00"
		if td.Nonce != 0 {
			hexEncodedNonce = hex.EncodeToString(big.NewInt(int64(td.Nonce)).Bytes())
		}

		txDataField = []byte(strings.Join([]string{string(txDataField), hexEncodedToken, hexEncodedNonce, dcdtValueEncoded}, "@"))
	}

	tx := &transaction.Transaction{
		Nonce:     nonce,
		SndAddr:   sndAdr,
		RcvAddr:   sndAdr,
		GasLimit:  10_000_000,
		GasPrice:  minGasPrice,
		Data:      txDataField,
		Value:     big.NewInt(0),
		Version:   1,
		Signature: []byte("dummySig"),
		ChainID:   []byte(configs.ChainID),
	}

	return tx
}

func TestChainSimulator_IssueToken_REWATicker(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	startTime := time.Now().Unix()
	roundDurationInMillis := uint64(6000)
	roundsPerEpoch := core.OptionalUint64{
		HasValue: true,
		Value:    20,
	}

	activationEpoch := uint32(4)

	baseIssuingCost := "1000"

	numOfShards := uint32(3)
	cs, err := chainSimulator.NewChainSimulator(chainSimulator.ArgsChainSimulator{
		BypassTxSignatureCheck:   true,
		TempDir:                  t.TempDir(),
		PathToInitialConfig:      defaultPathToInitialConfig,
		NumOfShards:              numOfShards,
		GenesisTimestamp:         startTime,
		RoundDurationInMillis:    roundDurationInMillis,
		RoundsPerEpoch:           roundsPerEpoch,
		ApiInterface:             api.NewNoApiInterface(),
		MinNodesPerShard:         3,
		MetaChainMinNodes:        3,
		NumNodesWaitingListMeta:  0,
		NumNodesWaitingListShard: 0,
		AlterConfigsFunction: func(cfg *config.Configs) {
			cfg.EpochConfig.EnableEpochs.REWAInMultiTransferEnableEpoch = activationEpoch
			cfg.SystemSCConfig.DCDTSystemSCConfig.BaseIssuingCost = baseIssuingCost
		},
	})
	require.Nil(t, err)
	require.NotNil(t, cs)

	defer cs.Close()

	addrs := createAddresses(t, cs, false)

	err = cs.GenerateBlocksUntilEpochIsReached(int32(activationEpoch) - 1)
	require.Nil(t, err)

	log.Info("Initial setup: Issue token (before the activation of REWAInMultiTransferFlag)")

	// issue NFT
	nftTicker := []byte("REWA")
	nonce := uint64(0)
	tx := issueNonFungibleTx(nonce, addrs[0].Bytes, nftTicker, baseIssuingCost)
	nonce++

	txResult, err := cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)
	require.Equal(t, "success", txResult.Status.String())

	nftTokenID := txResult.Logs.Events[0].Topics[0]

	roles := [][]byte{
		[]byte(core.DCDTRoleNFTCreate),
		[]byte(core.DCDTRoleTransfer),
	}
	setAddressDcdtRoles(t, cs, nonce, addrs[0], nftTokenID, roles)
	nonce++

	log.Info("Issued NFT token id", "tokenID", string(nftTokenID))

	nftMetaData := txsFee.GetDefaultMetaData()
	nftMetaData.Nonce = []byte(hex.EncodeToString(big.NewInt(1).Bytes()))

	tx = dcdtNftCreateTx(nonce, addrs[0].Bytes, nftTokenID, nftMetaData, 1)
	nonce++

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)

	require.Equal(t, "success", txResult.Status.String())

	err = cs.GenerateBlocksUntilEpochIsReached(int32(activationEpoch))
	require.Nil(t, err)

	err = cs.GenerateBlocks(10)
	require.Nil(t, err)

	log.Info("Issue token (after activation of REWAInMultiTransferFlag)")

	// should fail issuing token with REWA ticker
	tx = issueNonFungibleTx(nonce, addrs[0].Bytes, nftTicker, baseIssuingCost)

	txResult, err = cs.SendTxAndGenerateBlockTilTxIsExecuted(tx, maxNumOfBlockToGenerateWhenExecutingTx)
	require.Nil(t, err)
	require.NotNil(t, txResult)

	errMessage := string(txResult.Logs.Events[0].Topics[1])
	require.Equal(t, vm.ErrCouldNotCreateNewTokenIdentifier.Error(), errMessage)

	require.Equal(t, "success", txResult.Status.String())
}
