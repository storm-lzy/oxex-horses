package mq

import (
	"encoding/json"
	"log"

	"niuma-house/internal/model"
	"niuma-house/pkg/database"
	"niuma-house/pkg/queue"

	"gorm.io/gorm"
)

// StartExpConsumer 启动经验值消费者
func StartExpConsumer() {
	ch := queue.GetChannel()

	msgs, err := ch.Consume(
		"exp_queue", // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("Exp consumer started, waiting for messages...")

	for msg := range msgs {
		var expMsg ExpMessage
		if err := json.Unmarshal(msg.Body, &expMsg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			msg.Nack(false, false)
			continue
		}

		if err := processExpMessage(expMsg); err != nil {
			log.Printf("Failed to process exp message: %v", err)
			msg.Nack(false, true) // requeue
			continue
		}

		msg.Ack(false)
		log.Printf("Processed exp message: userID=%d, action=%s, exp=+%d",
			expMsg.UserID, expMsg.Action, expMsg.ExpAmount)
	}
}

// processExpMessage 处理经验值消息
func processExpMessage(msg ExpMessage) error {
	db := database.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		// 更新经验值
		if err := tx.Model(&model.User{}).
			Where("id = ?", msg.UserID).
			UpdateColumn("exp", gorm.Expr("exp + ?", msg.ExpAmount)).Error; err != nil {
			return err
		}

		// 获取更新后的用户
		var user model.User
		if err := tx.First(&user, msg.UserID).Error; err != nil {
			return err
		}

		// 计算新等级
		newLevel := model.CalculateLevel(user.Exp)
		if newLevel != user.Level {
			if err := tx.Model(&user).Update("level", newLevel).Error; err != nil {
				return err
			}
			log.Printf("User %d level up: %d -> %d", msg.UserID, user.Level, newLevel)
		}

		return nil
	})
}
