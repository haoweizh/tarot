package function

import (
	"github.com/songtianyi/rrframework/logs"
	"time"
	"tarot/wechat-go/wxweb"
	"tarot/model"
	"fmt"
)

func ProcessLogin() {
	for true {
		if err := model.AppBot.LoginAndServe(false); err != nil {
			logs.Error("session exit, %s", err)
			for i := 0; i < 3; i++ {
				logs.Info("trying re-login with cache")
				if err := model.AppBot.LoginAndServe(true); err != nil {
					logs.Error("re-login error, %s", err)
				}
				time.Sleep(3 * time.Second)
			}
			if model.AppBot, err = wxweb.CreateSession(nil, model.AppBot.HandlerRegister, wxweb.WEB_MODE); err != nil {
				logs.Error("create new session failed, %s", err)
				break
			}
		} else {
			logs.Info("closed by user")
			break
		}
	}
}

func SyncContact() {
	contacts := model.AppBot.Cm.GetAll()
	for _, contact := range contacts {
		myContact := &model.MyContact{NickName: contact.NickName, TarotNickName: model.AppBot.Bot.NickName}
		model.DB.Where("nick_name = ? AND tarot_nick_name = ?", myContact.NickName,
			model.AppBot.Bot.NickName).First(&myContact)
		myContact.Init(contact)
		if model.DB.NewRecord(&myContact) {
			myContact.TarotStatus = 3
			fmt.Println(myContact.NickName + " is added as nick of " + model.AppBot.Bot.NickName)
			model.DB.Create(&myContact)
		} else {
			fmt.Println(myContact.NickName + " do not updated as nick of " + model.AppBot.Bot.NickName)
			//model.DB.Save(myContact)
		}
	}
}
