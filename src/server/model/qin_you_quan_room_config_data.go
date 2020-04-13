package model

import (
	"pdk/src/server/lib/db"
	"time"
)

type QinYouQianRoomConfig struct {
	ID        uint32 `gorm:"primary_key;index;type:BIGINT AUTO_INCREMENT"`
	Qid 	  uint32 `gorm:"primary_key;index;type:MEDIUMINT"`
	Uid 	  uint32 `gorm:"primary_key;index;type:MEDIUMINT"`
	QName		string `gorm:"type:VARCHAR(16)"`
	UName    	string `gorm:"type:VARCHAR(16)"`
	Config		string `gorm:"type:VARCHAR(512)"`
	CreatedAt time.Time
	UpdatedAt time.Time

}

func (q *QinYouQianRoomConfig) FindAllByQid()(* []QinYouQianRoomConfig,error)  {
	var qf []QinYouQianRoomConfig
	dbResult := db.GetGormDB().Where("qid = ? ",q.Qid).Find(&qf)
	return &qf,dbResult.Error
}




