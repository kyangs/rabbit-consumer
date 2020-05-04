package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"consumer/common/rabbitmq"
	"consumer/config"
	"github.com/yakaa/log4g"
)

var configFileTest = flag.String("c", "config.json", "Please set config file")

func main() {
	flag.Parse()
	body, err := ioutil.ReadFile(*configFileTest)
	if err != nil {
		log.Fatalf("read file %s: %s", *configFileTest, err)
	}
	conf := new(config.Config)
	if err := json.Unmarshal(body, conf); err != nil {
		log.Fatalf("json.Unmarshal %s: %s", *configFileTest, err)
	}
	log4g.Init(conf.Log4g)
	publisher, err := rabbitmq.BuildPublisher(conf.RabbitMq)
	if err != nil {
		log.Fatalf("create publisher err %+v", err)
	}
	for i := 0; i < 20; i++ {
		if err := publisher.Push(&rabbitmq.Message{
			Type: rabbitmq.TypeMessage,
			Url:  "https://baidu.com",
			Data: []int{i},
		}); err != nil {
			log.Fatalf("send pong err %+v", err)
		}
	}

	publisher.Close()
	log4g.Info("SUCCESS")

}
