package main

import (
	"kp-collector/internal"
	"kp-collector/internal/pkg/conf"
	log2 "kp-collector/internal/pkg/log"
	"kp-collector/internal/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	internal.InitProjects()

	go server.Execute(conf.Conf.Kafka.Host)

	/// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log2.Logger.Info("注销成功")

}
