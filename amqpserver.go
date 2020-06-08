package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"consumer/common/rabbitmq"
	"consumer/config"
	"consumer/model"
	"consumer/service"

	"github.com/yakaa/log4g"
)

var configFile = flag.String("c", "config.json", "Please set config file")

func main() {
	flag.Parse()
	body, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("read file %s: %s", *configFile, err)
	}
	conf := new(config.Config)
	if err := json.Unmarshal(body, conf); err != nil {
		log.Fatalf("json.Unmarshal %s: %s", *configFile, err)
	}
	log4g.Init(conf.Log4g)
	publisher, err := rabbitmq.BuildPublisher(conf.RabbitMq)
	if err != nil {
		log.Fatalf("create publisher err %+v", err)
	}
	if err := publisher.Push(&rabbitmq.Message{Type: rabbitmq.TypePong}); err != nil {
		log.Fatalf("send pong err %+v", err)
	}
	dingTalkAlarmService := &service.DingTalkAlarmService{Conf: conf.DingTalk}

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
