package model

import (
	"fmt"
	"pdk/src/server/lib/db"
	"testing"
)
func init() {
	db.Init("pdk:WKwcyf66fTFKtip4@tcp(192.168.176.128:3306)/pdk?charset=utf8&parseTime=True&loc=Local")
}

func TestQinYouQuan_Create(t *testing.T) {
	qin := &QinYouQuan{Qid: 123555,
		Name: "跑得快555"}
	qin.Create()
}

func TestQinYouQuan_FindOne(t *testing.T) {
	qin := &QinYouQuan{Qid: 12345}
	count, err := qin.FindOneByQid()
	fmt.Println("TestQinYouQuan_FindOne:",count,err,qin)
}
