package tarot

import (
	"fmt"
	"strings"
	"tarot/model"
	"tarot/util"
	"tarot/wechat-go/wxweb"
	"time"
)

// register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(listenCmd), "tarotText")
	session.HandlerRegister.Add(wxweb.MSG_IMG, wxweb.Handler(listenCmd), "tarotImg")
	session.HandlerRegister.Add(wxweb.MSG_SYS, wxweb.Handler(listenCmd), "tarotSys")
	session.HandlerRegister.Add(wxweb.MSG_EMOTION, wxweb.Handler(listenCmd), "tarotEmotion")
	session.HandlerRegister.Add(wxweb.MSG_LINK, wxweb.Handler(listenCmd), "tarotLink")
	session.HandlerRegister.Add(wxweb.MSG_SHORT_VIDEO, wxweb.Handler(listenCmd), "tarotShortVideo")
	session.HandlerRegister.Add(wxweb.MSG_LOCATION, wxweb.Handler(listenCmd), "tarotLocation")
	session.HandlerRegister.Add(wxweb.MSG_VOICE, wxweb.Handler(listenCmd), "tarotVoice")
	session.HandlerRegister.Add(wxweb.MSG_VIDEO, wxweb.Handler(listenCmd), "tarotVideo")
	session.HandlerRegister.Add(wxweb.MSG_FV, wxweb.Handler(listenCmd), "system-fv")
	if err := session.HandlerRegister.EnableByName("tarotText"); err != nil {
		util.Notice(err.Error())
	}
	if err := session.HandlerRegister.EnableByName("tarotImg"); err != nil {
		util.Notice(err.Error())
	}
	if err := session.HandlerRegister.EnableByName("tarotSys"); err != nil {
		util.Notice(err.Error())
	}
	if err := session.HandlerRegister.EnableByName("tarotEmotion"); err != nil {
		util.Notice(err.Error())
	}
	if err := session.HandlerRegister.EnableByName("tarotLink"); err != nil {
		util.Notice(err.Error())
	}
	if err := session.HandlerRegister.EnableByName("tarotShortVideo"); err != nil {
		util.Notice(err.Error())
	}
	if err := session.HandlerRegister.EnableByName("tarotLocation"); err != nil {
		util.Notice(err.Error())
	}
	if err := session.HandlerRegister.EnableByName("tarotVoice"); err != nil {
		util.Notice(err.Error())
	}
	if err := session.HandlerRegister.EnableByName("tarotVideo"); err != nil {
		util.Notice(err.Error())
	}
	if err := session.HandlerRegister.EnableByName("system-fv"); err != nil {
		util.Notice(err.Error())
	}
}

func listenCmd(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	// contact filter
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	switch msg.MsgType {
	case wxweb.MSG_FV:
		model.VerifyChannel <- *msg
		return
	case wxweb.MSG_TEXT:
		if strings.Contains(msg.Content, "唧唧复唧唧") {
			model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
				contact.NickName, session.Bot.NickName).
				Update(map[string]interface{}{"tarot_status": 101})
			return
		}
	case wxweb.MSG_SYS:
		if strings.Contains(util.FilterByte(msg.Content, '\n'), "对方验证通过后") {
			if msg.RecommendInfo != nil {
				model.DB.Model(&model.MyContact{}).Where(`nick_name=?`, msg.RecommendInfo.NickName).
					Updates(map[string]interface{}{`tarot_status`: 1, `updated_at`: time.Now()})
			} else {
				util.Info(`msg do not have recommend Info`)
			}
			return
		}
	}
	if contact == nil {
		util.Notice(`nil contact`)
		return
	}
	var toTarotStatus = 0
	var sentenceType string
	var myContact model.MyContact
	model.DB.Where("nick_name = ? AND tarot_nick_name = ?", contact.NickName, session.Bot.NickName).First(&myContact)
	if &myContact == nil {
		util.Notice(`nil my_contact`)
		return
	}
	if msg.MsgType == wxweb.MSG_SYS {
		if myContact.TarotStatus <= 201 {
			util.Info(`ignore sys message before tarot status 201`)
			return
		} else if !(myContact.TarotStatus == 401 || myContact.TarotStatus == 402 || myContact.TarotStatus == 403 ||
			myContact.TarotStatus == 404 || myContact.TarotStatus == 504) {
			sentenceType = `All除※-All除※`
			toTarotStatus = myContact.TarotStatus
			event := model.TarotEvent{FromUserName: session.Bot.UserName, ToUserName: contact.UserName,
				SentenceType: sentenceType, NickName: contact.NickName, FromTarotStatus: myContact.TarotStatus,
				ToTarotStatus: toTarotStatus}
			util.Info(fmt.Sprintf("%s receive hongbao %s %d %s", contact.NickName, sentenceType, msg.MsgType, msg.Content))
			model.SendChannel <- event
			return
		}
	}
	if (myContact.TarotStatus >= 101 && myContact.TarotStatus <= 104) ||
		(myContact.TarotStatus >= 500 && myContact.TarotStatus <= 503) ||
		(myContact.TarotStatus >= 510 && myContact.TarotStatus <= 515) ||
		(myContact.TarotStatus >= 520 && myContact.TarotStatus <= 525) ||
		(myContact.TarotStatus >= 530 && myContact.TarotStatus <= 533) || myContact.TarotStatus == 584 ||
		myContact.TarotStatus == 585 || myContact.TarotStatus == 594 || myContact.TarotStatus == 595 ||
		(myContact.TarotStatus >= 600 && myContact.TarotStatus <= 602) {
		toTarotStatus = receiveAny(myContact.TarotStatus)
	} else if (myContact.TarotStatus >= 200 && myContact.TarotStatus <= 211) || myContact.TarotStatus == 603 {
		toTarotStatus = receiveCheckImg(myContact.TarotStatus, msg.MsgType)
	} else if myContact.TarotStatus == 212 {
		toTarotStatus = doNothing(myContact.TarotStatus)
	} else if (myContact.TarotStatus >= 301 && myContact.TarotStatus <= 313) || myContact.TarotStatus == 604 {
		toTarotStatus = checkNum(myContact.TarotStatus, msg.Content, msg.MsgType)
	} else if myContact.TarotStatus >= 401 && myContact.TarotStatus <= 404 {
		toTarotStatus = receiveHongbao(myContact.TarotStatus, msg.MsgType)
	} else if myContact.TarotStatus == 505 && msg.MsgType == wxweb.MSG_TEXT {
		toTarotStatus = receiveBeginTarot(myContact.TarotStatus, msg.Content)
	} else if myContact.TarotStatus == 504 {
		toTarotStatus = receiveHongbao(myContact.TarotStatus, msg.MsgType)
	}
	sentenceType = fmt.Sprintf(`%d-%d`, myContact.TarotStatus, toTarotStatus)
	tarotLog := &model.TarotLog{TarotNickName: session.Bot.NickName, UserNickName: myContact.NickName,
		MsgType: msg.MsgType, MsgContent: msg.Content, FromStatus: myContact.TarotStatus, ToStatus: toTarotStatus}
	model.DB.Save(tarotLog)
	event := model.TarotEvent{FromUserName: session.Bot.UserName, ToUserName: contact.UserName,
		SentenceType: sentenceType, NickName: contact.NickName, FromTarotStatus: myContact.TarotStatus, ToTarotStatus: toTarotStatus}
	util.Info(fmt.Sprintf("%s play tarot with sentence %s %d %s", contact.NickName, sentenceType, msg.MsgType, msg.Content))
	model.SendChannel <- event
}
