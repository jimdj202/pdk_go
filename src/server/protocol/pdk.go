package protocol

type UserInfo struct {
	Uid      uint32 // 用户id
	Account  string // 客户端玩家展示的账号
	Nickname string // 微信昵称
	Sex      uint8  // 微信性别 0-未知，1-男，2-女
	Profile  string // 微信头像
	Chips    uint32 // 筹码
}

type RoomInfo struct {
	Number    string
	Volume    uint8
	GameType  uint32 //游戏类型 即玩法
	PayValue  uint8  //倍数
	BaseMoney uint32 //最低资本 才能进房间
	RoomPwd   string //房间锁--密码
	RoomID    uint32

	SB       uint32   // 小盲注
	BB       uint32   // 大盲注
	Cards    []byte   //公共牌
	Pot      []uint32 // 当前奖池筹码数
	Timeout  uint8    // 倒计时超时时间(秒)
	Button   uint8    // 当前庄家座位号，从1开始
	Chips    []uint32 // 玩家本局下注的总筹码数，与occupants一一对应
	Bet      uint32   // 当前下注额
	Max      uint8    // 房间最大玩家人数
	MaxChips uint32
	MinChips uint32
}

type CreateRoom struct {
	ToomType uint32 `房间类型:跑得快,四川麻将等等`
	IsQinYouQuan bool
	QinYouQuanNum string
	CreatePdkRoomConf

}

type CreatePdkRoomConf struct {
	TotalPersion uint32

}

type CreateRoomResp struct {
	TotalPersion uint32
	RoomNum string
}

type StandUp struct {
	Uid uint32
}

type SitDown struct {
	Uid uint32
	Pos uint8
}

type LeaveRoom struct {
	RoomNumber string
	Uid        uint32
}


type JoinRoom struct {
	Uid        uint32
	RoomNumber string
	RoomPwd    string
}

type JoinRoomBroadcast struct {
	UserInfo *UserInfo
}

type JoinRoomResp struct {
	UserInfos []*UserInfo
	RoomInfo  *RoomInfo
}

//底牌
type PreFlop struct {
	Cards []byte
	Kind  uint8
}

// 翻牌
type Flop struct {
	Cards []byte
	Kind  uint8
}

// 转牌
type Turn struct {
	Card byte
	Kind uint8
}

//河牌
type River struct {
	Card byte
	Kind uint8
}

//通报本局庄家
type Button struct {
	Uid uint32
}

// 玩家提交下注数据
// 有四种下注方式，下注数分别对应为：
//弃牌: <0 (fold)
//跟注：等于单注额 (call)
//看注：= 0 表示看注 (check)
//加注：大于单注额 (raise)
//全押：等于玩家手中所有筹码 (allin)
type Bet struct {
	Value int32
}

// 提示指定的玩家下注
type BetPrompt struct {
}

// 通报玩家下注
type BetBroadcast struct {
	Value int32
	Kind  string
	Uid   uint32
}

type BetResp struct {
	Value int32
	Kind  string
	Uid   uint32
}

//通报奖池
type Pot struct {
	Pot []uint32
}

//摊牌和比牌
type Showdown struct {
	Showdown []*ShowdownItem
}

type ShowdownItem struct {
	Uid      uint32
	ChipsWin uint32
	Chips    uint32
}

type RoomList struct {
}

type Room struct {
	Rid             uint32
	Number          string // 给玩家展示的房间号
	State           uint8  //房间状态 0默认可用 1不可用
	Name            string //房间名字
	CreatedAt       uint32 //创建时间
	OriginalOwnerID uint32 //原始创建人的信息
	Owner           uint32 //房管
	Kind            uint32 //游戏类型 即玩法
	DraginChips     uint32 // 带入筹码
	Cap             uint8
	MaxCap          uint8
}

type RoomListResp struct {
	Room []*Room
}

type Chat struct {
	Uid  uint32
	Text string
}


func init() {

	//房间会话注册
	Processor.Register(&RoomInfo{})  //基本信息
	Processor.Register(&JoinRoom{})  //
	Processor.Register(&LeaveRoom{}) //

	Processor.Register(&Showdown{})
	Processor.Register(&PreFlop{})
	Processor.Register(&Pot{})
	Processor.Register(&Bet{})
	Processor.Register(&Button{})
	Processor.Register(&StandUp{})
	Processor.Register(&SitDown{})
	Processor.Register(&UserInfo{})
	Processor.Register(&JoinRoomResp{})
	Processor.Register(&JoinRoomBroadcast{})
	Processor.Register(&BetResp{})
	Processor.Register(&RoomList{})
	Processor.Register(&RoomListResp{})
	Processor.Register(&Chat{})
}