package rabbitmq

import (
	"encoding/json"
	"os"
	"os/signal"

	"consumer/config"
	"github.com/streadway/amqp"
	"github.com/yakaa/log4g"
)

type (
	Consumer struct {
		amqpDial     *amqp.Connection
		amqpDialCh   *amqp.Channel
		stop         chan bool
		consumerFunc ConsumerFunc
		conf         config.RabbitMq
	}
	ConsumerFunc func(message *Message) error
)

func BuildConsumer(conf config.RabbitMq, consumerFunc ConsumerFunc) (*Consumer, error) {
	amqpDial, err := amqp.Dial(conf.DataSource)
	if err != nil {
		return nil, err
	}
	ch, err := amqpDial.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{
		amqpDial:     amqpDial,
		stop:         make(chan bool),
		consumerFunc: consumerFunc,
		amqpDialCh:   ch,
		conf:         conf,
	}, nil
}
func (c *Consumer) StartConsume() error {
	if err := c.amqpDialCh.Qos(1, 0, false); err != nil {
		return err
	}
	response, err := c.amqpDialCh.Consume(
		c.conf.QueueName,
		c.conf.Consumer,
		c.conf.AutoAck,
		c.conf.Exclusive,
		c.conf.NoLocal,
		c.conf.NoWait,
		QueueDelayedTable,
	)
	if err != nil {
		return err
	}
	go func() {
		for d := range response {
			message := new(Message)
			if err := json.Unmarshal(d.Body, message); err != nil {
				log4g.ErrorFormat("Err Message format %+v", err)
				continue
			}
			log4g.InfoFormat("start Consume message %+v", message)
			if err := c.consumerFunc(message); err != nil {
				log4g.ErrorFormat("Consume Message Err %+v", err)
				continue
			}
		}
	}()
	<-c.stop
	return nil
}

func (c *Consumer) Run() error {
	log4g.InfoFormat("consumer start run..., listen queue name is %s", c.conf.QueueName)
	c.Close()
	if err := c.StartConsume(); err != nil {
		return err
	}
	return nil
}

func (c *Consumer) Close() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, DeadSignal...)
	go func() {
		log4g.InfoFormat("Consumer receive dead signal %+v ", <-ch)
		if err := c.amqpDialCh.Close(); err != nil {
			log4g.ErrorFormat("c.amqpDialCh.Close err %+v", err)
		}
		if err := c.amqpDial.Close(); err != nil {
			log4g.InfoFormat("Consumer conn Close err %+v by receive dead signal", err)
		}
		c.stop <- true
		os.Exit(1)
	}()
}
