package feishux

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/ruomm/goxframework/gox/httpx"
	"time"
)

//	type FsExample struct {
//		MsgType string `json:"msg_type"`
//		Content struct {
//			Post struct {
//				ZhCn struct {
//					Title   string `json:"title"`
//					Content [][]struct {
//						Tag    string `json:"tag"`
//						Text   string `json:"text,omitempty"`
//						Href   string `json:"href,omitempty"`
//						UserId string `json:"user_id,omitempty"`
//					} `json:"content"`
//				} `json:"zh_cn"`
//			} `json:"post"`
//		} `json:"content"`
//	}
type FeishuRobotConfigs struct {
	//# 聊天机器人的地址
	WebHookURL string `yaml:"webHookURL"`
	//# 聊天机器人的签名
	WebHookKey string `yaml:"webHookKey"`
	//# 飞书机器人消息标题
	RobotMsgTitle string `yaml:"robotMsgTitle"`
	// # 飞书机器人消息是否富文本模式
	RobotMsgRichEnable bool `yaml:"robotMsgRichEnable"`
}

type FsMessageResult struct {
	//StatusCode    int    `json:"StatusCode"`
	//StatusMessage string `json:"StatusMessage"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type FsMessageRichText struct {
	Timestamp string        `json:"timestamp"`
	Sign      string        `json:"sign"`
	MsgType   string        `json:"msg_type"`
	Content   FsContentRich `json:"content"`
}
type FsMessageText struct {
	Timestamp string        `json:"timestamp"`
	Sign      string        `json:"sign"`
	MsgType   string        `json:"msg_type"`
	Content   FsContextText `json:"content"`
}
type FsContextText struct {
	Text string `json:"text"`
}
type FsContentRich struct {
	Post FsPost `json:"post,omitempty"`
}
type FsPost struct {
	ZhCn *FsZhcnEnus `json:"zh_cn,omitempty"`
	EnUs *FsZhcnEnus `json:"en_us,omitempty"`
}

type FsZhcnEnus struct {
	Title   string           `json:"title,omitempty"`
	Content [][]FsSubContent `json:"content,omitempty"`
}

type FsSubContent struct {
	Tag      string `json:"tag,omitempty"`
	Text     string `json:"text,omitempty"`
	Href     string `json:"href,omitempty"`
	UserId   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

func (t *FsMessageRichText) GenSign(secret string) error {
	if len(secret) <= 0 {
		return nil
	}
	timeStamp := fmt.Sprintf("%v", time.Now().Unix())
	stringToSign := fmt.Sprintf(timeStamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	t.Timestamp = timeStamp
	t.Sign = signature
	return nil
}

func (t *FsMessageText) GenSign(secret string) error {
	if len(secret) <= 0 {
		return nil
	}
	timeStamp := fmt.Sprintf("%v", time.Now().Unix())
	stringToSign := fmt.Sprintf(timeStamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	t.Timestamp = timeStamp
	t.Sign = signature
	return nil
}

func (t *FsMessageText) SendMessageByRobot(webHookUrl string, webHookKey string) (*FsMessageResult, error) {
	err := t.GenSign(webHookKey)
	if err != nil {
		fmt.Println("Send Text Message By Feishu Rebot GenSign err:" + err.Error())
		return nil, err
	}
	feishuMessageResult := FsMessageResult{Code: -1, Msg: "发送失败"}
	_, err = httpx.DoPostJson(nil, webHookUrl, t, &feishuMessageResult)
	if err != nil {
		fmt.Println("Send Text Message By Feishu Rebot Request err:" + err.Error())
		return &feishuMessageResult, err
	} else {
		return &feishuMessageResult, err
	}
}

func (t *FsMessageRichText) SendMessageByRobot(webHookUrl string, webHookKey string) (*FsMessageResult, error) {
	err := t.GenSign(webHookKey)
	if err != nil {
		fmt.Println("Send Rich Text Message By Feishu Rebot GenSign err:" + err.Error())
		return nil, err
	}
	feishuMessageResult := FsMessageResult{Code: -1, Msg: "发送失败"}
	_, err = httpx.DoPostJson(nil, webHookUrl, t, &feishuMessageResult)
	if err != nil {
		fmt.Println("Send Rich Text Message By Feishu Rebot Request err:" + err.Error())
		return &feishuMessageResult, err
	} else {
		return &feishuMessageResult, nil
	}
}
