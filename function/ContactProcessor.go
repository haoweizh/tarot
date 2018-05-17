package function

import (
	"time"
	"fmt"
	"tarot/model"
	"math/rand"
	"tarot/util"
)

// 主动
//3 话术，表情 3/6(3次，间隔2天）
//4 请求分享 8
//7 发洗牌消息，选牌视频 10
//9 发结果语句和视频，求红包 12
//11 感谢、表情、金句 13
//12 4天后或周末 3设置welcome_no_resp为0
//13 四天后 3
func sendHandler(event model.TarotEvent) {
	contacts := model.AppBot.Cm.GetContactsByName(event.NickName)
	if contacts == nil || len(contacts) == 0 {
		fmt.Println("no this contact, ignore" + event.NickName)
		return
	}
	for _, value := range contacts {
		switch event.TarotStatus {
		case 3:
			if event.WelcomeNoResp >= 3 {
				model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
					value.NickName, model.AppBot.Bot.NickName).
					Update(map[string]interface{}{"tarot_status": 6, "new_status": true})
			} else if event.WelcomeNoResp == 0 || time.Now().Unix()-event.UpdatedAt.Unix() > 172800 {
				event.WelcomeNoResp++
				util.SendTarotText(`welcome`, 0, model.AppBot.Bot.UserName, value.UserName)
				model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
					value.NickName, model.AppBot.Bot.NickName).
					Update(map[string]interface{}{"welcome_no_resp": event.WelcomeNoResp, "new_status": true})
			} else {
				fmt.Println(fmt.Sprintf("%s ignore tarot status 3 with welcome no response times %d",
					event.NickName, event.WelcomeNoResp))
			}
		case 4:
			util.SendTarotText(`share`, 0, model.AppBot.Bot.UserName, value.UserName)
			model.AppBot.SendFile("./resource/share.jpg", model.AppBot.Bot.UserName, value.UserName)
			model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
				value.NickName, model.AppBot.Bot.NickName).
				Update(map[string]interface{}{"tarot_status": 8, "new_status": true})
		case 7:
			util.SendTarotText(`wash`, 0, model.AppBot.Bot.UserName, value.UserName)
			if rand.Intn(2) == 0 {
				model.AppBot.SendFile("./resource/00.mp4", model.AppBot.Bot.UserName, value.UserName)
			} else {
				model.AppBot.SendFile("./resource/01.mp4", model.AppBot.Bot.UserName, value.UserName)
			}
			util.SendTarotText(`choose`, 0, model.AppBot.Bot.UserName, value.UserName)
			model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
				value.NickName, model.AppBot.Bot.NickName).
				Update(map[string]interface{}{"tarot_status": 10, "new_status": true})
		case 9:
			index := rand.Intn(44) + 1
			path := fmt.Sprintf("./resource/%d.mp4", index)
			model.AppBot.SendFile(path, model.AppBot.Bot.UserName, value.UserName)
			util.SendTarotText(`answer`, index, model.AppBot.Bot.UserName, value.UserName)
			util.SendTarotText(`pay`, 0, model.AppBot.Bot.UserName, value.UserName)
			model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
				value.NickName, model.AppBot.Bot.NickName).
				Update(map[string]interface{}{"tarot_status": 12, "new_status": true})
		case 11:
			util.SendTarotText(`after_pay`, 0, model.AppBot.Bot.UserName, value.UserName)
			model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
				value.NickName, model.AppBot.Bot.NickName).
				Update(map[string]interface{}{"tarot_status": 13, "new_status": true})
		case 12:
			if time.Now().Unix()-event.UpdatedAt.Unix() > 345600 || time.Now().Weekday() == time.Saturday ||
				time.Now().Weekday() == time.Sunday {
				model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
					value.NickName, model.AppBot.Bot.NickName).
					Update(map[string]interface{}{"tarot_status": 3, "welcome_no_resp": 0, "new_status": true})
			}
		case 13:
			if time.Now().Unix()-event.UpdatedAt.Unix() > 345600 {
				model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
					value.NickName, model.AppBot.Bot.NickName).
					Update(map[string]interface{}{"tarot_status": 3, "welcome_no_resp": 0, "new_status": true})
			}
		}

		time.Sleep(time.Second * 3)
	}
}

func PlayTarot() {
	for true {
		rows, err := model.DB.Table("my_contacts").
			Select("city,sex,nick_name,tarot_status,updated_at,welcome_no_resp").
			Where("tarot_nick_name = ? AND tarot_status in (3,4,7,9,11,12,13) AND new_status = true",
			model.AppBot.Bot.NickName).Order("updated_at desc").Rows()
		if err != nil {
			fmt.Sprint(err)
			continue
		}
		for rows.Next() {
			var nickName, city string
			var tarotStatus, welcomeNoResp, sex int
			var updatedAt time.Time
			rows.Scan(&city, &sex, &nickName, &tarotStatus, &updatedAt, &welcomeNoResp)
			fmt.Println(fmt.Sprintf("%s at status %d play tarot no response time %d", nickName, tarotStatus, welcomeNoResp))
			event := model.TarotEvent{NickName: nickName, TarotStatus: tarotStatus, City: city, Sex: sex,
				UpdatedAt: updatedAt, WelcomeNoResp: welcomeNoResp}
			if tarotStatus != 3 {
				model.DB.Table("my_contacts").Where("nick_name = ? AND tarot_nick_name = ?",
					nickName, model.AppBot.Bot.NickName).Update(map[string]interface{}{"new_status": false})
			}
			model.SendChannel <- event
		}
		time.Sleep(time.Second * 10)
	}
}

func SendChannelServe() {
	for true {
		state := <-model.SendChannel
		fmt.Println(fmt.Sprintf("go to handle state %s status %d", state.NickName, state.TarotStatus))
		go sendHandler(state)
	}
}
