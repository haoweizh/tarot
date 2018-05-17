package model

import (
	"time"
	"tarot/wechat-go/wxweb"
)

type TarotSentence struct {
	ID               uint `gorm:"primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Content          string
	SentenceType     string
	SentenceScenario int
}

type TarotEvent struct {
	NickName      string
	TarotStatus   int
	UpdatedAt     time.Time
	WelcomeNoResp int
	City          string
	Sex           int
}

type MyContact struct {
	ID            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TarotStatus   int
	TarotNickName string
	WelcomeNoResp int
	NewStatus     bool
	///////////////////////////////////////////
	// from wxweb.User
	Uin               int
	NickName          string
	HeadImgUrl        string
	ContactFlag       int // 0 自己
	MemberCount       int
	RemarkName        string
	PYInitial         string
	PYQuanPin         string
	RemarkPYInitial   string
	RemarkPYQuanPin   string
	HideInputBarFlag  int
	StarFriend        int    // 是否为星标朋友  0-否  1-是
	Sex               int    // 0-未设置（公众号、保密），1-男，2-女
	Signature         string // 公众号的功能介绍 or 好友的个性签名
	AppAccountFlag    int
	Statues           int
	AttrStatus        uint32
	Province          string
	City              string
	Alias             string
	VerifyFlag        int
	OwnerUin          int
	WebWxPluginSwitch int
	HeadImgFlag       int
	SnsFlag           int
	UniFriend         int
	DisplayName       string
	ChatRoomId        int
	KeyWord           string
	EncryChatRoomId   string
	IsOwner           int
	MemberStatus      int
}

func (myContact *MyContact) Init(value *wxweb.User) {
	myContact.Uin = value.Uin
	//contact.NickName = value.NickName nick name is used as unique key
	myContact.HeadImgUrl = value.HeadImgUrl
	myContact.ContactFlag = value.ContactFlag
	myContact.MemberCount = value.MemberCount
	myContact.RemarkName = value.RemarkName
	myContact.PYInitial = value.PYInitial
	myContact.PYQuanPin = value.PYQuanPin
	myContact.RemarkPYInitial = value.RemarkPYInitial
	myContact.RemarkPYQuanPin = value.RemarkPYQuanPin
	myContact.HideInputBarFlag = value.HideInputBarFlag
	myContact.StarFriend = value.StarFriend
	myContact.Sex = value.Sex
	myContact.Signature = value.Signature
	myContact.AppAccountFlag = value.AppAccountFlag
	myContact.Statues = value.Statues
	myContact.AttrStatus = value.AttrStatus
	myContact.Province = value.Province
	myContact.City = value.City
	myContact.Alias = value.Alias
	myContact.VerifyFlag = value.VerifyFlag
	myContact.OwnerUin = value.OwnerUin
	myContact.WebWxPluginSwitch = value.WebWxPluginSwitch
	myContact.HeadImgFlag = value.HeadImgFlag
	myContact.SnsFlag = value.SnsFlag
	myContact.UniFriend = value.UniFriend
	myContact.DisplayName = value.DisplayName
	myContact.ChatRoomId = value.ChatRoomId
	myContact.KeyWord = value.KeyWord
	myContact.EncryChatRoomId = value.EncryChatRoomId
	myContact.IsOwner = value.IsOwner
	myContact.MemberStatus = value.MemberStatus
}
