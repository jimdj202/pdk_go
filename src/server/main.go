package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"pdk/src/server/conf"
	"pdk/src/server/game"
	"pdk/src/server/gate"
	"pdk/src/server/lib/db"
	"pdk/src/server/login"
	"pdk/src/server/model"
)

func main() {
	//Init the command-line flags. for glog
	flag.Parse()
	defer glog.Flush()
	//flag.Lookup("alsologtostderr").Value.Set("true")
	//flag.Lookup("log_dir").Value.Set("./log/")

	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath
	conf.Server.DBUrl = "pdk:WKwcyf66fTFKtip4@tcp(192.168.176.128:3306)/pdk?charset=utf8&parseTime=True&loc=Local"
	db.Init(conf.Server.DBUrl)

	//createDb()
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}

func createDb() {
	// 建表,只维护和服务器game里面有关的表
	//err := db.GetGormDB().AutoMigrate(&model.User{}, &model.Room{},&model.QinYouQuan{},&model.QinYouQuanMember{})
	err := db.GetGormDB().AutoMigrate(&model.QinYouQuan{},&model.QinYouQuanMember{})
	if err != nil {
		glog.Errorln(err)
	}
}