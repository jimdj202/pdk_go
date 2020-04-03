package internal

import (
	"github.com/golang/glog"
	"github.com/name5566/leaf/gate"
	"pdk/src/server/game"
	"pdk/src/server/game/room"
	"pdk/src/server/model"
	"pdk/src/server/protocol"
	"reflect"
	"time"
)

func init() {
	handler(&protocol.UserLoginInfo{}, handlLoginUser)
	handler(&protocol.Version{}, handlVersion)
	handler(&protocol.RoomList{}, onRoomList) //
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handlVersion(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0].(*protocol.UserLoginInfoResp)
	// 消息的发送者
	a := args[1].(gate.Agent)
	glog.Infoln(m)
	a.WriteMsg(m)
}

func handlLoginUser(args []interface{}) {
	// 收到的消息
	m := args[0].(*protocol.UserLoginInfo)
	// 消息的发送者
	a := args[1].(gate.Agent)

	user := &model.User{UnionId: m.UnionId}
	exist, err := user.GetByUnionId()
	if err != nil {
		a.WriteMsg(protocol.MSG_DB_Error)
		return
	}

	if !exist {
		user = &model.User{Nickname: m.Nickname,
			UnionId: m.UnionId,LastTime:time.Now()}
		err := user.Create()
		if err != nil {
			a.WriteMsg(protocol.MSG_User_Not_Exist)
			return
		}
	}else {
		//更新最近登录时间
		user.UpdateLoginTimeIp(a.LocalAddr().String())
	}

	resp := &protocol.UserLoginInfoResp{
		Nickname: user.Nickname,
		Account:  user.Account,
		UnionId:  user.UnionId,
	}

	a.WriteMsg(resp)
	game.ChanRPC.Go(model.Agent_Login, a,user)
}

func onRoomList(args []interface{}) {
	// 收到的消息
	//m := args[0].(*protocol.RoomList)
	// 消息的发送者
	a := args[1].(gate.Agent)

	msg := &protocol.RoomListResp{}

	array := room.GetRooms()
	rooms := make([]*protocol.Room, len(array))

	for k, v := range array {
		d := v.Data()
		data := d.(*model.Room)
		rooms[k] = &protocol.Room{

			Number:      data.Number,
			MaxCap:      v.Cap(),
			Cap:         v.Len(),
			DraginChips: data.DraginChips,
			CreatedAt:   data.CreatedTime(),
			Rid:         data.Rid,
		}
	}
	msg.Room = rooms
	a.WriteMsg(msg)
}
