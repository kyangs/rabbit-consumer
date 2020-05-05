package rabbitmq

import (
	"encoding/json"
	"os"
	"os/signal"
	"strconv"

	"consumer/config"
	"github.com/streadway/amqp"
	"github.com/yakaa/log4g"
)

type (
	Consumer struct {
		amqpDial     *amqp.Connection
		amqpDialCh   *amqp.Channel
		consumerFunc ConsumerFunc
		consumerName string
		conf         config.RabbitMq
	}
	ConsumerFunc func(message *Message) error

	ConsumerPool struct {
		Pool []*Consumer
		stop chan bool
	}
)

func BuildConsumerPool(conf config.RabbitMq, consumerFunc ConsumerFunc, amount int) (*ConsumerPool, error) {
	cp := &ConsumerPool{Pool: []*Consumer(nil), stop: make(chan bool)}
	for i := 0; i < amount; i++ {
		c, err := buildConsumer(conf, consumerFunc, conf.Consumer+"_"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		cp.put(c)
	}
	return cp, nil
}

func (cp *ConsumerPool) put(c *Consumer) {
	cp.Pool = append(cp.Pool, c)
}

func (cp *ConsumerPool) Run() error {
	log4g.InfoFormat("consumer pool start run...")
	for _, c := range cp.Pool {
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
			log4g.ErrorFormat("%s create consumer channel  fail %+v", c.consumerName, err)
			return err
		}
		log4g.InfoFormat("%s created ", c.consumerName)
		go func(consumer *Consumer, m <-chan amqp.Delivery) {
			for d := range m {
				message := new(Message)
				if err := json.Unmarshal(d.Body, message); err != nil {
					log4g.ErrorFormat("%s Err Message format %+v", consumer.consumerName, err)
					continue
				}
				log4g.InfoFormat("%s start Consume message %+v", consumer.consumerName, message)
				if err := c.consumerFunc(message); err != nil {
					log4g.ErrorFormat("%s Consume Message err %+v", consumer.consumerName, err)
					continue
				}
			}
		}(c, response)
	}
	cp.Close()
	<-cp.stop
	return nil
}

func (cp *ConsumerPool) Close() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, DeadSignal...)
	go func() {
		log4g.InfoFormat("ConsumerPool receive dead signal %+v ", <-ch)
		for _, consumer := range cp.Pool {
			if err := consumer.amqpDialCh.Close(); err != nil {
				log4g.ErrorFormat("%s ConsumerPool amqpDialCh.Close err %+v ", consumer.consumerName, err)
			} else {
				log4g.InfoFormat("%s ConsumerPool  Close channel success", consumer.consumerName)
			}

			if err := consumer.amqpDial.Close(); err != nil {
				log4g.ErrorFormat("%s ConsumerPool conn Close err %+v by receive dead signal", consumer.consumerName, err)
			} else {
				log4g.InfoFormat("%s ConsumerPool  Close Dial success", consumer.consumerName)
			}
		}
		cp.stop <- true
		os.Exit(1)
	}()
}

func buildConsumer(conf config.RabbitMq, consumerFunc ConsumerFunc, consumerName string) (*Consumer, error) {
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
		consumerFunc: consumerFunc,
		amqpDialCh:   ch,
		consumerName: consumerName,
		conf:         conf,
	}, nil
}
