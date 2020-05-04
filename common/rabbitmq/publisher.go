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
	Publisher struct {
		amqpDial   *amqp.Connection
		amqpDialCh *amqp.Channel
		conf       config.RabbitMq
	}
)

const (
	DelayedKind        = "x-delayed-message"
	DefaultDeliveryKey = "#"

	XDelayedTypeKey   = "x-delayed-type"
	XDelayedTypeValue = "direct"

	XDeadLetterExchangeKey   = "x-dead-letter-exchange"
	XDeadLetterExchangeValue = "delayed"
)

var (
	ExchangeTable = map[string]interface{}{
		XDelayedTypeKey: XDelayedTypeValue,
	}
	QueueDelayedTable = map[string]interface{}{
		XDeadLetterExchangeKey: XDeadLetterExchangeValue,
	}
)

func BuildPublisher(conf config.RabbitMq) (*Publisher, error) {
	amqpDial, err := amqp.Dial(conf.DataSource)
	if err != nil {
		return nil, err
	}
	if conf.Args == nil {
		conf.Args = make(map[string]interface{})
	}
	if conf.DeliveryKey == "" {
		conf.DeliveryKey = DefaultDeliveryKey
	}
	ch, err := amqpDial.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{amqpDial: amqpDial, conf: conf, amqpDialCh: ch}, nil
}

func (p *Publisher) Push(message *Message) error {
	if err := p.amqpDialCh.ExchangeDeclare(
		p.conf.Exchange,
		DelayedKind,
		p.conf.Durable,
		p.conf.AutoDelete,
		false,
		p.conf.NoWait,
		ExchangeTable,
	); err != nil {
		return err
	}
	q, err := p.amqpDialCh.QueueDeclare(
		p.conf.QueueName,
		p.conf.Durable,
		p.conf.AutoDelete,
		p.conf.Exclusive,
		p.conf.NoWait,
		p.conf.Args,
	)
	if err != nil {
		return err
	}
	if err := p.amqpDialCh.QueueBind(
		q.Name,
		p.conf.DeliveryKey,
		p.conf.Exchange,
		p.conf.NoWait,
		p.conf.Args,
	); err != nil {
		return err
	}
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Body:         body,
		Headers: map[string]interface{}{
			DefaultDelayKey: message.Delay,
		},
	}
	if err = p.amqpDialCh.Publish(p.conf.Exchange, p.conf.DeliveryKey, false, false, msg); err != nil {
		return err
	}
	return nil
}

func (p *Publisher) GetQueueName() string {
	return p.conf.QueueName
}

func (p *Publisher) Close() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, DeadSignal...)
	go func() {
		log4g.InfoFormat("Publisher receive dead signal %+v ", <-ch)
		if err := p.amqpDialCh.Close(); err != nil {
			log4g.ErrorFormat("c.amqpDialCh.Close err %+v", err)
		}
		if err := p.amqpDial.Close(); err != nil {
			log4g.InfoFormat("Publisher conn Close err %+v by receive dead signal", err)
		}
	}()
}
