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
		messageModel  *model.MessagesModel
		alarmServices []BaseAlarmService
	}
)

func NewMessageService(messageModel *model.MessagesModel, alarmServices ...BaseAlarmService) *MessageService {

	return &MessageService{
		messageModel:  messageModel,
		alarmServices: alarmServices,
	}
}

func (s *MessageService) ConsumerMessage(message *rabbitmq.Message) error {
	if message.Type == rabbitmq.TypePong {
		log4g.Info("receive pong message")
		return nil
	}
	if message.Delay > 0 && message.IsDelay == false {
		message.IsDelay = true
		err := s.messageModel.PublishDelayMessage(message)
		if err != nil {
			s.sendAlarm("【消息消费】发送延时消息对列失败，请检查对列")
		}
		return err
	}
	if responseStatus, err := utils.HttpRequest(httpx.HttpMethodPost, message.Url, message.Data); err != nil || responseStatus == false {
		log4g.ErrorFormat("http send message Error %+v", err)
		s.sendAlarm("【消息消费】发送请求到URL：" + message.Url + "失败，原因：" + err.Error())
		if message.RetryTime > 0 {
			message.Delay = message.RetryTime
			err := s.messageModel.PublishDelayMessage(message)
			if err != nil {
				s.sendAlarm("【消息消费】发送请求失败后启动延时重试，但发送消息到重试队列失败")
			}
			return err
		}
	}
	return nil
}

// 批量发送报警消息
func (s *MessageService) sendAlarm(content string) {
	if len(s.alarmServices) == 0 {
		return
	}
	for _, alarmService := range s.alarmServices {
		alarmService.send(content)
	}
}
