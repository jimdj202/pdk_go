package model

//亲友圈基本信息
import (
	_ "github.com/go-sql-driver/mysql"
	"pdk/src/server/lib/db"
	"time"
)

type QinYouQuan struct {
	ID        uint32 `gorm:"primary_key;AUTO_INCREMENT;index;unique_index;type:BIGINT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Name		string `gorm:"type:VARCHAR(16)"`


}

func (q *QinYouQuan)Create()(int64,error)  {
	dbResult := db.GetGormDB().Create(q)
	return dbResult.RowsAffected,dbResult.Error
}

func (q *QinYouQuan) Delete()(int64,error) {
	dbResult := db.GetGormDB().Delete(q)
	return dbResult.RowsAffected,dbResult.Error
}


