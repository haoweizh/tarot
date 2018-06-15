package model

import (
	"github.com/jinzhu/gorm"
	"tarot/wechat-go/wxweb"
)

//TODO
// 可能新加好友无法测试

//1 通过好友 3
var AppBot *wxweb.Session
var DB *gorm.DB
var ApplicationEvents *Events

var SendChannel = make(chan TarotEvent, 100)
var VerifyChannel = make(chan wxweb.ReceivedMessage, 500)

//const DBConnection = "host=139.196.84.85 port=5432 user=postgres dbname=postgres password=tarotpostgres sslmode=disable"

const DBConnection = "host=simonisgood.cfxtvaff0nlc.ap-northeast-1.rds.amazonaws.com port=54322 user=crook dbname=andrew password=zaq12WSX sslmode=disable"
