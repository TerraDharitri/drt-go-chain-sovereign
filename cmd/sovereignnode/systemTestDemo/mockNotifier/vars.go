package main

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain-core/hashing/blake2b"
	"github.com/TerraDharitri/drt-go-chain-core/marshal/factory"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
)

const (
	addressLen  = 32
	hashSize    = 32
	lenItemSize = 4
	u64Size     = 8
)

var (
	log                = logger.GetOrCreate("sovereign-mock-notifier")
	hasher             = blake2b.NewBlake2b()
	marshaller, _      = factory.NewMarshalizer("gogo protobuf")
	pubKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(addressLen, "drt")

	wsURL             = "localhost:22111"
	grpcAddress       = ":8085"
	subscribedAddress = "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"
)
