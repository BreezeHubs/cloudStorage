package mq

import (
	"cloudStorage/config"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
)

func initChannel() error {
	//判断channel是否存在
	if channel != nil {
		return nil
	}

	//获得rabitmq连接
	conn, err := amqp.Dial(config.RABBIT_URL)
	if err != nil {
		return errors.Wrap(err, "dial rabbitmq error")
	}

	//打开一个channel，用于消息的发送和接收
	channel, err = conn.Channel()
	if err != nil {
		return errors.Wrap(err, "open channel error")
	}

	return nil
}

func Publish(exchange, routingKey string, msg []byte) error {
	//判断channel是否正常
	if err := initChannel(); err != nil {
		return errors.Wrap(err, "init channel error")
	}

	//发送消息
	if err := channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		}); err != nil {
		return errors.Wrap(err, "publish message error")
	}

	return nil
}
