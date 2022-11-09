package main

import (
	"RunnerGo-collector/internal"
	"RunnerGo-collector/internal/pkg/conf"
	log2 "RunnerGo-collector/internal/pkg/log"
	"RunnerGo-collector/internal/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	internal.InitProjects()

	collectorService := &http.Server{
		Addr: conf.Conf.Http.Host,
	}
	go server.Execute(conf.Conf.Kafka.Host)

	go func() {
		if err := collectorService.ListenAndServe(); err != nil {
			log2.Logger.Error("collector:", err)
			return
		}
	}()

	/// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log2.Logger.Info("注销成功")

}
