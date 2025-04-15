package processorV2

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"sync"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain-core/data/smartContractResult"
	"github.com/TerraDharitri/drt-go-chain-core/data/typeConverters/uint64ByteSlice"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/sha256"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
	"github.com/TerraDharitri/drt-go-chain-vm-common/builtInFunctions"
	mockVm "github.com/TerraDharitri/drt-go-chain-vm-common/mock"
	"github.com/TerraDharitri/drt-go-chain-vm-common/parsers"
	"github.com/stretchr/testify/require"

	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/forking"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/TerraDharitri/drt-go-chain/process/coordinator"
	"github.com/TerraDharitri/drt-go-chain/process/mock"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/hooks"
	"github.com/TerraDharitri/drt-go-chain/process/smartContract/scrCommon"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/state/factory"
	"github.com/TerraDharitri/drt-go-chain/state/storagePruningManager"
	"github.com/TerraDharitri/drt-go-chain/state/storagePruningManager/evictionWaitingList"
	"github.com/TerraDharitri/drt-go-chain/state/syncer"
	"github.com/TerraDharitri/drt-go-chain/storage/txcache"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/economicsmocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genericMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	stateMock "github.com/TerraDharitri/drt-go-chain/testscommon/state"
	storageMock "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/TerraDharitri/drt-go-chain/trie"
)

var (
	dcdtKeyPrefix = []byte(core.ProtectedKeyPrefix + core.DCDTKeyIdentifier)
	sovHasher     = sha256.NewSha256()
	sovMarshaller = &marshal.GogoProtoMarshalizer{}
	sovPubKeyConv = createMockPubkeyConverter()
	sovShardCoord = mock.NewMultiShardsCoordinatorMock(1)

	sovEnableEpochsHandler = enableEpochsHandlerMock.NewEnableEpochsHandlerStub(
		common.SaveToSystemAccountFlag,
		common.OptimizeNFTStoreFlag,
		common.SendAlwaysFlag,
		common.DCDTNFTImprovementV1Flag,
	)
)

func createSovereignSmartContractProcessorArguments() scrCommon.ArgsNewSmartContractProcessor {
	accAdapter := createAccountsDB()
	dcdtStorage := createDCDTDataStorage(accAdapter)
	builtInFuncContainer := createBuiltInFuncContainer(accAdapter, dcdtStorage)
	sovBlock := createSovBlockchainHook(accAdapter, builtInFuncContainer, dcdtStorage)
	txTypeHandler := createTxTypeHandler(builtInFuncContainer)

	return scrCommon.ArgsNewSmartContractProcessor{
		VmContainer:         &mock.VMContainerMock{},
		ArgsParser:          smartContract.NewArgumentParser(),
		Hasher:              &hashingMocks.HasherMock{},
		Marshalizer:         &marshal.GogoProtoMarshalizer{},
		AccountsDB:          accAdapter,
		BlockChainHook:      sovBlock,
		BuiltInFunctions:    builtInFuncContainer,
		PubkeyConv:          sovPubKeyConv,
		ShardCoordinator:    sovShardCoord,
		ScrForwarder:        &mock.IntermediateTransactionHandlerMock{},
		BadTxForwarder:      &mock.IntermediateTransactionHandlerMock{},
		TxFeeHandler:        &mock.FeeAccumulatorStub{},
		TxLogsProcessor:     &mock.TxLogsProcessorStub{},
		EconomicsFee:        &economicsmocks.EconomicsHandlerStub{},
		TxTypeHandler:       txTypeHandler,
		GasHandler:          &testscommon.GasHandlerStub{},
		EnableEpochsHandler: sovEnableEpochsHandler,
		GasSchedule:         testscommon.NewGasScheduleNotifierMock(make(map[string]map[string]uint64)),
		WasmVMChangeLocker:  &sync.RWMutex{},
		VMOutputCacher:      txcache.NewDisabledCache(),
	}
}

func createAccountsDB() *state.AccountsDB {
	argsAccCreator := factory.ArgsAccountCreator{
		Hasher:              sovHasher,
		Marshaller:          sovMarshaller,
		EnableEpochsHandler: sovEnableEpochsHandler,
	}
	accCreator, _ := factory.NewAccountCreator(argsAccCreator)

	storageManagerArgs := storageMock.GetStorageManagerArgs()
	storageManagerArgs.Marshalizer = sovMarshaller
	storageManagerArgs.Hasher = sovHasher

	trieFactoryManager, _ := trie.CreateTrieStorageManager(storageManagerArgs, storageMock.GetStorageManagerOptions())
	tr, _ := trie.NewTrie(trieFactoryManager, sovMarshaller, sovHasher, sovEnableEpochsHandler, 5)
	ewlArgs := evictionWaitingList.MemoryEvictionWaitingListArgs{
		RootHashesSize: 100,
		HashesSize:     10000,
	}
	ewl, _ := evictionWaitingList.NewMemoryEvictionWaitingList(ewlArgs)
	spm, _ := storagePruningManager.NewStoragePruningManager(ewl, 10)

	args := state.ArgsAccountsDB{
		Trie:                  tr,
		Hasher:                sovHasher,
		Marshaller:            sovMarshaller,
		AccountFactory:        accCreator,
		StoragePruningManager: spm,
		SnapshotsManager:      &stateMock.SnapshotsManagerStub{},
		AddressConverter:      &testscommon.PubkeyConverterMock{},
	}
	adb, _ := state.NewAccountsDB(args)
	return adb
}

func createDCDTDataStorage(accountsDB state.AccountsAdapter) vmcommon.DCDTNFTStorageHandler {
	dcdtDataStorage, _ := builtInFunctions.NewDCDTDataStorage(builtInFunctions.ArgsNewDCDTDataStorage{
		Accounts:                      accountsDB,
		GlobalSettingsHandler:         &mockVm.GlobalSettingsHandlerStub{},
		Marshalizer:                   sovMarshaller,
		EnableEpochsHandler:           sovEnableEpochsHandler,
		ShardCoordinator:              sovShardCoord,
		CrossChainTokenCheckerHandler: &mockVm.CrossChainTokenCheckerMock{},
	})

	return dcdtDataStorage
}

func createBuiltInFuncContainer(
	accountsDB state.AccountsAdapter,
	dcdtDataStorage vmcommon.DCDTNFTStorageHandler,
) vmcommon.BuiltInFunctionContainer {
	builtInFuncContainer := builtInFunctions.NewBuiltInFunctionContainer()
	dcdtMultiTransfer, _ := builtInFunctions.NewDCDTNFTMultiTransferFunc(
		10,
		sovMarshaller,
		&mockVm.GlobalSettingsHandlerStub{},
		accountsDB,
		sovShardCoord,
		vmcommon.BaseOperationCost{},
		sovEnableEpochsHandler,
		&mockVm.DCDTRoleHandlerStub{},
		dcdtDataStorage,
	)

	_ = dcdtMultiTransfer.SetPayableChecker(&mockVm.PayableHandlerStub{})
	_ = builtInFuncContainer.Add(core.BuiltInFunctionMultiDCDTNFTTransfer, dcdtMultiTransfer)

	return builtInFuncContainer
}

func createSovBlockchainHook(
	accountsDB state.AccountsAdapter,
	builtInFuncContainer vmcommon.BuiltInFunctionContainer,
	dcdtStorage vmcommon.DCDTNFTStorageHandler,
) process.BlockChainHookHandler {
	args := hooks.ArgBlockChainHook{
		Accounts:                 accountsDB,
		PubkeyConv:               sovPubKeyConv,
		StorageService:           genericMocks.NewChainStorerMock(0),
		DataPool:                 dataRetriever.NewPoolsHolderMock(),
		BlockChain:               &testscommon.ChainHandlerMock{},
		ShardCoordinator:         sovShardCoord,
		Marshalizer:              sovMarshaller,
		Uint64Converter:          uint64ByteSlice.NewBigEndianConverter(),
		BuiltInFunctions:         builtInFuncContainer,
		NFTStorageHandler:        dcdtStorage,
		GlobalSettingsHandler:    &testscommon.DCDTGlobalSettingsHandlerStub{},
		CompiledSCPool:           &testscommon.TxCacherStub{},
		ConfigSCStorage:          config.StorageConfig{},
		EnableEpochs:             config.EnableEpochs{},
		EpochNotifier:            forking.NewGenericEpochNotifier(),
		EnableEpochsHandler:      sovEnableEpochsHandler,
		WorkingDir:               "",
		NilCompiledSCStore:       true,
		GasSchedule:              testscommon.NewGasScheduleNotifierMock(make(map[string]map[string]uint64)),
		Counter:                  &testscommon.BlockChainHookCounterStub{},
		MissingTrieNodesNotifier: syncer.NewMissingTrieNodesNotifier(),
	}

	blockChainHook, _ := hooks.NewBlockChainHookImpl(args)
	sovBlockchainHook, _ := hooks.NewSovereignBlockChainHook(blockChainHook)
	return sovBlockchainHook
}

func createTxTypeHandler(builtInFuncContainer vmcommon.BuiltInFunctionContainer) process.TxTypeHandler {
	dcdtParser, _ := parsers.NewDCDTTransferParser(sovMarshaller)
	txTypeHandler, _ := coordinator.NewTxTypeHandler(coordinator.ArgNewTxTypeHandler{
		PubkeyConverter:     sovPubKeyConv,
		ShardCoordinator:    sovShardCoord,
		BuiltInFunctions:    builtInFuncContainer,
		ArgumentParser:      parsers.NewCallArgsParser(),
		DCDTTransferParser:  dcdtParser,
		EnableEpochsHandler: sovEnableEpochsHandler,
	})

	return txTypeHandler
}

func requireTokenExists(
	t *testing.T,
	account vmcommon.AccountHandler,
	tokenName []byte,
	nonce uint64,
	expectedValue *big.Int,
) {
	tokenId := append(dcdtKeyPrefix, tokenName...)
	dcdtNFTTokenKey := computeDCDTNFTTokenKey(tokenId, nonce)
	dcdtData := &dcdt.DCDigitalToken{}
	marshaledData, _, err := account.(vmcommon.UserAccountHandler).AccountDataHandler().RetrieveValue(dcdtNFTTokenKey)
	require.Nil(t, err)

	err = sovMarshaller.Unmarshal(dcdtData, marshaledData)
	require.Nil(t, err)
	require.Equal(t, expectedValue, dcdtData.Value)
}

func computeDCDTNFTTokenKey(dcdtTokenKey []byte, nonce uint64) []byte {
	return append(dcdtTokenKey, big.NewInt(0).SetUint64(nonce).Bytes()...)
}

func createNFTMetaData(value *big.Int, nonce uint64, creator []byte) []byte {
	dcdtData := &dcdt.DCDigitalToken{
		Type:  uint32(core.NonFungible),
		Value: value,
		TokenMetaData: &dcdt.MetaData{
			URIs:       [][]byte{[]byte("uri1"), []byte("uri2"), []byte("uri3")},
			Nonce:      nonce,
			Hash:       []byte("NFT hash"),
			Name:       []byte("name nft"),
			Attributes: []byte("attributes"),
			Creator:    creator,
		},
	}

	nftMetaData, _ := sovMarshaller.Marshal(dcdtData)
	return nftMetaData
}

func generateRandomBytes(len int) []byte {
	randomBytes := make([]byte, len)
	_, _ = rand.Read(randomBytes)
	return randomBytes
}

func TestNewSovereignSCRProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil ArgsParser should err", func(t *testing.T) {
		args := createSSCProcessArgs()
		args.ArgsParser = nil

		sovProc, err := NewSovereignSCRProcessor(args)
		require.Nil(t, sovProc)
		require.Equal(t, process.ErrNilArgumentParser, err)
	})
	t.Run("nil TxTypeHandler should err", func(t *testing.T) {
		args := createSSCProcessArgs()
		args.TxTypeHandler = nil

		sovProc, err := NewSovereignSCRProcessor(args)
		require.Nil(t, sovProc)
		require.Equal(t, process.ErrNilTxTypeHandler, err)
	})
	t.Run("nil SmartContractProcessor should err", func(t *testing.T) {
		args := createSSCProcessArgs()
		args.SmartContractProcessor = nil

		sovProc, err := NewSovereignSCRProcessor(args)
		require.Nil(t, sovProc)
		require.Equal(t, process.ErrNilSmartContractResultProcessor, err)
	})
	t.Run("nil SCProcessorHelperHandler should err", func(t *testing.T) {
		args := createSSCProcessArgs()
		args.SCProcessorHelperHandler = nil

		sovProc, err := NewSovereignSCRProcessor(args)
		require.Nil(t, sovProc)
		require.Equal(t, process.ErrNilSCProcessorHelper, err)
	})
	t.Run("should work", func(t *testing.T) {
		args := createSSCProcessArgs()

		sovProc, err := NewSovereignSCRProcessor(args)
		require.NotNil(t, sovProc)
		require.Nil(t, err)
	})
}

func TestSovereignSCProcessor_ProcessSmartContractResultIncomingSCR(t *testing.T) {
	t.Parallel()

	arguments := createSovereignSmartContractProcessorArguments()
	sc, _ := NewSmartContractProcessorV2(arguments)
	scpHelper, _ := scrCommon.NewSCProcessorHelper(scrCommon.SCProcessorHelperArgs{
		Accounts:         arguments.AccountsDB,
		ShardCoordinator: arguments.ShardCoordinator,
		Marshalizer:      arguments.Marshalizer,
		Hasher:           arguments.Hasher,
		PubkeyConverter:  arguments.PubkeyConv,
	})
	sovProc, _ := NewSovereignSCRProcessor(SovereignSCProcessArgs{
		ArgsParser:               arguments.ArgsParser,
		TxTypeHandler:            arguments.TxTypeHandler,
		SmartContractProcessor:   sc,
		SCProcessorHelperHandler: scpHelper,
	})

	scAddress := generateRandomBytes(32)

	token1 := []byte("token1")
	nftTransferNonce := big.NewInt(4)
	nftTransferValue := big.NewInt(100)
	nftMetaData := createNFTMetaData(nftTransferValue, nftTransferNonce.Uint64(), scAddress)
	transferNFT :=
		hex.EncodeToString(token1) + "@" + // id
			hex.EncodeToString(nftTransferNonce.Bytes()) + "@" + // nonce != 0
			hex.EncodeToString(nftMetaData) + "@" // meta data

	token2 := []byte("token2")
	dcdtTransferNonce := big.NewInt(0)
	dcdtTransferVal := big.NewInt(50)
	transferDCDT :=
		hex.EncodeToString(token2) + "@" + // id
			hex.EncodeToString(dcdtTransferNonce.Bytes()) + "@" + // nonce = 0
			hex.EncodeToString(dcdtTransferVal.Bytes()) // value

	scr := smartContractResult.SmartContractResult{
		SndAddr: core.DCDTSCAddress,
		RcvAddr: scAddress,
		Data:    []byte("MultiDCDTNFTTransfer@02@" + transferNFT + transferDCDT),
		Value:   big.NewInt(0),
	}
	_, err := sovProc.ProcessSmartContractResult(&scr)
	require.Nil(t, err)

	acc, err := arguments.AccountsDB.LoadAccount(scAddress)
	require.Nil(t, err)
	requireTokenExists(t, acc, token1, nftTransferNonce.Uint64(), nftTransferValue)
	requireTokenExists(t, acc, token2, dcdtTransferNonce.Uint64(), dcdtTransferVal)
}

func createSSCProcessArgs() SovereignSCProcessArgs {
	arguments := createSovereignSmartContractProcessorArguments()
	sc, _ := NewSmartContractProcessorV2(arguments)
	scpHelper, _ := scrCommon.NewSCProcessorHelper(scrCommon.SCProcessorHelperArgs{
		Accounts:         arguments.AccountsDB,
		ShardCoordinator: arguments.ShardCoordinator,
		Marshalizer:      arguments.Marshalizer,
		Hasher:           arguments.Hasher,
		PubkeyConverter:  arguments.PubkeyConv,
	})
	return SovereignSCProcessArgs{
		ArgsParser:               arguments.ArgsParser,
		TxTypeHandler:            arguments.TxTypeHandler,
		SmartContractProcessor:   sc,
		SCProcessorHelperHandler: scpHelper,
	}
}
