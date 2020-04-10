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
	qin := &QinYouQuan{Uid: 12345,
		Name: "跑得快5556"}
	qin.Create()
}

func TestQinYouQuan_FindOne(t *testing.T) {
	qin := &QinYouQuan{Qid: 12345}
	count, err := qin.FindOneByQid()
	fmt.Println("TestQinYouQuan_FindOne:",count,err,qin)
}

func TestQinYouQuan_FindAllByUid(t *testing.T) {
	qin := &QinYouQuan{Uid: 12345}
	list, err := qin.FindAllByUid()
	fmt.Println("TestQinYouQuan_FindOne:",list,err,qin)
}