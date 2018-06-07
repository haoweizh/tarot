package function

import (
	"time"
	"fmt"
	"tarot/model"
	"tarot/util"
	"tarot/plugins/tarot"
	"tarot/wechat-go/wxweb"
)

func sendHandler(event *model.TarotEvent) {
	util.Info(`try to send event ` + event.NickName + event.SentenceType)
	if event == nil {
		return
	}
	if event.FromTarotStatus != event.ToTarotStatus {
		util.Info(`save db event ` + event.NickName + event.SentenceType)
		util.Info(fmt.Sprintf(`%s update status from %d to %d`, event.NickName, event.FromTarotStatus,
			event.ToTarotStatus))
		model.DB.Model(&model.MyContact{}).Where(`nick_name=?`, event.NickName).
			Updates(map[string]interface{}{`tarot_status`: event.ToTarotStatus, `updated_at`: time.Now()})
	}
	bytes := []byte(event.ToUserName)
	if bytes[0] == '@' && bytes[1] == '@' { //过滤掉@@开头的userName(微信群)
		return
	}
	util.SendTarotMsg(event.NickName, event.FromUserName, event.ToUserName, event.SentenceType, event.FromTarotStatus)
	model.ApplicationEvents.RemoveEvent(event.NickName)
}

func PlayTarot() {
	for true {
		rows, err := model.DB.Table("my_contacts").Select("nick_name,tarot_status,updated_at").
			Where("tarot_nick_name = ?", model.AppBot.Bot.NickName).Rows()
		if err != nil {
			util.Notice(`db error:` + fmt.Sprint(err))
			continue
		}
		for rows.Next() {
			var nickName string
			var tarotStatus int
			var updatedAt time.Time
			rows.Scan(&nickName, &tarotStatus, &updatedAt)
			toTarotStatus := tarot.CheckTime(tarotStatus, updatedAt)
			if toTarotStatus == 0 {
				continue
			}
			contacts := model.AppBot.Cm.GetContactsByName(nickName)
			for _, value := range contacts {
				sentenceType := fmt.Sprintf(`%d-%d`, tarotStatus, toTarotStatus)
				event := model.TarotEvent{FromUserName: model.AppBot.Bot.UserName, ToUserName: value.UserName,
					SentenceType: sentenceType, NickName: nickName, FromTarotStatus: tarotStatus, ToTarotStatus: toTarotStatus}
				util.Info(fmt.Sprintf("%s play tarot with sentence %s", nickName, sentenceType))
				model.SendChannel <- event
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func SendChannelServe() {
	for true {
		event := <-model.SendChannel
		util.Info(fmt.Sprintf(`get event of %s from channel with %s from %d to %d`,
			event.NickName, event.SentenceType, event.FromTarotStatus, event.ToTarotStatus))
		currentEvent := model.ApplicationEvents.GetEvent(event.NickName)
		if currentEvent != nil {
			util.SocketInfo(event.NickName + `old event ` + currentEvent.SentenceType + ` abandon ` + event.SentenceType)
			continue
		}
		model.ApplicationEvents.PutEvent(event.NickName, &event)
		go sendHandler(&event)
	}
}

func VerifyChannelServe()  {
	for true {
		msg := <-model.VerifyChannel
		util.Info(`get msg_fv` + msg.RecommendInfo.NickName)
		err := model.AppBot.AcceptFriend("", []*wxweb.VerifyUser{{Value: msg.RecommendInfo.UserName,
			VerifyUserTicket: msg.RecommendInfo.Ticket}})
		if err != nil {
			util.Notice(`fail to accept friend verification ` +err.Error())
			time.Sleep(time.Second * 3600)
		}
		model.AppBot.Cm.AddUser(&wxweb.User{NickName: msg.RecommendInfo.NickName,
			UserName: msg.RecommendInfo.UserName, City: msg.RecommendInfo.City, Sex: msg.RecommendInfo.Sex})
		myContact := model.MyContact{NickName: msg.RecommendInfo.NickName, TarotNickName: model.AppBot.Bot.NickName}
		model.DB.Where("nick_name = ? AND tarot_nick_name = ?", myContact.NickName, model.AppBot.Bot.NickName).
			First(&myContact)
		util.Info(fmt.Sprintf("accept user apply with name of %s", myContact.NickName))
		myContact.TarotStatus = 1
		if model.DB.NewRecord(&myContact) {
			util.Info(fmt.Sprintf("new contact added %s of %s", myContact.NickName, model.AppBot.Bot.NickName))
			model.DB.Create(&myContact)
		} else {
			model.DB.Save(&myContact)
		}
		event := model.TarotEvent{FromUserName: model.AppBot.Bot.UserName, ToUserName: msg.RecommendInfo.UserName,
			SentenceType: `1-101`, NickName: msg.RecommendInfo.NickName, FromTarotStatus: 1, ToTarotStatus: 101}
		model.SendChannel <- event
		time.Sleep(30*time.Second)
	}
}
