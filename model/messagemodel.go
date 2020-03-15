package model

import (
	"consumer/common/rabbitmq"
	"github.com/yakaa/log4g"
)

type (
	MessagesModel struct {
		publisher *rabbitmq.Publisher
	}
)

func NewMessagesModel(publisher *rabbitmq.Publisher) *MessagesModel {

	return &MessagesModel{publisher: publisher}
}

func (m *MessagesModel) PublishDelayMessage(message *rabbitmq.Message) error {
	log4g.InfoFormat("start publish delay message %+v", message)
	if err := m.publisher.Push(message); err != nil {
		log4g.ErrorFormat("send publish delay message Error %+v", err)
		return err
	}
	return nil
}
