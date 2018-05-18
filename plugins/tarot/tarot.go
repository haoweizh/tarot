package tarot

import (
	"tarot/wechat-go/wxweb"
	"strings"
	"tarot/model"
	"regexp"
	"strconv"
	"github.com/songtianyi/rrframework/logs"
	"math/rand"
	"time"
	"tarot/util"
)

// register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(listenCmd), "tarotText")
	session.HandlerRegister.Add(wxweb.MSG_IMG, wxweb.Handler(listenCmd), "tarotImg")
	session.HandlerRegister.Add(wxweb.MSG_SYS, wxweb.Handler(listenCmd), "tarotSys")
	if err := session.HandlerRegister.EnableByName("tarotText"); err != nil {
		logs.Error(err)
	}
	if err := session.HandlerRegister.EnableByName("tarotImg"); err != nil {
		logs.Error(err)
	}
}

//3 收到回应 4
//8 收到图片或随机概率 7
//8 收到非图片,抗议 8
//10 收到数字 9
//12 收到红包 11
func listenCmd(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	// contact filter
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("no this contact, ignore", msg.FromUserName)
		return
	}
	if msg.MsgType == wxweb.MSG_TEXT && strings.Contains(msg.Content, "重新") {
		model.AppBot.SendFile(`./resource/emoji1.png`, model.AppBot.Bot.UserName, contact.UserName)
		model.AppBot.SendFile(`./resource/cry.gif`, model.AppBot.Bot.UserName, contact.UserName)
		model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
			contact.NickName, session.Bot.NickName).
			Update(map[string]interface{}{"tarot_status": 3, "welcome_no_resp": 0, "new_status": true})
		return
	}

	var myContact model.MyContact
	model.DB.Where("nick_name = ? AND tarot_nick_name = ?",
		contact.NickName, session.Bot.NickName).First(&myContact)
	switch myContact.TarotStatus {
	case 3:
		model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
			contact.NickName, session.Bot.NickName).
			Update(map[string]interface{}{"tarot_status": 4, "welcome_no_resp": 0, "new_status": true})
	case 8:
		if msg.MsgType == wxweb.MSG_IMG || rand.Intn(5) == 1 { // 收到图片
			model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
				contact.NickName, model.AppBot.Bot.NickName).
				Update(map[string]interface{}{"tarot_status": 7, "new_status": true})
		} else {
			util.SendTarotText(`no_share`, 0, model.AppBot.Bot.NickName, contact.UserName)
			time.Sleep(10 * time.Second)
		}
	case 10:
		reg := regexp.MustCompile(`\d+`)
		choice, _ := strconv.ParseInt(reg.FindString(msg.Content), 10, 0)
		if choice >= 1 || choice <= 22 {
			model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
				contact.NickName, model.AppBot.Bot.NickName).
				Update(map[string]interface{}{"tarot_status": 9, "new_status": true})
		}
	case 12:
		if msg.MsgType == wxweb.MSG_SYS {
			model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
				contact.NickName, model.AppBot.Bot.NickName).
				Update(map[string]interface{}{"tarot_status": 11, "new_status": true})
		}
	}
}
