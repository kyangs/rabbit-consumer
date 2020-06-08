package config

import "github.com/yakaa/log4g"

//amqp.Dial("amqp://guest:guest@localhost:5672/")
type (
	Config struct {
		Log4g      log4g.Config
		RabbitMq   RabbitMq
		MpsMysql   Mysql
		AmqpMysql  Mysql
		ErpMysql   Mysql
		RomeoMysql Mysql
		RsaCert    RsaCert
		DingTalk   DingTalk
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
	Mysql struct {
		DataSource string
	}
	// 钉钉 机器 人

	DingTalk struct {
		WebHook []string
	}
)
