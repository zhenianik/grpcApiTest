package logger

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

type Logger struct {
	kafka *kafka.Writer
}

func New(address string) *Logger {
	w := &kafka.Writer{
		Addr:  kafka.TCP(address),
		Async: true,
	}

	return &Logger{
		kafka: w,
	}
}

type UserLogRequest struct {
	ID   int64  `json:"user_id"`
	Time int64  `json:"time"`
	Name string `json:"name"`
}

func (l *Logger) LogRegistration(userID int64, name string, time time.Time) error {
	dataJSON, err := json.Marshal(UserLogRequest{
		ID:   userID,
		Time: time.Unix(),
		Name: name,
	})
	if err != nil {
		return err
	}

	return l.kafka.WriteMessages(context.Background(), kafka.Message{
		Topic: "users",
		Value: dataJSON,
	})
}
