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

