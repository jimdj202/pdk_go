package gate

import (
	"pdk/src/server/protocol"
	"pdk/src/server/game"
	"pdk/src/server/login"
)

func init() {
	protocol.Processor.SetRouter(&protocol.UserLoginInfo{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Version{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.RoomList{}, login.ChanRPC)

	protocol.Processor.SetRouter(&protocol.CreateQinYouQuan{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.DeleteQinYouQuan{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.JoinQinYouQuan{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.LeaveQinYouQuan{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.GetAllQinYouQuan{}, login.ChanRPC)

	protocol.Processor.SetRouter(&protocol.CreateRoom{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.JoinRoom{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.LeaveRoom{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.SitDown{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.StandUp{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Bet{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.Chat{}, game.ChanRPC)
}
