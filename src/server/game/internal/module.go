package internal

import (
	"github.com/golang/glog"
	"github.com/name5566/leaf/module"
	"pdk/src/server/base"
	"pdk/src/server/protocol"
	"reflect"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}
func init() {
	handler(&protocol.CreateRoom{}, OnMessageCreateRoom)
	handler(&protocol.JoinRoom{}, OnMessage)
	handler(&protocol.LeaveRoom{}, OnMessage)
	handler(&protocol.Bet{}, OnMessage)
	handler(&protocol.SitDown{}, OnMessage) //
	handler(&protocol.StandUp{}, OnMessage) //
	handler(&protocol.Chat{}, OnMessage)    //
}

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton

}

func (m *Module) OnDestroy() {
	glog.Errorln("game OnDestroy")
}
