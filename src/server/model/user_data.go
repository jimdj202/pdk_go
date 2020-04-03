package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"math/rand"
	"pdk/src/server/lib/db"
	"time"
)

func (this *User) GetByUId() (bool, error) {
	//return db.C().Engine().Where("uid = ?", this.Uid).Get(this)
	db := db.GetGormDB().Find(this)
	return db.RowsAffected > 0,nil
}


func (this *User) GetByAccount() (bool, error) {
	//return db.C().Engine().Where("account = ?", this.Account).Get(this)
	db := db.GetGormDB().Where("account = ?",this.Account).Find(this)
	return db.RowsAffected > 0,db.Error
}

func (this *User) GetByUnionId() (bool, error) {
	//return db.C().Engine().Where("union_id = ?", this.UnionId).Get(this)
	db := db.GetGormDB().Where("union_id = ?",this.UnionId).Find(this)
	return db.RowsAffected > 0, nil
}

type User struct {
	Uid        uint32    `gorm:"column:uid;type:BIGINT;primary_key;unique_index;AUTO_INCREMENT"`            // 用户id
	Account    string    `gorm:"column:account;type:VARCHAR(16);index;unique_index"` // 客户端玩家展示的账号
	DeviceId   string    `gorm:"column:device_id;type:VARCHAR(32);index"`             // 设备id
	UnionId    string    `gorm:"column:union_id;type:VARCHAR(32)"`              // 微信联合id
	Nickname   string    `gorm:"column:nickname;type:VARCHAR(32)"`              // 微信昵称
	Sex        uint8     `gorm:"column:sex;type:SMALLINT"`                      // 微信性别 0-未知，1-男，2-女
	Profile    string    `gorm:"column:profile;type:VARCHAR(128)"`               // 微信头像
	PhoneNum    string    `gorm:"column:phone_num;type:VARCHAR(20)"`               // 微信头像
	InviteCode string    `gorm:"column:invite_code;type:VARCHAR(6)"`             // 绑定的邀请码
	Coin       uint32    `gorm:"column:coin;type:INT"`                              // 筹码
	Lv         uint8     `gorm:"column:lv;type:SMALLINT"`                       // 等级
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`         // 注册时间
	LastTime   time.Time `gorm:"column:last_time"`                         // 上次登录时间
	LastIp     string    `gorm:"column:last_ip;type:VARCHAR(20)"`                    // 最后登录ip
	Kind       uint8     `gorm:"column:kind;type:SMALLINT;NOT NULL"`           // 用户类型
	Disable    bool      `gorm:"column:disable;type:SMALLINT"`                           // 是否禁用
	Signature  string    `gorm:"column:signature;type:VARCHAR(64)"`             // 个性签名
	Gps        string    `gorm:"column:gps;type:VARCHAR(32)"`                   // gps定位数据
	Black      bool      `gorm:"column:black;type:SMALLINT"`                             // 黑名单列表
	RoomID     string    `gorm:"column:room_id;type:VARCHAR(12)"`                           // 当前所在房间号，0表示不在房间,用于掉线重连
}

func (u *User) Create() error {

	n := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(999999-100000) + 100000
	u.Account = fmt.Sprintf("%v", n)
	u.Uid = uint32(n)
	for{
		exist,err:= u.GetByUId()
		if err!= nil {
			return err
		}
		if !exist {
			break
		}
		n := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(999999-100000) + 100000
		u.Account = fmt.Sprintf("%v", n)
	}

	db := db.GetGormDB().Create(u)
	if db.Error != nil {
		glog.Errorln(db.Error)
		return db.Error
	}

	return nil
}
func (u *User) UpdateChips(value int32) error {
	//_, err := db.C().Engine().Exec(`UPDATE public.user SET
	//	chips = chips + $1 WHERE uid =$2 `, value, u.Uid)

	db := db.GetGormDB().Model(u).Update("chips",gorm.Expr("chips + ?", value))
	if db.Error != nil {
		glog.Errorln(db.Error)
	}
	return db.Error
	//s:=db.C().Engine().Table(u).Incr("chips",value)
	//return nil
}

func (u *User) UpdateLoginTimeIp(ip string) error {
	//sql := `UPDATE public.user SET
	//last_time =  $1 ,last_ip =  $2 WHERE uid = $3 `
	//_, err := db.C().Engine().Exec(sql, time.Now(), utils.InetToaton(ip), u.Uid)
	//if err != nil {
	//	glog.Errorln(err)
	//}
	db := db.GetGormDB().Model(u).Update(User{LastTime:time.Now(),LastIp:ip})
	return db.Error
}

func (u *User) UpdateRoomId() error {
	//sql := `UPDATE public.user SET
	//room_id =  $1 WHERE uid = $2 `
	//_, err := db.C().Engine().Exec(sql, u.RoomID, u.Uid)
	//if err != nil {
	//	glog.Errorln(err)
	//}
	db := db.GetGormDB().Model(u).Update("room_id",u.RoomID)
	return db.Error
}
