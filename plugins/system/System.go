package system

import (
	"github.com/songtianyi/rrframework/logs"
	"tarot/wechat-go/wxweb"
	"tarot/model"
)

func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_SYS, wxweb.Handler(system), "system-sys")
	session.HandlerRegister.Add(wxweb.MSG_WITHDRAW, wxweb.Handler(system), "system-withdraw")
	session.HandlerRegister.Add(wxweb.MSG_FV, wxweb.Handler(system), "system-fv")

	if err := session.HandlerRegister.EnableByName("system-sys"); err != nil {
		logs.Error(err)
	}

	if err := session.HandlerRegister.EnableByName("system-withdraw"); err != nil {
		logs.Error(err)
	}

	if err := session.HandlerRegister.EnableByName("system-fv"); err != nil {
		logs.Error(err)
	}
}

func system(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	switch msg.MsgType {
	case wxweb.MSG_FV:
		session.AcceptFriend("", []*wxweb.VerifyUser{{Value: msg.RecommendInfo.UserName,
			VerifyUserTicket: msg.RecommendInfo.Ticket}})
		myContact := model.MyContact{NickName: msg.RecommendInfo.NickName, TarotNickName: model.AppBot.Bot.NickName,
			NewStatus: true}
		model.DB.Where("nick_name = ? AND tarot_nick_name = ?", myContact.NickName, model.AppBot.Bot.NickName).
			First(&myContact)
		logs.Info("accept user apply with name of %s", myContact.NickName)
		if model.DB.NewRecord(&myContact) {
			myContact.TarotStatus = 3
			myContact.NewStatus = true
			myContact.WelcomeNoResp = 0
			logs.Info("new contact added %s of %s", myContact.NickName, model.AppBot.Bot.NickName)
			model.DB.Create(&myContact)
		} else {
			logs.Info("do not update contact user nick %s of %s", myContact.NickName, model.AppBot.Bot.NickName)
			//model.DB.Save(myContact)
		}
	}

	logs.Debug(msg)
}
