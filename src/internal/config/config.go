package config

type Consumer struct {
	Brokers string `yaml:"env" env-required:"true"`
	GroupId string `yaml:"consumer_group" env-required:"true"`
	Topic   string `yaml:"message_topic" env-required:"true"`
}
