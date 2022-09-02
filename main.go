package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"

	"kp-collector/internal"
	"kp-collector/internal/pkg/conf"
)

func main() {
	internal.InitProjects()

	consumer, err := sarama.NewConsumer([]string{conf.Conf.Kafka.Host}, sarama.NewConfig())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partition, err := consumer.ConsumePartition(conf.Conf.Kafka.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partition.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	fmt.Println("kafka consumer initialized")

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0

ConsumerLoop:
	for {
		select {
		case msg := <-partition.Messages():
			log.Printf("%+v", msg)
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
