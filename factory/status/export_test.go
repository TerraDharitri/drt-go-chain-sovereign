package status

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/outport"
	outportDriverFactory "github.com/TerraDharitri/drt-go-chain/outport/factory"
	"github.com/TerraDharitri/drt-go-chain/p2p"
)

// ComputeNumConnectedPeers -
func ComputeNumConnectedPeers(
	appStatusHandler core.AppStatusHandler,
	netMessenger p2p.Messenger,
	suffix string,
) {
	computeNumConnectedPeers(appStatusHandler, netMessenger, suffix)
}

// ComputeConnectedPeers -
func ComputeConnectedPeers(
	appStatusHandler core.AppStatusHandler,
	netMessenger p2p.Messenger,
	suffix string,
) {
	computeConnectedPeers(appStatusHandler, netMessenger, suffix)
}

// MakeHostDriversArgs -
func (scf *statusComponentsFactory) MakeHostDriversArgs() ([]outportDriverFactory.ArgsHostDriverFactory, error) {
	return scf.makeHostDriversArgs()
}

// EpochStartEventHandler -
func (pc *statusComponents) OutportHandler() outport.OutportHandler {
	return pc.outportHandler
}
