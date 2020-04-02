// 数据库客户端常用操作 （基于 xorm）
package db

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

// Client 数据库客户端
type Client struct {
	dbGorm *gorm.DB
	dbSql *sql.DB
}

var db *Client

//var database_addr string
//const configFilePath string = "config.ini"
func Init(url string) error{
	var err error
	if db == nil {
		//db, err = NewClient("postgres", url)
		db, err = NewClient("mysql", url)
		//db.ShowSQL(config.GetCfgDatabase().ShowSql)
		db.ShowSQL(true)
	}
	return err
}

func GetGormDB() *gorm.DB{
	return db.dbGorm
}

// NewClient 创建一个客户端链接
func NewClient(driver, connstr string) (*Client, error) {
	//db, err := gorm.Open("mysql", "user:password@tcp(IP:port)/dbname?charset=utf8&parseTime=True&loc=Local")
	//db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	dbGorm, err := gorm.Open(driver, connstr)
	if err == nil {
		dbSql := dbGorm.DB()
		dbSql.SetMaxIdleConns(50) //最大连接数设置为200
		dbSql.SetMaxOpenConns(200)
		ping := func() <-chan string {
			errmsg := make(chan string)
			go func() {
				err := dbGorm.DB().Ping()
				if err != nil {
					errmsg <- err.Error()
				} else {
					errmsg <- ""
				}
			}()
			return errmsg
		}
		go func() {
			select {
			case msg := <-ping():
				if len(msg) > 0 {
					//logs.Error(msg)
				} else {
					//logs.Info("connect %s succeed.", driver)
				}
			case <-time.After(time.Second * 5):
				//logs.Info("connect timeout 5 second.")
			}
		}()
		return &Client{dbGorm,dbSql}, err
	}
	return nil,err
}

func (cl *Client) Engine() *gorm.DB {
	return cl.dbGorm
}

// Close 关闭连接
func (cl *Client) Close() error {
	if cl.dbGorm != nil {
		err := cl.dbGorm.Close()
		cl.dbGorm = nil
		return err
	}
	return nil
}


// ShowSQL 显示 SQL
func (cl *Client) ShowSQL(show bool) {
	cl.dbGorm.LogMode(show)
}

