package function

import (
	"time"
	"fmt"
	"tarot/model"
	"tarot/util"
	"tarot/plugins/tarot"
)

func sendHandler(nickEvents map[string]*model.TarotEvent, event model.TarotEvent) {
	util.SendTarotMsg(event.FromUserName, event.ToUserName, event.SentenceType)
	if event.FromTarotStatus != 0 && event.ToTarotStatus != 0 && event.FromTarotStatus != event.ToTarotStatus {
		model.DB.Model(&model.MyContact{}).Where(`nick_name=?`, event.NickName).
			Updates(map[string]interface{}{`tarot_status`: event.ToTarotStatus, `updated_at`: time.Now()})
	}
	nickEvents[event.NickName] = nil
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

func checkEventAmount(events map[string]*model.TarotEvent) (amount int) {
	amount = 0
	for _, value := range events {
		if value != nil {
			amount++
		}
	}
	return amount
}

func SendChannelServe() {
	nickEvents := make(map[string]*model.TarotEvent)
	for true {
		event := <-model.SendChannel
		if nickEvents[event.NickName] != nil {
			util.SocketInfo(fmt.Sprintf(`can not send %s msg, abandon %s`, event.NickName, event.SentenceType))
			continue
		}
		nickEvents[event.NickName] = &event
		for checkEventAmount(nickEvents) > 2 {
			util.SocketInfo(fmt.Sprintf(`%d events in nickEvents, sleep 3 seconds`, len(nickEvents)))
			time.Sleep(time.Second * 3)
		}
		go sendHandler(nickEvents, event)
	}
}
