package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zhenianik/grpcApiTest/internal"
)

type eventSender struct {
	kafka *kafka.Writer
}

func New(address string) *eventSender {
	w := &kafka.Writer{
		Addr:  kafka.TCP(address),
		Async: true,
	}

	return &eventSender{
		kafka: w,
	}
}

type UserEvent struct {
	ID   int64  `json:"user_id"`
	Time int64  `json:"time"`
	Name string `json:"name"`
}

func (es *eventSender) Send(userID internal.UserID, name string, time time.Time) error {
	dataJSON, err := json.Marshal(UserEvent{
		ID:   int64(userID),
		Time: time.Unix(),
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("marshal user event error: %w", err)
	}

	return es.kafka.WriteMessages(context.Background(), kafka.Message{
		Topic: "users",
		Value: dataJSON,
	})
}
