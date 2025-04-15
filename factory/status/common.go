package status

import (
	nodeData "github.com/TerraDharitri/drt-go-chain-core/data"
	outportCore "github.com/TerraDharitri/drt-go-chain-core/data/outport"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/epochStart"
	"github.com/TerraDharitri/drt-go-chain/epochStart/notifier"
	"github.com/TerraDharitri/drt-go-chain/outport"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
)

// CreateSaveValidatorsPubKeysEventHandler creates an epoch start action handler able to save validators pub keys
// in outport handler
func CreateSaveValidatorsPubKeysEventHandler(
	nodesCoordinator nodesCoordinator.NodesCoordinator,
	outportHandler outport.OutportHandler,
) epochStart.ActionHandler {
	subscribeHandler := notifier.NewHandlerForEpochStart(func(hdr nodeData.HeaderHandler) {
		currentEpoch := hdr.GetEpoch()
		validatorsPubKeys, err := nodesCoordinator.GetAllEligibleValidatorsPublicKeys(currentEpoch)
		if err != nil {
			log.Warn("pc.nodesCoordinator.GetAllEligibleValidatorPublicKeys for current epoch failed",
				"epoch", currentEpoch,
				"error", err.Error())
		}

		outportHandler.SaveValidatorsPubKeys(&outportCore.ValidatorsPubKeys{
			ShardID:                hdr.GetShardID(),
			ShardValidatorsPubKeys: outportCore.ConvertPubKeys(validatorsPubKeys),
			Epoch:                  currentEpoch,
		})

	}, func(_ nodeData.HeaderHandler) {}, common.IndexerOrder)

	return subscribeHandler
}
