package kafka_consumer

type Config struct {
	BrokerAddress string
	GroupID       string
	Topic         string
}

func NewConfig(brokerAddress, groupID, topic string) *Config {
	return &Config{
		BrokerAddress: brokerAddress,
		GroupID:       groupID,
		Topic:         topic,
	}
}
