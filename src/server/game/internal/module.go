package internal

import (
	"github.com/name5566/leaf/module"
	"pdk/src/server/base"
	"github.com/golang/glog"
	"pdk/src/server/protocol"
	"pdk/src/server/game/room"
	"reflect"
	"pdk/src/server/model"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}
func init() {
	handler(&protocol.JoinRoom{}, room.OnMessage)
	handler(&protocol.LeaveRoom{}, room.OnMessage)
	handler(&protocol.Bet{}, room.OnMessage)
	handler(&protocol.SitDown{}, room.OnMessage) //
	handler(&protocol.StandUp{}, room.OnMessage) //
	handler(&protocol.Chat{}, room.OnMessage)    //
}

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	room.Init(&Creator{})
}

func (m *Module) OnDestroy() {
	glog.Errorln("game OnDestroy")
}

type Creator struct{}

// 对玩家未进入房间，或者没房间数据的处理
func (this *Creator) Create(m interface{}) room.IRoom {
	if msg, ok := m.(*protocol.JoinRoom); ok {
		if len(msg.RoomNumber) == 0 {
			r := room.FindRoom()
			return r
		}
		r := room.GetRoom(msg.RoomNumber)
		if r != nil {
			return r
		}
		room := NewRoom(9, 5, 10, 1000, model.Timeout)
		room.Insert()
		return room
	}
	return nil
}
