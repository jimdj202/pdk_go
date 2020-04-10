package model

import (
	"pdk/src/server/lib/db"
	"testing"
	//"github.com/dolotech/lib/db"
)

func init() {
	//db.Init("postgres://postgres:haosql@127.0.0.1:5432/postgres?sslmode=disable")
	db.Init("pdk:WKwcyf66fTFKtip4@tcp(192.168.176.128:3306)/pdk?charset=utf8&parseTime=True&loc=Local")
}

func TestUser_UpdateChips(t *testing.T) {
	room := &Room{
	}

	t.Log(room.Insert())

	room = &Room{Rid: 5}


	id,err:= room.GetById()

	t.Log(room.CreatedAt)
	t.Logf("%v %v %#+v",id,err, room)
}
