package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// RabbitMQ RabbitMQ连接封装
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// QueueConfig 队列配置
type QueueConfig struct {
	Name         string
	Durable      bool
	DeleteUnused bool
	Exclusive    bool
	NoWait       bool
	Args         amqp.Table
}

// NewRabbitMQ 创建RabbitMQ连接
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ连接失败: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建Channel失败: %w", err)
	}

	logx.Info("RabbitMQ连接成功")
	return &RabbitMQ{conn: conn, channel: channel}, nil
}

// Close 关闭连接
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// GetChannel 获取Channel（注意：不是线程安全的，多协程请使用NewChannel）
func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.channel
}

// NewChannel 创建新Channel（线程安全）
func (r *RabbitMQ) NewChannel() (*amqp.Channel, error) {
	return r.conn.Channel()
}

// DeclareExchange 声明交换机
func (r *RabbitMQ) DeclareExchange(name, kind string) error {
	return r.channel.ExchangeDeclare(
		name,  // name
		kind,  // type: direct, fanout, topic, headers
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // args
	)
}

// DeclareQueue 声明队列
func (r *RabbitMQ) DeclareQueue(config QueueConfig) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		config.Name,
		config.Durable,
		config.DeleteUnused,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)
}

// DeclareDLXQueue 声明带死信队列的普通队列
func (r *RabbitMQ) DeclareDLXQueue(queueName, dlxExchange, dlxQueue string) error {
	// 1. 声明死信交换机
	err := r.DeclareExchange(dlxExchange, "direct")
	if err != nil {
		return fmt.Errorf("声明死信交换机失败: %w", err)
	}

	// 2. 声明死信队列
	dlxQ, err := r.DeclareQueue(QueueConfig{
		Name:    dlxQueue,
		Durable: true,
	})
	if err != nil {
		return fmt.Errorf("声明死信队列失败: %w", err)
	}

	// 3. 绑定死信队列到死信交换机
	err = r.channel.QueueBind(dlxQ.Name, dlxQueue, dlxExchange, false, nil)
	if err != nil {
		return fmt.Errorf("绑定死信队列失败: %w", err)
	}

	// 4. 声明主队列，配置死信路由
	args := amqp.Table{
		"x-dead-letter-exchange":    dlxExchange,
		"x-dead-letter-routing-key": dlxQueue,
	}
	_, err = r.DeclareQueue(QueueConfig{
		Name:    queueName,
		Durable: true,
		Args:    args,
	})
	if err != nil {
		return fmt.Errorf("声明主队列失败: %w", err)
	}

	return nil
}

// PublishMessage 发布消息
func (r *RabbitMQ) PublishMessage(ctx context.Context, exchange, routingKey string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}

	err = r.channel.PublishWithContext(ctx,
		exchange,     // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent, // 持久化
			Timestamp:    time.Now(),
			MessageId:    fmt.Sprintf("%d", time.Now().UnixNano()),
		},
	)
	if err != nil {
		return fmt.Errorf("发布消息失败: %w", err)
	}

	logx.Infof("消息已发布: exchange=%s, routingKey=%s", exchange, routingKey)
	return nil
}

// PublishWithRetry 带重试的消息发布（指数退避）
func (r *RabbitMQ) PublishWithRetry(ctx context.Context, exchange, routingKey string, body interface{}, maxRetry int) error {
	var lastErr error
	for i := 0; i <= maxRetry; i++ {
		err := r.PublishMessage(ctx, exchange, routingKey, body)
		if err == nil {
			return nil
		}
		lastErr = err
		if i < maxRetry {
			// 指数退避：1s, 2s, 4s, 8s, 16s...
			backoff := time.Duration(1<<uint(i)) * time.Second
			logx.Infof("消息发布失败，%v后重试 (第%d次): %v", backoff, i+1, err)
			time.Sleep(backoff)
		}
	}
	return fmt.Errorf("消息发布失败，已达最大重试次数%d: %w", maxRetry, lastErr)
}

// ConsumeMessages 消费消息
func (r *RabbitMQ) ConsumeMessages(queueName string, handler func(body []byte) error) error {
	msgs, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // autoAck (手动确认)
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("启动消费者失败: %w", err)
	}

	go func() {
		for msg := range msgs {
			logx.Infof("收到消息: queue=%s, messageId=%s", queueName, msg.MessageId)
			err := handler(msg.Body)
			if err != nil {
				logx.Errorf("消息处理失败: %v", err)
				// 拒绝消息，重新入队或进入死信队列
				if msg.DeliveryTag > 0 {
					_ = msg.Nack(false, false) // false=不重新入队，进入死信队列
				}
			} else {
				// 手动确认
				if msg.DeliveryTag > 0 {
					_ = msg.Ack(false)
				}
			}
		}
	}()

	logx.Infof("消费者已启动: queue=%s", queueName)
	return nil
}

// SetQoS 设置预取数量
func (r *RabbitMQ) SetQoS(prefetchCount int) error {
	return r.channel.Qos(prefetchCount, 0, false)
}

// DiagnosisMessage 诊断消息体
type DiagnosisMessage struct {
	TaskID    uint64 `json:"task_id"`
	TaskNo    string `json:"task_no"`
	UserID    uint64 `json:"user_id"`
	ImagePath string `json:"image_path"`
	ImageURL  string `json:"image_url"`
	Age       int32  `json:"age,omitempty"`
	Gender    int8   `json:"gender,omitempty"`
	// 幂等键
	IdempotentKey string `json:"idempotent_key"`
	// 重试次数
	RetryCount int `json:"retry_count"`
	// 时间戳
	Timestamp int64 `json:"timestamp"`
}
