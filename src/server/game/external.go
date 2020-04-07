package game

import (
	"pdk/src/server/game/internal"

)

var (
	Module  = new(internal.Module)//建立模块新的
	ChanRPC = internal.ChanRPC
)
func GetRoomsEx() []internal.IRoom {
	//r := make([]internal.IRoom, len(rooms.M))
	//rooms.RLock()
	//var n = 0
	//for _, v := range rooms.M {
	//	r[n] = v
	//	n ++
	//}
	//
	//rooms.RUnlock()
	//return r

	return internal.GetRooms()
}