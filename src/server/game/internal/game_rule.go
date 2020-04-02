package internal

import (
	"github.com/golang/glog"
	"pdk/src/server/protocol"
	"pdk/src/server/model"
	"pdk/src/server/algorithm"
	"time"
	"pdk/src/server/lib/utils"
)

func (r *Room) startDelay(startDelay *startDelay, o *Occupant) {
	if startDelay.kind == 0 {
		r.Info()
		r.start()
	}
}
func (r *Room) start() {
	if r.status == RUNNING {
		return
	}
	// 产生庄
	var dealer *Occupant
	button := r.Button - 1
	r.Each((button+1)%r.Cap(), func(o *Occupant) bool {
		r.Button = o.Pos
		dealer = o
		return false
	})

	if dealer == nil {
		return
	}

	r.remain = 0
	r.allin = 0
	// 剔除筹码小于大盲和离线的玩家

	n := 0
	r.Each(0, func(o *Occupant) bool {
		if o.chips < r.BB || o.IsOffline() {
			o.SetSitdown()
			return true
		}
		o.SetGameing()
		n ++
		return true
	})

	// 2人及以上才开始游戏
	if n < 2 {
		return
	}

	r.status = RUNNING
	// 洗牌
	r.Cards.Shuffle()

	// 产生小盲
	sb := r.next(dealer.Pos)
	if n == 2 { // one-to-one
		sb = dealer
	}
	// 产生大盲
	bb := r.next(sb.Pos)
	bbPos := bb.Pos

	// 通报本局庄家
	r.WriteMsg(&protocol.Button{Uid: dealer.Uid})

	// 小大盲下注
	r.betting(sb, int32(r.SB))
	r.betting(bb, int32(r.BB))

	// Round 1 : preflop
	r.ready()
	r.Each(0, func(o *Occupant) bool {
		o.cards = algorithm.Cards{r.Cards.Take(), r.Cards.Take()}

		kind, _ := algorithm.De(o.cards.GetType())
		m := &protocol.PreFlop{
			Cards: o.cards.Bytes(),
			Kind:  kind,
		}
		o.WriteMsg(m)
		return true
	})
	r.Broadcast(&protocol.PreFlop{}, false)

	r.action(0)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

	// Round 2 : Flop
	r.ready()
	r.Cards = algorithm.Cards{r.Cards.Take(), r.Cards.Take(), r.Cards.Take()}
	r.Each(0, func(o *Occupant) bool {
		cs := r.Cards.Append(o.cards...)

		kind, _ := algorithm.De(cs.GetType())
		m := &protocol.Flop{
			Cards: cs.Bytes(),
			Kind:  kind,
		}
		o.WriteMsg(m)
		return true
	})
	r.Broadcast(&protocol.Flop{Cards: r.Cards.Bytes()}, false)

	r.action(0)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

	// Round 3 : Turn
	r.ready()
	r.Cards = r.Cards.Append(r.Cards.Take())
	r.Each(0, func(o *Occupant) bool {
		cs := r.Cards.Append(o.cards...)
		kind, _ := algorithm.De(cs.GetType())
		m := &protocol.Turn{
			Card: r.Cards[3],
			Kind: kind,
		}
		o.WriteMsg(m)
		return true
	})
	r.Broadcast(&protocol.Turn{Card: r.Cards[3]}, false)

	r.action(0)

	if r.remain <= 1 {
		goto showdown
	}
	r.calc()

	// Round 4 : River
	r.ready()
	r.Cards = r.Cards.Append(r.Cards.Take())
	r.Each(0, func(o *Occupant) bool {
		cs := r.Cards.Append(o.cards...)
		value := cs.GetType()
		kind, _ := algorithm.De(value)
		m := &protocol.River{
			Card: r.Cards[4],
			Kind: kind,
		}
		o.WriteMsg(m)
		o.HandValue = value
		return true
	})
	r.Broadcast(&protocol.River{Card: r.Cards[4]}, false)

	r.action(0)

showdown:
	r.showdown()
	showdown := &protocol.Showdown{}
	for _, o := range r.Occupants {
		if o != nil && o.IsGameing() {
			o.SetSitdown()

			item := &protocol.ShowdownItem{
				Uid:      o.Uid,
				ChipsWin: r.Chips[o.Pos-1],
				Chips:    o.chips,
			}
			showdown.Showdown = append(showdown.Showdown, item)
		}
	}
	r.Broadcast(showdown, true)
	r.Info(sb.Pos, bbPos)

	r.status = GAMEOVER

	time.AfterFunc(time.Second*2, func() {
		defer utils.PrintPanicStack()
		r.Send(nil, &startDelay{})
	})
}

func (r *Room) calc() (pots []handPot) {
	pots = calcPot(r.Chips)
	r.Pot = r.Pot[:]
	var ps []uint32
	for _, pot := range pots {
		r.Pot = append(r.Pot, pot.Pot)
		ps = append(ps, pot.Pot)
	}
	r.Broadcast(&protocol.Pot{Pot: ps}, true)
	return
}

func (r *Room) action(pos uint8) {
	if r.allin+1 >= r.remain {
		return
	}
	var skip uint8
	if pos == 0 { // start from left hand of button
		pos = (r.Button)%r.Cap() + 1
	}

	for {
		var raised uint8
		r.Each(pos-1, func(o *Occupant) bool {
			if r.remain <= 1 {
				return false
			}
			if o.Pos == skip || o.chips == 0 {
				return true
			}
			r.WriteMsg(&protocol.BetPrompt{})
			n := o.GetAction(r.Timeout)
			if r.remain <= 1 {
				return false
			}
			if r.betting(o, n) {
				raised = o.Pos
				return false
			}
			return true
		})
		if raised == 0 {
			break
		}
		pos = raised
		skip = pos
	}
}

func (r *Room) ready() {
	r.Bet = 0
	r.Each(0, func(o *Occupant) bool {
		o.Bet = 0
		o.waitAction = false
		r.remain++
		o.HandValue = 0
		return true
	})
}

// 比牌
func (r *Room) showdown() {
	pots := r.calc()

	for i, _ := range r.Chips {
		r.Chips[i] = 0
	}

	for _, pot := range pots {
		var maxO *Occupant
		for _, pos := range pot.OPos {
			o := r.Occupants[pos-1]
			if o != nil && len(o.cards) > 0 {
				if maxO == nil {
					maxO = o
					continue
				}
				if o.HandValue > maxO.HandValue {
					maxO = o
				}
			}
		}

		var winners []uint8

		for _, pos := range pot.OPos {
			o := r.Occupants[pos-1]
			if o != nil && o.HandValue == maxO.HandValue && o.IsGameing() {
				winners = append(winners, o.Pos)
			}
		}

		if len(winners) == 0 {
			glog.Errorln("!!!no winners!!!")
			return
		}

		for _, winner := range winners {
			r.Chips[winner-1] += pot.Pot / uint32(len(winners))
		}
		r.Chips[winners[0]-1] += pot.Pot % uint32(len(winners)) // odd chips
	}

	for i, _ := range r.Chips {
		if r.Occupants[i] != nil {
			r.Occupants[i].chips += r.Chips[i]
		}
	}
}

func (r *Room) betting(o *Occupant, n int32) (raised bool) {
	if n > int32(o.chips) || // 手上筹码不足
		(n == 0 && o.Bet != r.Bet) || // 让牌
		(n > 0 && n != int32(o.chips) && ((n + int32(o.Bet)) < int32(r.Bet))) {
		glog.Errorf("下注筹码不合法!!！ n:%d  p.Bet:%d  p.Chips:%d  t.Bet:%d", n, o.Bet, o.chips, r.Bet)
		return
	}

	value := n
	actionName := ""
	if n < 0 {
		actionName = model.BET_FOLD
		n = 0
		r.remain--
		o.SetSitdown()
	} else if n == 0 {
		actionName = model.BET_CHECK
	} else if uint32(n)+o.Bet <= r.Bet {
		actionName = model.BET_CALL
		o.chips -= uint32(n)
		o.Bet += uint32(n)
	} else {
		actionName = model.BET_RAISE
		o.chips -= uint32(n)
		o.Bet += uint32(n)
		r.Bet = o.Bet
		raised = true
	}
	if o.chips == 0 {
		r.allin++
		actionName = model.BET_ALLIN
	}
	r.Chips[o.Pos-1] += uint32(n)

	r.Broadcast(&protocol.BetBroadcast{
		Uid:   o.Uid,
		Kind:  actionName,
		Value: value,
	}, true)

	return
}

func (r *Room) next(pos uint8) *Occupant {
	volume := r.Cap()
	for i := (pos) % volume; i != pos-1; i = (i + 1) % volume {
		if r.Occupants[i] != nil && r.Occupants[i].IsGameing() {
			return r.Occupants[i]
		}
	}
	return nil
}
