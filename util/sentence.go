package util

import (
	"regexp"
	"strconv"
	"tarot/model"
	"time"
	"fmt"
	"math/rand"
)

// 一共有welcome share no_share wash choose pay after_pay answer 这几大类语句,每一大类叫做一个sentence_type(列名)
// 1. 每一类语句可以有多个版本，插入到数据库的不同行里就行，系统会随机从同一个大类中选一个发送，删除掉就不会发送这一句了
// 2. 发送的内容在content(列名)里，每一个content中加入 sleep数字，会顺序发送多条
//    例如  你好我是你的占卜师sleep5现在开始 这样就会发两条中间间隔5秒，可以支持加入无数个 sleep数字，放在一个content的
//    前面后面中间，都可以
// 3. 同一条信息需要分段的，内容里打 \n
type TarotTextMsg struct {
	Content string
	Sleep   int64
}

func parseWxMsg(content string) []TarotTextMsg {
	reg := regexp.MustCompile(`sleep\d+`)
	regNum := regexp.MustCompile(`\d+`)
	contents := reg.Split(content, -1)
	sleeps := reg.FindAllString(content, -1)
	messages := make([]TarotTextMsg, len(contents))
	for index, value := range contents {
		messages[index] = TarotTextMsg{}
		messages[index].Content = value
		if index < len(sleeps) {
			messages[index].Sleep, _ = strconv.ParseInt(regNum.FindString(sleeps[index]), 10, 64)
		} else if index == len(sleeps) {
			messages[index].Sleep = 0
		}
	}
	return messages
}

func SendTarotText(sentenceType string, sentenceScenario int, from, to string) {
	content := getSentence(sentenceType, sentenceScenario)
	sendTarotText(content, from, to)
}

func sendTarotText(content, from, to string) {
	tarotTextMessages := parseWxMsg(content)
	for _, value := range tarotTextMessages {
		if value.Content != `` {
			model.AppBot.SendText(value.Content, from, to)
		}
		time.Sleep(time.Second * time.Duration(value.Sleep))
	}
}

func getSentence(sentenceType string, sentenceScenario int) (sentence string) {
	rows, err := model.DB.Table("tarot_sentences").Select("content").
		Where("sentence_type = ? AND sentence_scenario = ?", sentenceType, sentenceScenario).Rows()
	if err != nil {
		fmt.Printf("can not get row from %s with scenario %d", sentenceType, sentenceScenario)
		return ``
	}
	contents := make([]string, 0)
	for rows.Next() {
		var content string
		rows.Scan(&content)
		contents = append(contents, content)
	}
	return contents[rand.Intn(len(contents))]
}
