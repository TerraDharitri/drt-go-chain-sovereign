package config

import "github.com/TerraDharitri/drt-go-chain/config"

// SovereignConfig holds sovereign node config
type SovereignConfig struct {
	*config.Configs
	SovereignExtraConfig *config.SovereignConfig
	SovereignEpochConfig *config.SovereignEpochConfig
}
