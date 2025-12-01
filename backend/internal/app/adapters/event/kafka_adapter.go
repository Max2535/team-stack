package event

import (
    "context"

    "github.com/example/team-stack/backend/internal/app/ports"
    "github.com/segmentio/kafka-go"
)

type KafkaAdapter struct {
    writer *kafka.Writer
}

func NewKafka(brokers []string, defaultTopic string) ports.EventBus {
    return &KafkaAdapter{
        writer: &kafka.Writer{
            Addr:     kafka.TCP(brokers...),
            Topic:    defaultTopic,
            Balancer: &kafka.LeastBytes{},
        },
    }
}

func (k *KafkaAdapter) Publish(subject string, data []byte) error {
    msg := kafka.Message{
        Topic: subject,
        Value: data,
    }
    return k.writer.WriteMessages(context.Background(), msg)
}
