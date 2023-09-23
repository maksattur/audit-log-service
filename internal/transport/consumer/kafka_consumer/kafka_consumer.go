package kafka_consumer

import (
	"context"
	"encoding/json"
	"github.com/maksattur/audit-log-service/internal/transport/consumer"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(config *Config) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.BrokerAddress},
		GroupID: config.GroupID,
		Topic:   config.Topic,
	})

	return &KafkaConsumer{
		reader: reader,
	}
}

func (kc *KafkaConsumer) Receive(ctx context.Context, eventChan chan<- consumer.Event, errChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			close(eventChan)
			close(errChan)
			return
		default:
			msg, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				errChan <- err
				continue
			}
			var event consumer.Event
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				errChan <- err
				continue
			}
			eventChan <- event
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}
