package internal

import (
	"github.com/name5566/leaf/gate"
	"pdk/src/server/conf"
	"pdk/src/server/game"
	"pdk/src/server/protocol"
	"github.com/golang/glog"
)

type Module struct {
	*gate.Gate
}

func (m *Module) OnInit() {
	m.Gate = &gate.Gate{
		MaxConnNum:      conf.Server.MaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		WSAddr:          conf.Server.WSAddr,
		HTTPTimeout:     conf.HTTPTimeout,
		CertFile:        conf.Server.CertFile,
		KeyFile:         conf.Server.KeyFile,
		TCPAddr:         conf.Server.TCPAddr,
		LenMsgLen:       conf.LenMsgLen,
		LittleEndian:    conf.LittleEndian,
		Processor:       protocol.Processor,
		AgentChanRPC:    game.ChanRPC,
	}
}
func (gate *Module) OnDestroy() {
	glog.Errorln("gate OnDestroy")
}


