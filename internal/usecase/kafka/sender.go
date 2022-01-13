package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zhenianik/grpcApiTest/internal/model"
)

type es struct {
	kafka *kafka.Writer
}

func New(address string) *es {
	w := &kafka.Writer{
		Addr:  kafka.TCP(address),
		Async: true,
	}

	return &es{
		kafka: w,
	}
}

type UserEvent struct {
	ID   int64  `json:"user_id"`
	Time int64  `json:"time"`
	Name string `json:"name"`
}

func (es *es) Send(userID model.UserID, name string, time time.Time) error {
	dataJSON, err := json.Marshal(UserEvent{
		ID:   int64(userID),
		Time: time.Unix(),
		Name: name,
	})
	if err != nil {
		return err
	}

	return es.kafka.WriteMessages(context.Background(), kafka.Message{
		Topic: "users",
		Value: dataJSON,
	})
}
