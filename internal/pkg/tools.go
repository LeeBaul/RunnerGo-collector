package pkg

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	services "kp-collector/api"
	log2 "kp-collector/internal/pkg/log"
	"net/http"
	"strconv"
	"strings"
)

// SendStopMsg 发送结束任务消息
func SendStopMsg(host, reportId string) {
	ctx := context.TODO()

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		panic(errors.Wrap(err, "cannot load root CA certs"))
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(creds))
	defer func() {
		grpcErr := conn.Close()
		if grpcErr != nil {
			log2.Logger.Error("关闭grpc连接失败:", grpcErr)
		}
	}()
	grpcClient := services.NewKpControllerClient(conn)
	req := new(services.NotifyStopStressReq)
	log2.Logger.Error("reportId", reportId)
	req.ReportID, err = strconv.ParseInt(reportId, 10, 64)
	if err != nil {
		log2.Logger.Error("reportId转换失败", err)
		return
	}

	_, err = grpcClient.NotifyStopStress(ctx, req)
	if err != nil {
		log2.Logger.Error("发送停止任务失败", err)
		return
	}
	log2.Logger.Info(reportId, "任务结束， 消息已发送")
}

func Post(url, body string) (err error) {

	request, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		fmt.Println("http请求创建失败：   ", err)
		return
	}
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		fmt.Println("http发送请求失败：    ", err)
		return
	}
	return
}
