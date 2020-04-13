package internal

import (
	"github.com/golang/glog"
	"github.com/name5566/leaf/gate"
	"pdk/src/server/common"
	"pdk/src/server/game"
	//"pdk/src/server/game/internal"
	"pdk/src/server/model"
	"pdk/src/server/protocol"
	"reflect"
	"time"
)

func init() {
	handler(&protocol.UserLoginInfo{}, handlLoginUser)
	handler(&protocol.Version{}, handlVersion)
	handler(&protocol.RoomList{}, onRoomList) //

	handler(&protocol.CreateQinYouQuan{},handlerCreateQinYouQuan)
	handler(&protocol.DeleteQinYouQuan{},handlerDeleteQinYouQuan)
	handler(&protocol.JoinQinYouQuan{},handlerJoinQinYouQuan)
	handler(&protocol.LeaveQinYouQuan{},handlerLeaveQinYouQuan)

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
		Uid: user.Uid,
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

	//array := internal.GetRooms()
	array := game.GetRoomsEx()
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
			Rid:         data.ID,
		}
	}
	msg.Room = rooms
	a.WriteMsg(msg)
}

func handlerCreateQinYouQuan(args []interface{}){
	// 收到的消息
	m := args[0].(*protocol.CreateQinYouQuan)
	// 消息的发送者
	a := args[1].(gate.Agent)
	if m.Uid <0 || len(m.Name) <0 {
		a.WriteMsg(protocol.MSG_Param_Error)
		return
	}

	modelQ := &model.QinYouQuan{Uid: m.Uid,Name: m.Name}
	count,err := modelQ.GetCountsByUid()
	if (count !=0 && err != nil) || count>4 {
		a.WriteMsg(protocol.MSG_Max_Created_Error)
		return
	}

	count,err = modelQ.Create()
	if err!=nil{
		a.WriteMsg(protocol.MSG_DB_Error)
		return
	}

	msg := &protocol.CreateQinYouQuanResp{Qid:modelQ.Qid,Name:modelQ.Name}
	msg.Qid = modelQ.Qid

	a.WriteMsg(msg)

}

func handlerDeleteQinYouQuan(args []interface{}){
	// 收到的消息
	m := args[0].(*protocol.DeleteQinYouQuan)
	// 消息的发送者
	a := args[1].(gate.Agent)
	o := a.UserData().(*common.Occupant)
	if m.Uid <1 || m.Qid <1 || o.Uid != m.Uid {
		a.WriteMsg(protocol.MSG_Param_Error)
		return
	}
	modelQ := &model.QinYouQuan{Uid: m.Uid,Qid: m.Qid}
	err := modelQ.Delete()
	if err!=nil {
		a.WriteMsg(protocol.MSG_DB_Error)
		return
	}
	msg := &protocol.DeleteQinYouQuanResp{Uid: m.Uid,Qid:m.Qid}
	a.WriteMsg(msg)
}

func handlerJoinQinYouQuan(args []interface{}){
	// 收到的消息
	m := args[0].(*protocol.JoinQinYouQuan)
	// 消息的发送者
	a := args[1].(gate.Agent)

	if m.Uid <1 || m.Qid <1  {
		a.WriteMsg(protocol.MSG_Param_Error)
		return
	}
	//查找亲友圈
	modelQ := &model.QinYouQuan{Qid: m.Qid}
	_,err := modelQ.FindOneByQid()
	if err != nil {
		a.WriteMsg(protocol.MSG_Param_Error)
		return
	}

	modelQm := &model.QinYouQuanMember{Qid:m.Qid,Uid:m.Uid}
	_ ,err = modelQm.FindOrCreate()
	if err != nil {
		a.WriteMsg(protocol.MSG_DB_Error)
		return
	}
	msg := &protocol.JoinQinYouQuanResp{Qid: m.Qid,Uid: m.Uid}
	a.WriteMsg(msg)

}

func handlerLeaveQinYouQuan(args []interface{}){
	// 收到的消息
	m := args[0].(*protocol.LeaveQinYouQuan)
	// 消息的发送者
	a := args[1].(gate.Agent)

	if m.Uid < 1 || m.Qid < 1  {
		a.WriteMsg(protocol.MSG_Param_Error)
		return
	}

	modelQ := &model.QinYouQuanMember{Qid:m.Qid,Uid:m.Uid}
	_ ,err := modelQ.DeleteByQidAndUid()
	if err != nil {
		a.WriteMsg(protocol.MSG_DB_Error)
	}

	msg := &protocol.LeaveQinYouQuanResp{Uid: m.Uid,Qid: m.Qid}
	a.WriteMsg(msg)

}

