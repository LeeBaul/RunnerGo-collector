package server

import (
	"github.com/Shopify/sarama"
	"kp-collector/internal/pkg/conf"
	"kp-collector/internal/pkg/dal/es"
	log2 "kp-collector/internal/pkg/log"
	"sync"
)

func Execute(host string) {
	var topicMap sync.Map
	var consumerMap = &sync.Map{}
	// 创建client
	client, err := sarama.NewClient([]string{host}, sarama.NewConfig())
	if err != nil {
		log2.Logger.Error("创建kafka客户端失败:", err)
		return
	}

	defer func() {
		if clientErr := client.Close(); clientErr != nil {
			log2.Logger.Error("关闭kafka客户端失败:", clientErr)
		}
	}()
	for {
		// 获取所有的topic
		topics, topicsErr := client.Topics()
		if topicsErr != nil {
			log2.Logger.Error("获取topics失败：", err)
			continue
		}
		index := 0
		for _, topic := range topics {
			if _, ok := topicMap.Load(topic); !ok {
				consumer, consumerErr := sarama.NewConsumer([]string{host}, sarama.NewConfig())
				if consumerErr != nil {
					log2.Logger.Error("topic  :"+topic+", 创建消费者失败:", consumerErr)
					continue
				}
				partitions, consumerErr := consumer.Partitions(topic)
				if consumerErr != nil {
					log2.Logger.Error("获取", conf.Conf.Kafka.Topic, "主题失败：", consumerErr)
					if consumerErr = consumer.Close(); consumerErr != nil {
						log2.Logger.Error("关闭消费者失败：", consumerErr)
					}
					continue
				}
				index++
				topicMap.Store(topic, true)
				for _, p := range partitions {
					if v, partitionOk := consumerMap.Load(p); !partitionOk {
						consumerMap.Store(v, true)
						partition, partitionConsumerErr := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
						if partitionConsumerErr != nil {
							log2.Logger.Error("创建消费者失败：", partitionConsumerErr)
							consumerMap.Delete(p)
							continue
						}
						go ReceiveMessage(partition, consumerMap, p)
					}

				}

			}
		}
	}
}

func ReceiveMessage(partition sarama.PartitionConsumer, consumerMap *sync.Map, p int32) {
	defer partition.Close()
	defer consumerMap.Delete(p)
	if partition == nil || consumerMap == nil {
		return
	}
	var machineMap *sync.Map
	for msg := range partition.Messages() {
		err := es.InsertTestData(msg.Value, machineMap)
		if err != nil {
			log2.Logger.Error("es写入失败:", err)
			return
		}
	}
}
