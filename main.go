package main

import (
	"github.com/Shopify/sarama"
	"kp-collector/internal"
	"kp-collector/internal/pkg/conf"
	log2 "kp-collector/internal/pkg/log"
	"kp-collector/internal/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	go server.Execute(conf.Conf.Kafka.Topic, consumer)

	/// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log2.Logger.Info("注销成功")

}
