/*
Copyright 2017 wechat-go Authors. All Rights Reserved.
MIT License

Copyright (c) 2017

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package wxweb

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/songtianyi/rrframework/config"
	"crypto/md5"
	"encoding/hex"
	"gopkg.in/h2non/filetype.v1"
)

type ApiV2 struct {
	httpClient *Client
}

// JsLogin: jslogin api
func (api *ApiV2) JsLogin(common *Common) (string, error) {

	km := url.Values{}
	km.Add("appid", common.AppId)
	km.Add("fun", "new")
	km.Add("lang", common.Lang)
	km.Add("redirect_uri", common.RedirectUri)
	km.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	uri := common.LoginUrl + "/jslogin?" + km.Encode()

	body, _ := api.httpClient.Get(uri, nil)

	ss := strings.Split(string(body), "\"")
	if len(ss) < 2 {
		return "", fmt.Errorf("jslogin response invalid, %s", string(body))
	}
	return ss[1], nil
}

// QrCode: get qrcode
func (api *ApiV2) QrCode(common *Common, uuid string) ([]byte, error) {
	km := url.Values{}
	km.Add("t", "webwx")
	km.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	uri := common.LoginUrl + "/qrcode/" + uuid + "?" + km.Encode()
	body, _ := api.httpClient.Get(uri, nil)
	return body, nil
}

// Login: login api
func (api *ApiV2) Login(common *Common, uuid, tip string) (string, error) {
	km := url.Values{}
	km.Add("tip", tip)
	km.Add("uuid", uuid)
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	uri := common.LoginUrl + "/cgi-bin/mmwebwx-bin/login?" + km.Encode()
	body, _ := api.httpClient.Get(uri, nil)
	strb := string(body)
	if strings.Contains(strb, "window.code=200") &&
		strings.Contains(strb, "window.redirect_uri") {
		ss := strings.Split(strb, "\"")
		if len(ss) < 2 {
			return "", fmt.Errorf("parse redirect_uri fail, %s", strb)
		}
		return ss[1], nil
	}

	return "", fmt.Errorf("login response, %s", strb)
}

// WebNewLoginPage: webwxnewloginpage api
func (api *ApiV2) WebNewLoginPage(common *Common, xc *XmlConfig, uri string) ([]*http.Cookie, error) {
	u, _ := url.Parse(uri)
	km := u.Query()
	km.Add("fun", "new")
	uri = common.CgiUrl + "/webwxnewloginpage?" + km.Encode()
	resp, _ := api.httpClient.fetchReponse("GET", uri, []byte(""), Header{})
	body, _ := ioutil.ReadAll(resp.Body)
	if err := xml.Unmarshal(body, xc); err != nil {
		return nil, err
	}
	if xc.Ret != 0 {
		return nil, fmt.Errorf("xc.Ret != 0, %s", string(body))
	}
	return resp.Cookies(), nil
}

// WebWxInit: webwxinit api
func (api *ApiV2) WebWxInit(common *Common, ce *XmlConfig) ([]byte, error) {
	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("skey", ce.Skey)
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))

	uri := common.CgiUrl + "/webwxinit?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
	}

	b, _ := json.Marshal(js)

	body, _ := api.httpClient.PostJsonByte(uri, b)

	return body, nil
}

// SyncCheck: synccheck api
func (api *ApiV2) SyncCheck(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	server string, skl *SyncKeyList) (retcode int, selector int, err error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix()*1000, 10))
	km.Add("sid", ce.Wxsid)
	km.Add("uin", ce.Wxuin)
	km.Add("skey", ce.Skey)
	km.Add("deviceid", common.DeviceID)
	km.Add("synckey", skl.String())
	km.Add("_", strconv.FormatInt(time.Now().Unix()*1000, 10))
	uri := "https://" + server + "/cgi-bin/mmwebwx-bin/synccheck?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
	}

	b, _ := json.Marshal(js)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	api.httpClient.SetJar(jar)

	body, _ := api.httpClient.GetByte(uri, b)

	strb := string(body)
	reg := regexp.MustCompile("window.synccheck={retcode:\"(\\d+)\",selector:\"(\\d+)\"}")
	sub := reg.FindStringSubmatch(strb)
	retcode = 0
	selector = 0
	if len(sub) >= 2 {
		retcode, _ = strconv.Atoi(sub[1])
		selector, _ = strconv.Atoi(sub[2])
	}

	return retcode, selector, nil
}

// WebWxSync: webwxsync api
func (api *ApiV2) WebWxSync(common *Common, ce *XmlConfig, cookies []*http.Cookie, msg chan []byte, skl *SyncKeyList) ([]*http.Cookie, error) {

	km := url.Values{}
	km.Add("skey", ce.Skey)
	km.Add("sid", ce.Wxsid)
	km.Add("lang", common.Lang)
	km.Add("pass_ticket", ce.PassTicket)

	uri := common.CgiUrl + "/webwxsync?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		SyncKey: skl,
		rr:      ^int(time.Now().Unix()) + 1,
	}

	b, _ := json.Marshal(js)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	api.httpClient.SetJar(jar)

	resp, body, _ := api.httpClient.PostJsonByteForResp(uri, b)

	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		return nil, err
	}
	retcode, err := jc.GetInt("BaseResponse.Ret")
	if err != nil {
		return nil, err
	}
	if retcode != 0 {
		return nil, fmt.Errorf("BaseResponse.Ret %d", retcode)
	}

	msg <- body

	skl.List = skl.List[:0]
	skl1, _ := GetSyncKeyListFromJc(jc)
	skl.Count = skl1.Count
	skl.List = append(skl.List, skl1.List...)
	return resp.Cookies(), nil
}
func (api *ApiV2) WebWxSyncFlushCookie(common *Common, ce *XmlConfig, cookies []*http.Cookie, skl *SyncKeyList) ([]*http.Cookie, error) {

	km := url.Values{}
	km.Add("skey", ce.Skey)
	km.Add("sid", ce.Wxsid)
	km.Add("lang", common.Lang)
	km.Add("pass_ticket", ce.PassTicket)

	uri := common.CgiUrl + "/webwxsync?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		SyncKey: skl,
		rr:      ^int(time.Now().Unix()) + 1,
	}

	b, _ := json.Marshal(js)
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	api.httpClient.SetJar(jar)
	resp, body, _ := api.httpClient.PostJsonByteForResp(uri, b)

	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		return nil, err
	}
	retcode, err := jc.GetInt("BaseResponse.Ret")
	if err != nil {
		return nil, err
	}
	if retcode != 0 {
		return nil, fmt.Errorf("BaseResponse.Ret %d", retcode)
	}

	skl.List = skl.List[:0]
	skl1, _ := GetSyncKeyListFromJc(jc)
	skl.Count = skl1.Count
	skl.List = append(skl.List, skl1.List...)
	return resp.Cookies(), nil
}

// WebWxStatusNotify: webwxstatusnotify api
func (api *ApiV2) WebWxStatusNotify(common *Common, ce *XmlConfig, bot *User) (int, error) {
	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("lang", common.Lang)
	uri := common.CgiUrl + "/webwxstatusnotify?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		Code:         3,
		FromUserName: bot.UserName,
		ToUserName:   bot.UserName,
		ClientMsgId:  int(time.Now().Unix()),
	}

	b, _ := json.Marshal(js)

	body, _ := api.httpClient.PostJsonByte(uri, b)
	jc, _ := rrconfig.LoadJsonConfigFromBytes(body)
	ret, _ := jc.GetInt("BaseResponse.Ret")
	return ret, nil
}

// WebWxGetContact: webwxgetcontact api
func (api *ApiV2) WebWxGetContact(common *Common, ce *XmlConfig, cookies []*http.Cookie) ([]byte, error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("seq", "0")
	km.Add("skey", ce.Skey)
	uri := common.CgiUrl + "/webwxgetcontact?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
	}

	b, _ := json.Marshal(js)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	api.httpClient.SetJar(jar)
	body, _ := api.httpClient.PostJsonByte(uri, b)
	return body, nil
}

// WebWxSendMsgImg: webwxsendmsgimg api
func (api *ApiV2) WebWxSendFile(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	from, to, media string, message Msg) {
	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("fun", "async")
	km.Add("f", "json")
	km.Add("lang", common.Lang)

	//uri := common.CgiUrl + "/webwxsendmsgimg?" + km.Encode()
	uri := fmt.Sprintf(`%s/%s?%s`, common.CgiUrl, message.Path(), km.Encode())
	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		Msg: &MediaMessage{
			Type:         3,
			Content:      "",
			FromUserName: from,
			ToUserName:   to,
			LocalID:      int(time.Now().Unix() * 1e4),
			ClientMsgId:  int(time.Now().Unix() * 1e4),
			MediaId:      media,
		},
		Scene: 0,
	}

	b, _ := json.Marshal(js)
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)

	api.httpClient.SetJar(jar)
	api.httpClient.PostJsonByte(uri, b)
}

// WebWxSendMsgImg: webwxsendmsgimg api
func (api *ApiV2) WebWxSendMsgImg(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	from, to, media string) (int, error) {

	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("fun", "async")
	km.Add("f", "json")
	km.Add("lang", common.Lang)

	uri := common.CgiUrl + "/webwxsendmsgimg?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		Msg: &MediaMessage{
			Type:         3,
			Content:      "",
			FromUserName: from,
			ToUserName:   to,
			LocalID:      int(time.Now().Unix() * 1e4),
			ClientMsgId:  int(time.Now().Unix() * 1e4),
			MediaId:      media,
		},
		Scene: 0,
	}

	b, _ := json.Marshal(js)
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)

	api.httpClient.SetJar(jar)
	body, _ := api.httpClient.PostJsonByte(uri, b)
	jc, _ := rrconfig.LoadJsonConfigFromBytes(body)
	ret, _ := jc.GetInt("BaseResponse.Ret")
	return ret, nil
}

// WebWxSendMsg: webwxsendmsg api
func (api *ApiV2) WebWxSendMsg(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	from, to string, msg string) ([]byte, error) {

	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)

	uri := common.CgiUrl + "/webwxsendmsg?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		Msg: &TextMessage{
			Type:         1,
			Content:      msg,
			FromUserName: from,
			ToUserName:   to,
			LocalID:      int(time.Now().Unix() * 1e4),
			ClientMsgId:  int(time.Now().Unix() * 1e4),
		},
	}

	b, _ := json.Marshal(js)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	api.httpClient.SetJar(jar)
	body, _ := api.httpClient.PostJsonByte(uri, b)
	return body, nil
}

// WebWxUploadMedia: webwxuploadmedia api

// UploadMedia is a convernice method to upload attachment to wx cdn.
func (api *ApiV2) WebWxUploadMedia(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	filename string, content []byte, from string, to string) (string, error) {

	// Only the first 261 bytes are used to sniff the content type.
	head := content[:261]

	var mediatype string
	if filetype.IsImage(head) {
		mediatype = `pic`
	} else if filetype.IsVideo(head) {
		mediatype = `video`
	} else {
		mediatype = `doc`
	}
	var dataTicket string
	for _, v := range cookies {
		if strings.Contains(v.String(), "webwx_data_ticket") {
			dataTicket = strings.Split(v.String(), "=")[1]
			break
		}
	}
	kind, _ := filetype.Get(content)
	fields := map[string]string{
		`id`:                `WU_FILE_` + strconv.Itoa(int(common.MediaCount)),
		`name`:              filename,
		`type`:              kind.MIME.Value,
		`lastModifiedDate`:  "Mon Feb 13 2017 17:27:23 GMT+0800 (CST)",
		`size`:              strconv.Itoa(len(content)),
		`mediatype`:         mediatype,
		`pass_ticket`:       ce.PassTicket,
		`webwx_data_ticket`: dataTicket,
	}

	h := md5.New()
	h.Write(content)
	mediaMap := map[string]interface{}{
		`BaseRequest`: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		`ClientMediaId`: strconv.FormatInt(time.Now().Unix(), 10),
		`TotalLen`:      len(content),
		`StartPos`:      0,
		`DataLen`:       len(content),
		`MediaType`:     4,
		`UploadType`:    2,
		`ToUserName`:    to,
		`FromUserName`:  from,
		`FileMd5`:       hex.EncodeToString(h.Sum(nil)),
	}
	chunks := int((len(content)-1)/524288) + 1
	var err error
	var resp []byte
	if chunks == 1 {
		resp, err = api.uploadMediaChunk(common, cookies, fields, mediaMap, filename, content)
	} else {
		fields[`chunks`] = strconv.Itoa(chunks)
		for chunk := 0; chunk < chunks; chunk++ {
			fields[`chunk`] = strconv.Itoa(chunk)
			end := (chunk + 1) * 524288
			if end > len(content) {
				end = len(content)
			}
			subBuf := content[chunk*524288 : end]
			resp, err = api.uploadMediaChunk(common, cookies, fields, mediaMap, filename, subBuf)
			if err != nil {
				break
			}
		}
	}
	jc, err := rrconfig.LoadJsonConfigFromBytes(resp)
	if err != nil {
		return "", err
	}
	ret, _ := jc.GetInt("BaseResponse.Ret")
	if ret != 0 {
		return "", fmt.Errorf("BaseResponse.Ret=%d", ret)
	}
	return jc.GetString("MediaId")
}

func (api *ApiV2) uploadMediaChunk(common *Common, cookies []*http.Cookie, fields map[string]string,
	mediaMap map[string]interface{}, fileName string, buf []byte) (resp []byte, error error) {

	media, err := json.Marshal(&mediaMap)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile(`filename`, fileName)
	if err != nil {
		return nil, err
	}
	fw.Write(buf)
	for k, v := range fields {
		writer.WriteField(k, v)
	}

	writer.WriteField(`uploadmediarequest`, string(media))
	writer.Close()

	urlOBJ, err := url.Parse(common.UploadUrl)
	if err != nil {
		return nil, err
	}
	host := urlOBJ.Host
	if err != nil {
		return nil, err
	}
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(common.UploadUrl)
	jar.SetCookies(u, cookies)
	api.httpClient.SetJar(jar)

	return api.httpClient.fetchWithReader("POST",
		fmt.Sprintf(`https://%s/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json`, host), body,
		Header{`Content-Type`: writer.FormDataContentType()})
}

// WebWxGetMsgImg: webwxgetmsgimg api
func (api *ApiV2) WebWxGetMsgImg(common *Common, ce *XmlConfig, cookies []*http.Cookie, msgId string) ([]byte, error) {
	km := url.Values{}
	km.Add("MsgID", msgId)
	km.Add("skey", ce.Skey)
	km.Add("type", "slave")

	uri := common.CgiUrl + "/webwxgetmsgimg?" + km.Encode()

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)

	api.httpClient.SetJar(jar)
	body, _ := api.httpClient.GetWithHeader(uri, Header{"Content-Type": "image/jpeg"})

	return body, nil
}

// WebWxSendEmoticon: webwxsendemoticon api
func (api *ApiV2) WebWxSendEmoticon(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	from, to, media string) (int, error) {

	km := url.Values{}
	km.Add("fun", "sys")
	km.Add("lang", common.Lang)

	uri := common.CgiUrl + "/webwxsendemoticon?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		Msg: &EmotionMessage{
			Type:         47,
			EmojiFlag:    2,
			FromUserName: from,
			ToUserName:   to,
			LocalID:      int(time.Now().Unix() * 1e4),
			ClientMsgId:  int(time.Now().Unix() * 1e4),
			MediaId:      media,
		},
		Scene: 0,
	}

	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return -1, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jc, _ := rrconfig.LoadJsonConfigFromBytes(body)
	ret, _ := jc.GetInt("BaseResponse.Ret")
	return ret, nil
}

// WebWxGetIcon: webwxgeticon api
func (api *ApiV2) WebWxGetIcon(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	username, chatroomid string) ([]byte, error) {
	km := url.Values{}
	km.Add("seq", "0")
	km.Add("username", username)
	if chatroomid != "" {
		km.Add("chatroomid", chatroomid)
	}
	km.Add("skey", ce.Skey)
	uri := common.CgiUrl + "/webwxgeticon?" + km.Encode()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

// WebWxGetIconByHeadImgUrl: get head img
func (api *ApiV2) WebWxGetIconByHeadImgUrl(common *Common, ce *XmlConfig, cookies []*http.Cookie, headImgUrl string) ([]byte, error) {
	uri := common.CgiDomain + headImgUrl

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

// WebWxBatchGetContact: webwxbatchgetcontact api
func (api *ApiV2) WebWxBatchGetContact(common *Common, ce *XmlConfig, cookies []*http.Cookie, cl []*User) ([]byte, error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("type", "ex")
	uri := common.CgiUrl + "/webwxbatchgetcontact?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		Count: len(cl),
		List:  cl,
	}

	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

// WebWxVerifyUser: webwxverifyuser api
func (api *ApiV2) WebWxVerifyUser(common *Common, ce *XmlConfig, cookies []*http.Cookie, opcode int, verifyContent string, vul []*VerifyUser) ([]byte, error) {
	var body []byte
	i := 0
	for i++; i <= 10; {
		km := url.Values{}
		km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
		km.Add("pass_ticket", ce.PassTicket)

		uri := common.CgiUrl + "/webwxverifyuser?" + km.Encode()
		js := InitReqBody{
			BaseRequest: &BaseRequest{
				ce.Wxuin,
				ce.Wxsid,
				ce.Skey,
				common.DeviceID,
			},
			Opcode:             opcode,
			SceneList:          []int{33},
			SceneListCount:     1,
			VerifyContent:      verifyContent,
			VerifyUserList:     vul,
			VerifyUserListSize: len(vul),
			skey:               ce.Skey,
		}
		b, _ := json.Marshal(js)

		jar, _ := cookiejar.New(nil)
		u, _ := url.Parse(uri)
		jar.SetCookies(u, cookies)
		api.httpClient.SetJar(jar)
		var err error
		body, err = api.httpClient.PostJsonByte(uri, b)
		if err != nil {
			if i >= 10 {
				return nil, err
			} else {
				continue
			}
		}

		break
	}
	return body, nil
}

// WebWxCreateChatroom: webwxcreatechatroom api
func (api *ApiV2) WebWxCreateChatroom(common *Common, ce *XmlConfig, cookies []*http.Cookie, users []*User, topic string) (interface{}, error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("pass_ticket", ce.PassTicket)

	uri := common.CgiUrl + "/webwxcreatechatroom?" + km.Encode()
	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		MemberCount: len(users),
		MemberList:  users,
		Topic:       topic,
	}
	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

// WebWxRevokeMsg: webwxrevokemsg api
func (api *ApiV2) WebWxRevokeMsg(common *Common, ce *XmlConfig, cookies []*http.Cookie, clientMsgId, svrMsgId, toUserName string) error {
	km := url.Values{}
	km.Add("lang", common.Lang)

	uri := common.CgiUrl + "/webwxrevokemsg?" + km.Encode()
	js := RevokeReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		ClientMsgId: clientMsgId,
		SvrMsgId:    svrMsgId,
		ToUserName:  toUserName,
	}
	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		return err
	}
	retcode, _ := jc.GetInt("BaseResponse.Ret")
	if retcode != 0 {
		return fmt.Errorf("BaseResponse.Ret %d", retcode)
	}
	return nil
}

// WebWxlogout: webwxlogout api
func (api *ApiV2) WebWxLogout(common *Common, ce *XmlConfig, cookies []*http.Cookie) error {
	km := url.Values{}
	km.Add("redirect", "1")
	km.Add("type", "1")
	km.Add("skey", ce.Skey)

	uri := common.CgiUrl + "/webwxlogout?" + km.Encode()
	js := LogoutReqBody{
		uin: ce.Wxuin,
		sid: ce.Wxsid,
	}
	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
func NewApiV2() *ApiV2 {
	return &ApiV2{httpClient: NewClient()}
}
