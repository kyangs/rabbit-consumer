package service

import (
	"consumer/common/httpx"
	"consumer/common/rabbitmq"
	"consumer/common/utils"
	"consumer/model"
	"github.com/yakaa/log4g"
)

type (
	MessageService struct {
		messageModel *model.MessagesModel
	}
)

func NewMessageService(messageModel *model.MessagesModel) *MessageService {

	return &MessageService{messageModel: messageModel}
}

func (s *MessageService) ConsumerMessage(message *rabbitmq.Message) error {
	if message.Delay > 0 && message.IsDelay == false {
		message.IsDelay = true
		return s.messageModel.PublishDelayMessage(message)
	}
	if responseStatus, err := utils.HttpRequest(httpx.HttpMethodPost, message.Url, message.Data); err != nil || responseStatus == false {
		log4g.ErrorFormat("http send message Error %+v", err)
		if message.RetryTime > 0 {
			message.Delay = message.RetryTime
			return s.messageModel.PublishDelayMessage(message)
		}
	}
	return nil
}
