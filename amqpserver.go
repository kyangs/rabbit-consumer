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

var configFile = flag.String("c", "config/config.json", "Please set config file")

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
	consumer, err := rabbitmq.BuildConsumer(
		conf.RabbitMq,
		service.NewMessageService(model.NewMessagesModel(publisher)).ConsumerMessage,
	)
	publisher.Close()
	if err != nil {
		log.Fatalf("create publisher fail %+v", err)
	}
	log4g.Error(consumer.Run())
}
