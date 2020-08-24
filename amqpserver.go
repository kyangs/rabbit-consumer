package main

import (
	"log"

	"consumer/common/rabbitmq"
	"consumer/config"
	"consumer/model"
	"consumer/service"

	"github.com/yakaa/log4g"
)

func main() {

	conf, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("ParseConfig %+v", err)
	}

	log4g.Init(conf.Log4g)
	publisher, err := rabbitmq.BuildPublisher(conf.RabbitMq)
	if err != nil {
		log.Fatalf("create publisher err %+v", err)
	}
	if err := publisher.Push(&rabbitmq.Message{Type: rabbitmq.TypePong}); err != nil {
		log.Fatalf("send pong err %+v", err)
	}
	dingTalkAlarmService := &service.WebHookAlarmService{Conf: conf.Hook}

	consumerPool, err := rabbitmq.BuildConsumerPool(
		conf.RabbitMq,
		service.NewMessageService(model.NewMessagesModel(publisher), dingTalkAlarmService).ConsumerMessage,
		conf.RabbitMq.ConsumerAmount,
	)
	publisher.Close()
	if err != nil {
		log.Fatalf("create publisher fail %+v", err)
	}
	log.Fatal(consumerPool.Run())
}
