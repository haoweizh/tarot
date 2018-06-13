package tarot

import (
	"regexp"
	"strconv"
	"strings"
	"tarot/wechat-go/wxweb"
	"math/rand"
	"time"
)

func triggerBySecondOfDay(updatedAt time.Time, waitSeconds int64, startSecond, endSecond int) bool {
	nowUnixSeconds := time.Now().Unix()
	if nowUnixSeconds-updatedAt.Unix() < waitSeconds {
		return false
	}
	nowSecond := time.Now().Hour()*3600 + time.Now().Minute()*60 + time.Now().Second()
	if rand.Intn(endSecond-startSecond+1) > nowSecond-startSecond {
		return false
	}
	return true
}

func triggerByWaitTime(updatedAt time.Time, waitSeconds, endSeconds int64) bool {
	nowUnixSeconds := time.Now().Unix()
	if nowUnixSeconds-updatedAt.Unix() < waitSeconds {
		return false
	}
	if endSeconds <= waitSeconds || nowUnixSeconds-updatedAt.Unix() > endSeconds {
		return true
	}
	if rand.Int63n(endSeconds-waitSeconds+1) > nowUnixSeconds-updatedAt.Unix() {
		return false
	}
	return true
}

func CheckTime(fromTarotStatus int, updatedAt time.Time) (toTarotStatus int) {
	// 进入该状态满30天，接下来的21点～24点间随机择时
	if triggerBySecondOfDay(updatedAt, 2592000, 75600, 86400) {
		if fromTarotStatus == 600 {
			return 601
		}
	}
	// 进入该状态满30天
	if triggerByWaitTime(updatedAt, 2592000, 2592000) {
		if fromTarotStatus == 601 {
			return 602
		}
	}
	// 进入该状态满10天，接下来的21点～24点间随机择时
	if triggerBySecondOfDay(updatedAt, 864000, 75600, 86400) {
		if fromTarotStatus >= 530 && fromTarotStatus <= 533 {
			return 600
		}
	}
	// 进入该状态满7天，接下来的21点～24点间随机择时
	if triggerBySecondOfDay(updatedAt, 604800, 75600, 86400) {
		if fromTarotStatus == 504 || fromTarotStatus==584 || fromTarotStatus == 594 {
			return 514
		}
		if fromTarotStatus == 505 || fromTarotStatus==585 || fromTarotStatus == 595 {
			return 515
		}
	}
	// 进入该状态满4天，接下来的21点～24点间随机择时
	if triggerBySecondOfDay(updatedAt, 345600, 75600, 86400) {
		if fromTarotStatus == 524 || fromTarotStatus == 525 {
			return 600
		}
	}
	// 进入该状态满3天，接下来的21点～24点间随机择时
	if triggerBySecondOfDay(updatedAt, 259200, 75600, 86400) {
		if (fromTarotStatus >= 520 && fromTarotStatus <= 523) || fromTarotStatus == 514 || fromTarotStatus == 515 {
			return fromTarotStatus + 10
		}
	}
	// 进入该状态满48小时，接下来的21点～24点间随机择时
	if triggerBySecondOfDay(updatedAt, 172800, 75600, 86400) {
		if (fromTarotStatus >= 500 && fromTarotStatus <= 503) || (fromTarotStatus >= 510 && fromTarotStatus <= 513) {
			return fromTarotStatus + 10
		}
	}
	// 进入该状态满12小时后
	if triggerByWaitTime(updatedAt, 43200, 43200) {
		if fromTarotStatus == 603 || fromTarotStatus == 604 {
			return 500
		}
	}
	// 用户3小时没有回复
	if triggerByWaitTime(updatedAt, 10800, 10800) {
		if fromTarotStatus == 104 {
			return 501
		}
		if fromTarotStatus == 205 {
			return 502
		}
		if fromTarotStatus == 305 {
			return 503
		}
	}
	// 用户20～30分钟没有回复
	if triggerByWaitTime(updatedAt, 1200, 1800) {
		if fromTarotStatus == 401 {
			return 402
		}
	}
	// 用户8～12分钟没有回复
	if triggerByWaitTime(updatedAt, 480, 720) {
		if fromTarotStatus == 102 || fromTarotStatus == 103 || (fromTarotStatus >= 201 && fromTarotStatus <= 204) ||
			(fromTarotStatus >= 301 && fromTarotStatus <= 304) || fromTarotStatus == 402 {
			return fromTarotStatus + 1
		}
		if (fromTarotStatus >= 206 && fromTarotStatus <= 209) || fromTarotStatus == 311 || fromTarotStatus == 312 {
			return 203
		}
		if fromTarotStatus == 210 || fromTarotStatus == 211 {
			return 502
		}
		if fromTarotStatus >= 306 && fromTarotStatus <= 310 {
			return 303
		}
		if fromTarotStatus == 313 {
			return 304
		}
		if fromTarotStatus == 403 || fromTarotStatus == 404 {
			return 504
		}
		if fromTarotStatus == 200 {
			return 202
		}
	}
	// 用户6分钟没有回复
	if triggerByWaitTime(updatedAt, 360, 360) {
		if fromTarotStatus == 101 {
			return fromTarotStatus + 1
		}
	}
	return fromTarotStatus
}

func doNothing(fromTarotStatus int) (toTarotStatus int) {
	if fromTarotStatus == 212 {
		return 301
	}
	return fromTarotStatus
}

func receiveHongbao(fromTarotStatus int, msgType int) (toTarotStatus int) {
	if msgType == wxweb.MSG_SYS { // 用户给红包
		if (fromTarotStatus >= 401 && fromTarotStatus <= 404) || fromTarotStatus == 504 {
			return 505
		}
	} else { //红包以外的任何回复
		if fromTarotStatus >= 401 && fromTarotStatus <= 403 {
			return 404
		}
		if fromTarotStatus == 404 {
			return 504
		}
	}
	return fromTarotStatus
}

func receiveCheckImg(fromTarotStatus int, msgType int) (toTarotStatus int) {
	if msgType == wxweb.MSG_IMG { //用户回复图片
		if (fromTarotStatus >= 200 && fromTarotStatus <= 211) || fromTarotStatus == 603 {
			return 301
		}
	} else { //用户回复图片以外的其他
		if fromTarotStatus >= 200 && fromTarotStatus <= 205 {
			return 206
		}
		if fromTarotStatus >= 206 && fromTarotStatus <= 210 {
			return fromTarotStatus + 1
		}
		if fromTarotStatus == 211 { //用户回复图片以外的其他（随机）
			random := rand.Intn(2)
			if random == 0 {
				return 212
			} else {
				return 603
			}
		}
	}
	return fromTarotStatus
}

func receiveBeginTarot(fromTarotStatus int, content string) (toTarotStatus int) {
	if strings.Contains(content, `占卜`) { //用户回复占卜
		if fromTarotStatus == 504 {
			return 594
		}
		if fromTarotStatus == 505 {
			return 595
		}
	} else { //用户回复占卜以外的其他
		if fromTarotStatus == 504 {
			return 584
		}
		if fromTarotStatus == 505 {
			return 585
		}
	}
	return fromTarotStatus
}

func receiveAny(fromTarotStatus int) (toTarotStatus int) {
	if fromTarotStatus >= 101 && fromTarotStatus <= 104 {
		return 201
	}
	if (fromTarotStatus >= 500 && fromTarotStatus <= 503) || (fromTarotStatus >= 510 && fromTarotStatus <= 515) ||
		(fromTarotStatus >= 520 && fromTarotStatus <= 525) || (fromTarotStatus >= 530 && fromTarotStatus <= 533) ||
		(fromTarotStatus >= 600 && fromTarotStatus <= 602) {
		return 200
	}
	return fromTarotStatus
}

func checkNum(fromTarotStatus int, content string, msgType int) (toTarotStatus int) {
	others := false
	if msgType != wxweb.MSG_TEXT {
		others = true
	} else {
		if strings.Contains(content, `[`) && strings.Contains(content, `]`) {
			others = true
		} else {
			content = strings.Replace(content, `一个`, ``, -1)
			content = strings.Replace(content, `一张`, ``, -1)
			content = strings.Replace(content, `一次`, ``, -1)
			content = strings.Replace(content, `一下`, ``, -1)
			num, err := parseNum(content)
			if err == nil && num > 0 {
				if num <= 22 { //用回复包含1～22的数字（及汉字）
					if (fromTarotStatus >= 301 && fromTarotStatus <= 313) || fromTarotStatus == 604 {
						return 401
					}
				} else { //用回复1～22以外的数字
					if (fromTarotStatus >= 301 && fromTarotStatus <= 313) || fromTarotStatus == 604 {
						return 313
					}
				}
			} else {
				others = true
			}
		}
	}
	if others {
		//用户回复数字以外的信息
		if fromTarotStatus >= 306 && fromTarotStatus <= 311 {
			return fromTarotStatus + 1
		}
		if fromTarotStatus >= 301 && fromTarotStatus <= 305 {
			return 306
		}
		if fromTarotStatus == 312 || fromTarotStatus == 313 {
			return 604
		}
	}
	return fromTarotStatus
}

func parseNum(content string) (num int64, err error) {
	regNum := regexp.MustCompile(`\d+`)
	numStr := regNum.FindString(content)
	if numStr != `` {
		return strconv.ParseInt(numStr, 10, 64)
	}
	return int64(checkTarotNum(content)), nil
}

var tarotNumStr = []string{`二十二`, `二十一`, `二十`, `十九`, `十八`, `十七`, `十六`, `十五`, `十四`, `十三`,
	`十二`, `十一`, `十`, `九`, `八`, `七`, `六`, `五`, `四`, `三`, `二`, `一`}

func checkTarotNum(content string) (tarotNum int) {
	num := 0
	for key, value := range tarotNumStr {
		if strings.Contains(content, value) {
			content = strings.Replace(content, value, strconv.Itoa(22-key), 1)
			num = 22 - key
			break
		}
	}
	for _, value := range tarotNumStr {
		if strings.Contains(content, value) {
			// 当有多个字符串时，认为用户输入了大于二十二的汉字数
			return 100
		}
	}
	return num
}
