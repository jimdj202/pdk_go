package internal

import (
	"errors"
	"github.com/name5566/leaf/gate"
	"pdk/src/server/algorithm"
	"pdk/src/server/model"
	"time"
)

type Occupant struct {
	*model.User
	gate.Agent
	room   IRoom
	cards  algorithm.Cards
	Pos    uint8 // 玩家座位号，从1开始
	status int32 // 1为离线状态

	Bet        uint32 // 当前下注
	actions    chan int32
	waitAction bool

	chips     uint32 // 带入的筹码
	HandValue      uint32
}

const (
	Occupant_status_InGame  int32 = 3
	Occupant_status_Offline int32 = 1
	Occupant_status_Observe int32 = 2
	Occupant_status_Sitdown int32 = 0
)

func (o *Occupant) GetRoom() IRoom {
	return o.room
}
func (o *Occupant) SetRoom(m IRoom) {
	o.room = m
}
func (o *Occupant) SetAction(n int32)error {
	if o.waitAction {
		o.actions <- n
		return  nil
	}
	return  errors.New("not your action")
}
func (o *Occupant) GetAction(timeout time.Duration) int32 {
	timer := time.NewTimer(timeout)
	o.waitAction = true
	select {
	case n := <-o.actions:
		timer.Stop()
		o.waitAction = false
		return n
	case <-o.room.Closed():
		timer.Stop()
		o.waitAction = false
		return -1
	case <-timer.C:
		o.waitAction = false
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
	if o.status != Occupant_status_Offline {
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
	o.status = Occupant_status_Observe
}

func (o *Occupant) IsObserve() bool {
	return o.status == Occupant_status_Observe
}

func (o *Occupant) SetOffline() {
	o.status = Occupant_status_Offline
}

func (o *Occupant) IsOffline() bool {
	return o.status == Occupant_status_Offline
}

func (o *Occupant) SetSitdown() {
	o.status = Occupant_status_Sitdown
}

func (o *Occupant) IsSitdown() bool {
	return o.status == Occupant_status_Sitdown
}

func (o *Occupant) SetGameing() {
	o.status = Occupant_status_InGame
}

func (o *Occupant) IsGameing() bool {
	return o.status == Occupant_status_InGame
}

func (o *Occupant) Replace(value *Occupant) {
	o.Pos = value.Pos
	o.cards = value.cards
	o.room = value.room
}

func NewOccupant(data *model.User, conn gate.Agent) *Occupant {
	o := &Occupant{
		User:    data,
		Agent:   conn,
		actions: make(chan int32),
	}
	return o
}
