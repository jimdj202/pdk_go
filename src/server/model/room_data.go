package model

import (
	_ "github.com/go-sql-driver/mysql"
	"pdk/src/server/lib/db"
	"time"
)

//房间基本信息

type Room struct {
	Rid             uint32    `gorm:"column:rid;primary_key;AUTO_INCREMENT;index;unique_index;type:BIGINT"`
	Number          string    `gorm:"column:number;index;not null;type:VARCHAR(8)"` // 给玩家展示的房间号
	Pwd             string    `gorm:"column:pwd;type:VARCHAR(16)"`                  //房间锁--密码
	State           uint8     `gorm:"column:state;type:SMALLINT"`                   //房间状态 0默认可用 1不可用
	Name            string    `gorm:"column:name;type:VARCHAR(16)"`                 //房间名字
	CreatedAt       time.Time `gorm:"column:created_at;index;default:CURRENT_TIMESTAMP"`        //创建时间
	OriginalOwnerID uint32    `gorm:"column:original_owner_id;type:BIGINT;"`                //原始创建人的信息
	Owner           uint32    `gorm:"column:owner;type:BIGINT;"`                            //房管
	Kind            uint32    `gorm:"column:kind;type:SMALLINT;"`                             //游戏类型 即玩法
	DraginChips     uint32    `gorm:"column:dragin_chips;type:BIGINT;"`                     //带入筹码

	//Occupants       []*uint32 `xorm:"'occupants'"`                        // 玩家列表，列表第一项为庄家
}

func (u *Room) Insert() (int64, error) {
	db := db.GetGormDB().Create(u)
	return db.RowsAffected,db.Error
}

func (this *Room) GetById() (bool, error) {
	db := db.GetGormDB().Where("rid = ?", this.Rid).Find(this)
	return db.RowsAffected>0,db.Error
}


func (r *Room) CreatedTime() uint32 {
	return uint32(r.CreatedAt.Unix())
}