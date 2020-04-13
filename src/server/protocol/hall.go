package protocol

// 版本号
type Version struct {
	Version string
}

type CodeState struct {
	Code    int    // const
	Message string //警告信息
}


type Hello struct {
	Name string
}

//登录
type UserLoginInfo struct {
	UnionId  string
	Nickname string
}

//登录
type UserLoginInfoResp struct {
	UnionId  string
	Uid      uint32 // 用户id
	Account  string // 客户端玩家展示的账号
	Nickname string // 微信昵称
	Sex      uint8  // 微信性别 0-未知，1-男，2-女
	Profile  string // 微信头像
	Chips    uint32 // 筹码
}

type CreateRoom struct {
	ToomType uint32 `房间类型:跑得快,四川麻将等等`
	IsQinYouQuan bool
	QinYouQuanNum string
	CreatePdkRoomConf

}

type CreateRoomResp struct {
	TotalPersion uint32
	RoomNum string
}

type CreateQinYouQuan struct {
	Name string
	Uid uint32

}

type DeleteQinYouQuan struct {
	Qid uint32
	Uid uint32
}

func init() {
	Processor.Register(&Hello{})
	Processor.Register(&UserLoginInfo{})
	Processor.Register(&UserLoginInfoResp{})

	Processor.Register(&CodeState{})
	Processor.Register(&Version{})

}