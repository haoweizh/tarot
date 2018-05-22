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

var SendChannel = make(chan TarotEvent, 10)

const DBConnection = "host=139.196.84.85 port=5432 user=postgres dbname=postgres password=tarotpostgres sslmode=disable"

//const DBConnection = "host=172.20.48.174 port=5432 user=crook dbname=tarot password=zaq12WSX sslmode=disable"
