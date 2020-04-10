package model

import "time"

type QinYouQuanMember struct {
	ID        uint32 `gorm:"primary_key;AUTO_INCREMENT;type:BIGINT"`
	CreatedAt time.Time
	Qid		uint32 `gorm:"primary_key;index;type:MEDIUMINT"`
	Uid		uint32 `gorm:"primary_key;index;type:BIGINT"`
	Qname   string `gorm:"type:VARCHAR(16)"`
}

