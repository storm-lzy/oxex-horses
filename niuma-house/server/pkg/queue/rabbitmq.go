package queue

import (
	"fmt"
	"log"
	"sync"

	"niuma-house/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
	once    sync.Once
)

// InitRabbitMQ 初始化 RabbitMQ 连接
func InitRabbitMQ(cfg *config.RabbitMQConfig) *amqp.Channel {
	once.Do(func() {
		url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.VHost,
		)

		var err error
		conn, err = amqp.Dial(url)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}

		channel, err = conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open RabbitMQ channel: %v", err)
		}

		// 声明交换机
		err = channel.ExchangeDeclare(
			"user_activity", // name
			"direct",        // type
			true,            // durable
			false,           // auto-deleted
			false,           // internal
			false,           // no-wait
			nil,             // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare exchange: %v", err)
		}

		// 声明队列
		_, err = channel.QueueDeclare(
			"exp_queue", // name
			true,        // durable
			false,       // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare queue: %v", err)
		}

		// 绑定队列到交换机
		err = channel.QueueBind(
			"exp_queue",     // queue name
			"exp",           // routing key
			"user_activity", // exchange
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to bind queue: %v", err)
		}

		log.Println("RabbitMQ connected successfully")
	})
	return channel
}

// GetChannel 获取 RabbitMQ Channel 单例
func GetChannel() *amqp.Channel {
	if channel == nil {
		log.Fatal("RabbitMQ not initialized. Call InitRabbitMQ first.")
	}
	return channel
}

// Close 关闭连接
func Close() {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
