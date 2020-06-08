package service

import (
	"consumer/common/httpx"
	"consumer/common/utils"
	"consumer/config"
	"github.com/yakaa/log4g"
)

type (
	DingTalkAlarmService struct {
		Conf config.DingTalk
	}
	TextMessage struct {
		Content string `json:"content"`
	}
	DingTalkMessage struct {
		MsgType string       `json:"msgtype"`
		Text    *TextMessage `json:"text"`
	}
)

// 发送DingTalk 报警消息
func (d *DingTalkAlarmService) send(content string) {

	if len(d.Conf.WebHook) == 0 {
		return
	}

	for _, v := range d.Conf.WebHook {

		b, err := utils.HttpRequest(httpx.HttpMethodPost, v, &DingTalkMessage{
			MsgType: "text",
			Text: &TextMessage{
				Content: content,
			},
		})
		if b == false || err != nil {
			log4g.ErrorFormat("发送消费报警到钉钉失败", err)
		}
	}
	return
}
