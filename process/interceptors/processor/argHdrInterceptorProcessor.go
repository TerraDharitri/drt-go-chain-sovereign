package processor

import (
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/process"
)

// ArgHdrInterceptorProcessor is the argument for the interceptor processor used for headers (shard, meta and so on)
type ArgHdrInterceptorProcessor struct {
	Headers        dataRetriever.HeadersPool
	BlockBlackList process.TimeCacher
}
