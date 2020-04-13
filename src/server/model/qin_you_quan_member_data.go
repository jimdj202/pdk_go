package model

import (
	"pdk/src/server/lib/db"
	"time"
)

type QinYouQuanMember struct {
	ID        uint32 `gorm:"primary_key;type:BIGINT AUTO_INCREMENT"`
	CreatedAt time.Time
	Qid		uint32 `gorm:"primary_key;index;type:MEDIUMINT"`
	Uid		uint32 `gorm:"primary_key;index;type:BIGINT"`
	Qname   string `gorm:"type:VARCHAR(16)"`
	Uname	string `gorm:"type:VARCHAR(16)"`
	Status  uint8 `gorm:"type:SMALLINT"`
	IsManager uint8 `gorm:"type:SMALLINT"`

}

func (q *QinYouQuanMember) Create () (int64,error){
	dbResult := db.GetGormDB().Create(q)
	return dbResult.RowsAffected,dbResult.Error
}

func (q *QinYouQuanMember) Delete() (int64,error) {
	dbResult := db.GetGormDB().Delete(q)
	return dbResult.RowsAffected,dbResult.Error
}

func (q *QinYouQuanMember) DeleteByQidAndUid() (int64,error) {
	dbResult := db.GetGormDB().Where("qid = ? AND uid = ? ",q.Qid,q.Uid).Delete(QinYouQuanMember{})
	return dbResult.RowsAffected,dbResult.Error
}

func (q *QinYouQuanMember) GetMembersByQid()( *[]QinYouQuanMember,error)  {
	var mebs  []QinYouQuanMember
	dbResult := db.GetGormDB().Where("qid = ?",q.Qid).Find(&mebs)
	return &mebs,dbResult.Error
}
