package runType

import (
	"fmt"

	"github.com/TerraDharitri/drt-go-sdk-abi/abi"

	"github.com/TerraDharitri/drt-go-chain/cmd/sovereignnode/dataCodec"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/factory/runType"
	"github.com/TerraDharitri/drt-go-chain/process/block/sovereign/incomingHeader"
)

const (
	separator = "@"
)

// CreateSovereignArgsRunTypeComponents creates the args for run type component
func CreateSovereignArgsRunTypeComponents(
	argsRunType runType.ArgsRunTypeComponents,
	configs config.SovereignConfig,
) (*runType.ArgsSovereignRunTypeComponents, error) {
	runTypeComponentsFactory, err := runType.NewRunTypeComponentsFactory(argsRunType)
	if err != nil {
		return nil, fmt.Errorf("NewRunTypeComponentsFactory failed: %w", err)
	}

	serializer, err := abi.NewSerializer(abi.ArgsNewSerializer{
		PartsSeparator: separator,
	})
	if err != nil {
		return nil, err
	}

	dataCodecHandler, err := dataCodec.NewDataCodec(serializer)
	if err != nil {
		return nil, err
	}

	return &runType.ArgsSovereignRunTypeComponents{
		RunTypeComponentsFactory: runTypeComponentsFactory,
		Config:                   configs,
		DataCodec:                dataCodecHandler,
		TopicsChecker:            incomingHeader.NewTopicsChecker(),
	}, nil
}
