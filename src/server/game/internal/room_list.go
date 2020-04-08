package internal

import (
	"github.com/name5566/leaf/gate"
	"math/rand"
	"pdk/src/server/model"
	"pdk/src/server/protocol"
	"strconv"
	"sync"
	"time"
)

func OnMessageCreateRoom(args []interface{}){
	// 收到的 Hello 消息
	m := args[0]
	// 消息的发送者
	a := args[1].(gate.Agent)
	if msg, ok := m.(*protocol.CreateRoom); ok {
		num :=rooms.createNumber()

		room := NewRoom(9, msg.TotalPersion, 10, 1000, model.Timeout)
		room.SetNumber(num)
		room.Insert()

		a.WriteMsg(&protocol.CreateRoomResp{})
		return
	}
	a.WriteMsg(&protocol.CreateRoomResp{})
}

func OnMessage(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0]
	// 消息的发送者
	a := args[1].(gate.Agent)

	o := a.UserData().(IOccupant)
	if o.GetRoom() != nil {
		o.GetRoom().Send(o, m)
	} else {
		if r := CreateRoom(m); r == nil {
			a.WriteMsg(protocol.MSG_NOT_IN_ROOM)
		} else {
			SetRoom(r)
			r.Send(o, m)
		}
	}
}

var rooms *roomlist



func init() {
	rooms = &roomlist{
		M: make(map[string]IRoom, 1000),
	}
}

type roomlist struct {
	M map[string]IRoom
	sync.RWMutex
}

func CreateRoom(m interface{}) IRoom {
	if msg, ok := m.(*protocol.JoinRoom); ok {
		if len(msg.RoomNumber) == 0 {
			r := FindRoom()
			return r
		}
		r := GetRoom(msg.RoomNumber)
		if r != nil {
			return r
		}
		room := NewRoom(9, 5, 10, 1000, model.Timeout)
		room.Insert()
		return room
	}
	return nil
}

func FindRoom() IRoom {
	rooms.Lock()
	for _, v := range rooms.M {
		if v.Len() < v.Cap() {
			return v
		}
	}
	rooms.Unlock()
	return nil
}

func GetRoom(rid string) IRoom {
	rooms.RLock()
	r := rooms.M[rid]
	rooms.RUnlock()
	return r
}

func SetRoom(room IRoom) string {

	rooms.Lock()
	id := rooms.createNumber()
	room.SetNumber(id)
	rooms.M[id] = room
	rooms.Unlock()
	return id
}
func DelRoom(room IRoom) {
	rooms.Lock()
	delete(rooms.M, room.GetNumber())
	rooms.Unlock()
}

func (this *roomlist) createNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var n string
	for i := 0; i < 100; i++ {
		n = strconv.Itoa(int(r.Int31n(999999-100000) + 100000))
		if _, ok := rooms.M[n]; !ok {
			return n
		}
	}
	return n
}

func Each(f func(o IRoom) bool) {
	rooms.RLock()
	for _, v := range rooms.M {
		if !f(v) {
			break
		}
	}
	rooms.RUnlock()
}
func GetRooms() []IRoom {
	r := make([]IRoom, len(rooms.M))
	rooms.RLock()
	var n = 0
	for _, v := range rooms.M {
		r[n] = v
		n ++
	}

	rooms.RUnlock()
	return r
}
