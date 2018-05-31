package main

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"tarot/model"
	"tarot/wechat-go/wxweb"
	"tarot/plugins/tarot"
	"tarot/function"
)

func main() {
	var err error
	// create session
	model.AppBot, err = wxweb.CreateSession(nil, nil, wxweb.WEB_MODE)
	if err != nil {
		logs.Error(err)
		return
	}
	tarot.Register(model.AppBot)
	model.DB, err = gorm.Open("postgres", model.DBConnection)
	if err != nil {
		logs.Warn(err)
		return
	}
	defer model.DB.Close()
	model.DB.AutoMigrate(&model.MyContact{})
	model.DB.AutoMigrate(&model.TarotSentence{})
	model.DB.AutoMigrate(&model.TarotLog{})
	model.ApplicationEvents = model.NewEvents()

	go function.ProcessLogin()
	model.AppBot.SetAfterLogin(func() (err error) {
		go function.SendChannelServe()
		go function.PlayTarot()
		function.SyncContact()
		return nil
	})
	select {}
}
