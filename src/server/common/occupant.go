package common

import (
	"errors"
	"github.com/name5566/leaf/gate"
	"pdk/src/server/algorithm/dezhou"

	//"pdk/src/server/game/internal"
	"pdk/src/server/model"
	"time"
)

type Occupant struct {
	*model.User
	gate.Agent
	Room   IRoom
	Cards  dezhou.Cards
	Pos    uint8 // 玩家座位号，从1开始
	Status int32 // 1为离线状态

	Bet        uint32 // 当前下注
	Actions    chan int32
	WaitAction bool

	Chips     uint32 // 带入的筹码
	HandValue uint32
}

const (
	Occupant_status_InGame  int32 = 3
	Occupant_status_Offline int32 = 1
	Occupant_status_Observe int32 = 2
	Occupant_status_Sitdown int32 = 0
)

func (o *Occupant) GetRoom() IRoom {
	return o.Room
}
func (o *Occupant) SetRoom(m IRoom) {
	o.Room = m
}
func (o *Occupant) SetAction(n int32)error {
	if o.WaitAction {
		o.Actions <- n
		return  nil
	}
	return  errors.New("not your action")
}
func (o *Occupant) GetAction(timeout time.Duration) int32 {
	timer := time.NewTimer(timeout)
	o.WaitAction = true
	select {
	case n := <-o.Actions:
		timer.Stop()
		o.WaitAction = false
		return n
	case <-o.Room.Closed():
		timer.Stop()
		o.WaitAction = false
		return -1
	case <-timer.C:
		o.WaitAction = false
		timer.Stop()
		return -1 // 超时弃牌
	}
}
func (o *Occupant)SetPos(pos uint8)  {
	 o.Pos = pos
}
func (o *Occupant) GetPos() uint8 {
	return o.Pos
}
func (o *Occupant) GetUid() uint32 {
	return  o.Uid
}
func (o *Occupant) WriteMsg(msg interface{}) {
	if o.Status != Occupant_status_Offline {
		o.Agent.WriteMsg(msg)
	}
}

func (o *Occupant) SetData(d interface{}) {
	o.User = d.(*model.User)
}
func (o *Occupant) GetId() uint32 {
	return o.Uid
}

func (o *Occupant) SetObserve() {
	o.Status = Occupant_status_Observe
}

func (o *Occupant) IsObserve() bool {
	return o.Status == Occupant_status_Observe
}

func (o *Occupant) SetOffline() {
	o.Status = Occupant_status_Offline
}

func (o *Occupant) IsOffline() bool {
	return o.Status == Occupant_status_Offline
}

func (o *Occupant) SetSitdown() {
	o.Status = Occupant_status_Sitdown
}

func (o *Occupant) IsSitdown() bool {
	return o.Status == Occupant_status_Sitdown
}

func (o *Occupant) SetGameing() {
	o.Status = Occupant_status_InGame
}

func (o *Occupant) IsGameing() bool {
	return o.Status == Occupant_status_InGame
}

func (o *Occupant) Replace(value *Occupant) {
	o.Pos = value.Pos
	o.Cards = value.Cards
	o.Room = value.Room
}

func NewOccupant(data *model.User, conn gate.Agent) *Occupant {
	o := &Occupant{
		User:    data,
		Agent:   conn,
		Actions: make(chan int32),
	}
	return o
}
