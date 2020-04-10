package model

//亲友圈基本信息
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"math/rand"
	"pdk/src/server/lib/db"
	"time"
)

type QinYouQuan struct {
	//ID        uint32 `gorm:"primary_key;index;type:BIGINT AUTO_INCREMENT"`
	Qid 	  uint32 `gorm:"primary_key;index;type:MEDIUMINT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Name		string `gorm:"type:VARCHAR(16)"`


}

func (q *QinYouQuan)Create()(int64,error)  {
	n := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(999999-100000) + 100000
	q.Qid = uint32(n)
	for  {
		count,_:= q.FindOneByQid()
		if count  == 0 {
			break
		}
		q.Qid = uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(999999-100000) + 100000)
	}

	dbResult := db.GetGormDB().Create(q)
	return dbResult.RowsAffected,dbResult.Error
}

func (q *QinYouQuan) Delete() error {
	return db.GetGormDB().Transaction(func(tx *gorm.DB) error {
		// 在事务中做一些数据库操作 (这里应该使用 'tx' ，而不是 'db')
		dbResult := tx.Where("qid = ? ",q.Qid).Delete(QinYouQuanMember{})
		if err := dbResult.Error;err != nil {
			return err
		}
		dbResult2 := tx.Delete(q)
		if err := dbResult2.Error;err != nil {
			return err
		}
		//if err := tx.Delete(&Animal{Name: "Giraffe"}).Error; err != nil {
		//	// 返回任意 err ，整个事务都会 rollback
		//	return err
		//}
		//
		//if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
		//	return err
		//}

		// 返回 nil 提交事务
		return nil
	})
}

func (q *QinYouQuan) FindOneByQid()(int64,error)  {
	dbResult := db.GetGormDB().Where("qid = ?",q.Qid).First(q)
	return dbResult.RowsAffected,dbResult.Error
}
