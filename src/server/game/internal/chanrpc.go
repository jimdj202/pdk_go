package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/golang/glog"
	"pdk/src/server/model"
	"pdk/src/server/game/room"
)

func init() {
	skeleton.RegisterChanRPC(model.Agent_New, rpcNewAgent)
	skeleton.RegisterChanRPC(model.Agent_Close, rpcCloseAgent)
	skeleton.RegisterChanRPC(model.Agent_Login, rpcLoginAgent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	glog.Infof("新建链接 %T", a)
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	glog.Infof("链接关闭 %T", a)
}

func rpcLoginAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	u := args[1].(*model.User)
	o := NewOccupant(u, a)
	a.SetUserData(o)

	if len(u.RoomID) > 0 {
		o.room = room.GetRoom(u.RoomID)
	}
	glog.Infof("rpcLoginAgent %T", u)
}
