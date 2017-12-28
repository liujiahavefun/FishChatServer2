package dao

import (
	"FishChatServer2/common/dao/kafka"
	"FishChatServer2/jobs/msg_job/conf"
)

type Kafka struct {
	Consumer *kafka.Consumer
}

func NewKafka() (k *Kafka) {
	consumer := kafka.NewConsumer(conf.Conf.KafkaConsumer)
	k = &Kafka{
		Consumer: consumer,
	}
	return
}
