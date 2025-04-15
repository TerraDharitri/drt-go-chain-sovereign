package antiflood

import "github.com/TerraDharitri/drt-go-chain/process"

func (af *p2pAntiflood) Debugger() process.AntifloodDebugger {
	return af.debugger
}
