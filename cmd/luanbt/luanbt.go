package main

import (
	lua "github.com/yuin/gopher-lua"
)

// luanbt is called to get a Lua environment with nbt manipulation ability// lua vm memory limit; 0 is no limit
const memoryLimitMb = 100

func luanbt() *lua.LState {
	L := lua.NewState()
	if memoryLimitMb > 0 {
		L.SetMx(memoryLimitMb)
	}
	return L
}
