package util

import (
	"regexp"
	"strconv"
	"tarot/model"
	"time"
	"fmt"
	"math/rand"
	"github.com/bitly/go-simplejson"
	"strings"
	"hello/util"
	"github.com/pkg/errors"
)

func SendTarotMsg(from, to string, sentenceType string) {
	content, err := getSentence(sentenceType)
	if err != nil {
		util.Notice(err.Error())
		return
	}
	j, err := simplejson.NewJson([]byte(content))
	if err != nil {
		util.Notice(err.Error())
		return
	}
	sentences, _ := j.Get(`data`).StringArray()
	regNum := regexp.MustCompile(`\d+`)
	for _, value := range sentences {
		if strings.Contains(value, `tarotsleep`) {
			value = strings.Replace(value, `\n`, `                    `, -1)
			sleepSeconds := regNum.FindAllString(value, -1)
			if len(sleepSeconds) == 1 {
				sleepTime, _ := strconv.ParseInt(sleepSeconds[0], 10, 64)
				time.Sleep(time.Duration(sleepTime) * time.Second)
			} else if len(sleepSeconds) == 2 {
				right, _ := strconv.ParseInt(sleepSeconds[1], 10, 64)
				left, _ := strconv.ParseInt(sleepSeconds[0], 10, 64)
				time.Sleep(time.Duration(left+rand.Int63n(right-left)) * time.Second)
			} else {
				util.Notice(`wrong format for sleep from content: ` + value)
			}
		} else if strings.Contains(value, `tarotfile`) {
			bytes := []byte(value)
			path := bytes[5:]
			model.AppBot.SendFile(`./resource/`+string(path), from, to)
		} else if strings.Contains(value, `tarotjump`) {
			// do nothing
		} else {
			model.AppBot.SendText(value, from, to)
		}
	}
}

func getSentence(sentenceType string) (sentence string, err error) {
	rows, err := model.DB.Table("tarot_sentences").Select("content").
		Where("sentence_type = ?", sentenceType).Rows()
	if err != nil {
		return ``, errors.New(fmt.Sprintf("can not get row from %s", sentenceType))
	}
	contents := make([]string, 0)
	for rows.Next() {
		var content string
		rows.Scan(&content)
		contents = append(contents, content)
	}
	if len(contents) == 0 {
		return ``, errors.New(`no sentence found for ` + sentenceType)
	}
	return contents[rand.Intn(len(contents))], nil
}
