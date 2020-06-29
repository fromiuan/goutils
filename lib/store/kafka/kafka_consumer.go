package db

import (
	"github.com/Shopify/sarama"
)

type KafkaMessageConsumerFunc func(int64, string, []byte) error

var KafkaConsumerInstance *KafkaConsumer

func SetupKafkaConsumer(brokerURL string, topic string, partition int32,
	offset int64, handler KafkaMessageConsumerFunc) error {
	KafkaConsumerInstance = MakeKafkaConsumer(brokerURL, topic, partition, offset, handler)
	err := KafkaConsumerInstance.Start()
	return err
}

type KafkaConsumer struct {
	Running   bool
	BrokerURL string
	Topic     string
	Partition int32
	Offset    int64
	Handler   KafkaMessageConsumerFunc
}

func MakeKafkaConsumer(brokerURL string, topic string, partition int32,
	offset int64, handler KafkaMessageConsumerFunc) *KafkaConsumer {
	myoffset := offset
	//-1: start from the latest; -2: start from zero;
	if offset == -1 {
		myoffset = sarama.OffsetNewest
	} else if offset == -2 {
		myoffset = sarama.OffsetOldest
	}

	item := &KafkaConsumer{
		Running:   false,
		BrokerURL: brokerURL,
		Topic:     topic,
		Partition: partition,
		Offset:    myoffset,
		Handler:   handler,
	}
	return item
}

func (p *KafkaConsumer) Start() error {
	config := sarama.NewConfig()
	client, err := sarama.NewClient([]string{p.BrokerURL}, config)
	if err == nil {
		consumer, err2 := sarama.NewConsumerFromClient(client)
		err = err2
		if err == nil {
			defer consumer.Close()
			p.Running = true
			pc, err3 := consumer.ConsumePartition(p.Topic, p.Partition, p.Offset)
			err = err3
			if err == nil {
				go func() {
					defer pc.Close()
					for p.Running {
						for msg := range pc.Messages() {
							p.Offset = msg.Offset
							p.Handler(msg.Offset, msg.Topic, msg.Value)
						}
					}
				}()
			}
		}
	}
	return err
}

func (p *KafkaConsumer) Stop() error {
	p.Running = false
	return nil
}

func (p *KafkaConsumer) GetOffset() int64 {
	return p.Offset
}

func (p *KafkaConsumer) IsRunning() bool {
	return p.Running
}

func (p *KafkaConsumer) Status() string {
	if p.Running {
		return "running"
	}
	return "stopped"
}
