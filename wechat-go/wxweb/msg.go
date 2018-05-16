package wxweb

import "fmt"

// Msg implement this interface, can added addition send by wechat
type Msg interface {
	Path() string
	To() string
	Content() map[string]interface{}

	Description() string
}

// FileMsg struct
type FileMsg struct {
	to      string
	mediaID string
	path    string
	ftype   int
	fname   string
	ext     string
}

// Path is text msg's api path
func (msg *FileMsg) Path() string {
	return msg.path
}

// To destination
func (msg *FileMsg) To() string {
	return msg.to
}

// Content text msg's content
func (msg *FileMsg) Content() map[string]interface{} {
	content := make(map[string]interface{}, 0)

	content[`Type`] = msg.ftype

	if msg.ftype == 6 {
		content[`Content`] = fmt.Sprintf(`<appmsg appid='wxeb7ec651dd0aefa9' sdkver=''><title>%s</title><des></des><action></action><type>6</type><content></content><url></url><lowurl></lowurl><appattach><totallen>10</totallen><attachid>%s</attachid><fileext>%s</fileext></appattach><extinfo></extinfo></appmsg>`, msg.fname, msg.mediaID, msg.ext)
	} else {
		content[`MediaId`] = msg.mediaID
	}

	return content
}

func (msg *FileMsg) Description() string {
	return fmt.Sprintf(`[FileMsg] %s`, msg.fname)
}

// NewFileMsg construct a new FileMsg's instance
func NewFileMsg(mediaID, to, name, ext string) *FileMsg {
	return &FileMsg{to, mediaID, `webwxsendappmsg?fun=async&f=json`, 6, name, ext}
}

// NewImageMsg ..
func NewImageMsg(mediaID, to string) *FileMsg {
	return &FileMsg{to, mediaID, `webwxsendmsgimg?fun=async&f=json`, 3, ``, ``}
}

// NewVideoMsg ..
func NewVideoMsg(mediaID, to string) *FileMsg {
	return &FileMsg{to, mediaID, `webwxsendvideomsg?fun=async&f=json`, 43, ``, ``}
}

func (msg *FileMsg) String() string {
	if msg.ftype == 3 {
		return `IMAGE`
	} else if msg.ftype == 4 {
		return `GIF EMOTICON`
	} else {
		return `FILE`
	}
}

// TextMsg struct
type TextMsg struct {
	to      string
	content string
}

// Path is text msg's api path
func (msg *TextMsg) Path() string {
	return `webwxsendmsg`
}

// To destination
func (msg *TextMsg) To() string {
	return msg.to
}

// Content text msg's content
func (msg *TextMsg) Content() map[string]interface{} {
	content := make(map[string]interface{}, 0)

	content["Type"] = 1
	content["Content"] = msg.content

	return content
}

func (msg *TextMsg) Description() string {
	return fmt.Sprintf(`[TextMessage] %s`, msg.content)
}

// NewTextMsg construct a new TextMsg's instance
func NewTextMsg(text, to string) *TextMsg {
	return &TextMsg{to, text}
}

func (msg *TextMsg) String() string {
	return msg.content
}

// EmoticonMsg is wechat emoticon msg
type EmoticonMsg struct {
	to      string
	mediaID string
}

// Path is text msg's api path
func (msg *EmoticonMsg) Path() string {
	return `webwxsendemoticon?fun=sys`
}

// To destination
func (msg *EmoticonMsg) To() string {
	return msg.to
}

// Content text msg's content
func (msg *EmoticonMsg) Content() map[string]interface{} {
	content := make(map[string]interface{}, 0)

	content[`Type`] = 47
	content[`MediaId`] = msg.mediaID
	content[`EmojiFlag`] = 2

	return content
}

func (msg *EmoticonMsg) Description() string {
	return fmt.Sprintf(`[TextMessage] %s`, msg.mediaID)
}

// NewEmoticonMsgMsg create a new instance
func NewEmoticonMsgMsg(mid, to string) *EmoticonMsg {
	return &EmoticonMsg{to, mid}
}

func (msg *EmoticonMsg) String() string {
	return `GIF EMOTICON`
}
