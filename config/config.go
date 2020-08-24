package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"consumer/common/utils"
	"github.com/yakaa/log4g"
)

//amqp.Dial("amqp://guest:guest@localhost:5672/")
type (
	Config struct {
		Log4g    log4g.Config
		RabbitMq RabbitMq
		RsaCert  RsaCert
		Hook     Hook
	}

	RabbitMq struct {
		DataSource     string
		QueueName      string
		Consumer       string
		ConsumerAmount int
		Exchange       string
		DeliveryKey    string
		Durable        bool
		AutoDelete     bool
		AutoAck        bool
		Exclusive      bool
		NoLocal        bool
		NoWait         bool
		Args           map[string]interface{}
	}
	RsaCert struct {
		PublicKeyPath  string
		PrivateKeyPath string
	}
	// 钉钉 机器 人

	Hook struct {
		WebHook []string
	}
)

var configFile = flag.String("c", "", "Please set config file")

func ParseConfig() (*Config, error) {
	flag.Parse()

	fmt.Printf("configFile %s\n", *configFile)
	if utils.Exists(*configFile) {
		return parseConfigFormFile(*configFile)
	}

	from := os.Getenv("CONFIG_FROM")

	if from == "" {
		from = "file"
	}
	fmt.Printf(from + "-----\n")
	switch from {
	case "file":
		return parseConfigFormFile(os.Getenv("CONFIG_FILE"))
	case "env":
		return parseConfigFormEnv()
	}
	return nil, errors.New("config is error")
}

func parseConfigFormFile(filePath string) (*Config, error) {
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	conf := new(Config)
	if err := json.Unmarshal(body, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func parseConfigFormEnv() (*Config, error) {
	host := os.Getenv("MQ_HOST")
	queueName := os.Getenv("QUEUE_NAME")
	exchange := os.Getenv("EXCHANGE_NAME")
	fmt.Printf("MQ_HOST : %s ,QUEUE_NAME: %s , EXCHANGE_NAME : %s", host, queueName, exchange)
	if host == "" || queueName == "" || exchange == "" {
		return nil, errors.New("请设置环境变量MQ_HOST=xxxxxx:xx ,QUEUE_NAME=xx, EXCHANGE_NAME=xxx")
	}

	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "logs"
	}

	consumer := os.Getenv("MQ_CONSUMER_NAME")
	if consumer == "" {
		consumer = "consumer"
	}
	consumerAmount, _ := strconv.Atoi(os.Getenv("MQ_CONSUMER_NUM"))
	if consumerAmount == 0 {
		consumerAmount = 3
	}
	autoAck, _ := strconv.ParseBool(os.Getenv("MQ_AUTO_ACK"))
	durable, _ := strconv.ParseBool(os.Getenv("MQ_DURABLE"))     // Durable
	autoDelete, _ := strconv.ParseBool(os.Getenv("MQ_AUTO_DEL")) // AutoDelete
	exclusive, _ := strconv.ParseBool(os.Getenv("MQ_EXCLUSIVE")) // Exclusive
	noLocal, _ := strconv.ParseBool(os.Getenv("MQ_NO_LOCAL"))    // NoLocal
	noWait, _ := strconv.ParseBool(os.Getenv("MQ_NO_WAIT"))      // NoWait
	hooks := strings.Split(",", os.Getenv("WEB_HOOK"))
	log4gStdout, _ := strconv.ParseBool(os.Getenv("LOG_STDOUT"))

	return &Config{
		Log4g: log4g.Config{
			Path:   logPath,
			Stdout: log4gStdout,
		},
		RabbitMq: RabbitMq{
			DataSource:     host,
			QueueName:      queueName,
			Consumer:       consumer,
			ConsumerAmount: consumerAmount,
			Exchange:       exchange,
			AutoAck:        autoAck,
			Durable:        durable,
			AutoDelete:     autoDelete,
			Exclusive:      exclusive,
			NoLocal:        noLocal,
			NoWait:         noWait,
		},
		Hook: Hook{
			WebHook: hooks,
		},
	}, nil
}
