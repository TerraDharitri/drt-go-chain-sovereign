package config

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"

	"github.com/TerraDharitri/drt-go-chain/config"
)

// LoadSovereignGeneralConfig returns the extra config necessary by sovereign by reading it from the provided file
func LoadSovereignGeneralConfig(filepath string) (*config.SovereignConfig, error) {
	cfg := &config.SovereignConfig{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadSovereignEpochConfig returns the epoch config necessary by sovereign by reading it from the provided file
func LoadSovereignEpochConfig(filepath string) (*config.SovereignEpochConfig, error) {
	cfg := &config.SovereignEpochConfig{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
