package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"niuma-house/pkg/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 活动类型
const (
	ActionLogin     = "login"
	ActionPost      = "post"
	ActionLiked     = "liked"
	ActionCommented = "commented"
)

// ExpMessage 经验值消息
type ExpMessage struct {
	UserID    uint   `json:"user_id"`
	Action    string `json:"action"`
	ExpAmount int    `json:"exp_amount"`
	Timestamp int64  `json:"timestamp"`
}

// PublishExpMessage 发布经验值消息
func PublishExpMessage(userID uint, action string, expAmount int) error {
	ch := queue.GetChannel()

	msg := ExpMessage{
		UserID:    userID,
		Action:    action,
		ExpAmount: expAmount,
		Timestamp: time.Now().Unix(),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal exp message: %v", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"user_activity", // exchange
		"exp",           // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Printf("Failed to publish exp message: %v", err)
		return err
	}

	log.Printf("Published exp message: userID=%d, action=%s, exp=%d", userID, action, expAmount)
	return nil
}
