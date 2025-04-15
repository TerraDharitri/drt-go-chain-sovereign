package requesterscontainer

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/requestHandlers/requesters"
	"github.com/TerraDharitri/drt-go-chain/errors"
	"github.com/TerraDharitri/drt-go-chain/process/factory"
)

type sovereignShardRequestersContainerFactory struct {
	*shardRequestersContainerFactory
}

// NewSovereignShardRequestersContainerFactory creates a new container filled with topic requesters for sovereign shards
func NewSovereignShardRequestersContainerFactory(shardReqContainerFactory *shardRequestersContainerFactory) (*sovereignShardRequestersContainerFactory, error) {
	if check.IfNil(shardReqContainerFactory) {
		return nil, errors.ErrNilShardRequesterContainerFactory
	}

	f := &sovereignShardRequestersContainerFactory{
		shardRequestersContainerFactory: shardReqContainerFactory,
	}

	f.numIntraShardPeers = f.numTotalPeers
	f.numCrossShardPeers = 0

	return f, nil
}

// Create returns a requesters container that will hold all sovereign requesters
func (srcf *sovereignShardRequestersContainerFactory) Create() (dataRetriever.RequestersContainer, error) {
	err := srcf.generateCommonRequesters()
	if err != nil {
		return nil, err
	}

	err = srcf.generateHeaderRequesters(core.SovereignChainShardId)
	if err != nil {
		return nil, err
	}

	err = srcf.generateAccountAndValidatorTrieNodesRequesters(core.SovereignChainShardId)
	if err != nil {
		return nil, err
	}

	err = srcf.generateExtendedShardHeaderRequesters()
	if err != nil {
		return nil, err
	}

	return srcf.container, nil
}

func (srcf *sovereignShardRequestersContainerFactory) generateCommonRequesters() error {
	err := srcf.generateTxRequesters(factory.TransactionTopic)
	if err != nil {
		return err
	}

	err = srcf.generateTxRequesters(factory.UnsignedTransactionTopic)
	if err != nil {
		return err
	}

	err = srcf.generateMiniBlocksRequesters()
	if err != nil {
		return err
	}

	err = srcf.generatePeerAuthenticationRequester()
	if err != nil {
		return err
	}

	validatorInfoTopicID := common.ValidatorInfoTopic + srcf.shardCoordinator.CommunicationIdentifier(core.SovereignChainShardId)
	err = srcf.generateValidatorInfoRequester(validatorInfoTopicID)
	if err != nil {
		return err
	}

	return nil
}

func (srcf *sovereignShardRequestersContainerFactory) generateTxRequesters(topic string) error {
	shardC := srcf.shardCoordinator
	identifierTx := topic + shardC.CommunicationIdentifier(core.SovereignChainShardId)

	requester, err := srcf.createTxRequester(identifierTx, EmptyExcludePeersOnTopic, core.SovereignChainShardId, srcf.numCrossShardPeers, srcf.numIntraShardPeers)
	if err != nil {
		return err
	}

	return srcf.container.Add(identifierTx, requester)
}

func (srcf *sovereignShardRequestersContainerFactory) generateMiniBlocksRequesters() error {
	shardC := srcf.shardCoordinator

	identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(shardC.SelfId())
	requester, err := srcf.createMiniBlocksRequester(identifierMiniBlocks, EmptyExcludePeersOnTopic, srcf.shardCoordinator.SelfId(), srcf.numCrossShardPeers, srcf.numIntraShardPeers)
	if err != nil {
		return err
	}

	return srcf.container.Add(identifierMiniBlocks, requester)
}

func (srcf *sovereignShardRequestersContainerFactory) generateExtendedShardHeaderRequesters() error {
	shardC := srcf.shardCoordinator
	shardID := shardC.SelfId()

	identifierHdr := factory.ExtendedHeaderProofTopic + shardC.CommunicationIdentifier(shardID)
	requestSender, err := srcf.createOneRequestSenderWithSpecifiedNumRequests(identifierHdr, EmptyExcludePeersOnTopic, shardID, 0, srcf.numIntraShardPeers)
	if err != nil {
		return err
	}

	arg := requesters.ArgHeaderRequester{
		ArgBaseRequester: requesters.ArgBaseRequester{
			RequestSender: requestSender,
			Marshaller:    srcf.marshaller,
		},
		NonceConverter: srcf.uint64ByteSliceConverter,
	}
	requester, err := requesters.NewHeaderRequester(arg)
	if err != nil {
		return err
	}

	return srcf.container.Add(identifierHdr, requester)
}

// IsInterfaceNil returns true if there is no value under the interface
func (srcf *sovereignShardRequestersContainerFactory) IsInterfaceNil() bool {
	return srcf == nil
}
