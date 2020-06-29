package db

import (
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
	kafka "github.com/Shopify/sarama"
)

var KafkaPublisherInstance *KafkaPublisher

func SetupKafkaPublisher(brokerURL string, topic string, partition int32) error {
	KafkaPublisherInstance = MakeKafkaPublisher(brokerURL, topic, partition)
	err := KafkaPublisherInstance.Start()
	return err
}

type KafkaPublisher struct {
	BrokerURL string
	Topic     string
	Partition int32
	Offset    int64
	Writer    kafka.SyncProducer
}

func MakeKafkaPublisher(brokerURL string, topic string, partition int32) *KafkaPublisher {
	item := &KafkaPublisher{
		BrokerURL: brokerURL,
		Topic:     topic,
		Partition: partition,
		Offset:    0,
	}
	return item
}

func (p *KafkaPublisher) Start() error {
	var err error
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_11_0_0

	p.Writer, err = sarama.NewSyncProducer([]string{p.BrokerURL}, config)
	if err != nil {
		fmt.Errorf("failed to start kafkaPublisher:%v\n", err)
	}
	return err
}

func (p *KafkaPublisher) Stop() error {
	var err error
	if p.Writer != nil {
		err = p.Writer.Close()
	}
	return err
}

func (p *KafkaPublisher) SendPacket(msg []byte) error {
	var err error
	if p.Writer == nil {
		err = errors.New("kafka writer is nil")
	}
	kfmsg := &kafka.ProducerMessage{Topic: p.Topic, Partition: p.Partition, Value: kafka.StringEncoder(msg)}
	_, p.Offset, err = p.Writer.SendMessage(kfmsg)
	return err
}

func (p *KafkaPublisher) SendTopicPacket(topic string, msg []byte) error {
	var err error
	if p.Writer == nil {
		err = errors.New("kafka writer is nil")
	}
	kfmsg := &kafka.ProducerMessage{Topic: topic, Partition: p.Partition, Value: kafka.StringEncoder(msg)}
	_, p.Offset, err = p.Writer.SendMessage(kfmsg)
	return err
}
