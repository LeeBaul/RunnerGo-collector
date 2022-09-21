package es

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	services "kp-collector/api"
	"kp-collector/internal/pkg/conf"
	"kp-collector/internal/pkg/dal/kao"
	log2 "kp-collector/internal/pkg/log"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var Client *elastic.Client

func InitEsClient(host, user, password string) {
	Client, _ = elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(user, password),
		elastic.SetErrorLog(log.New(os.Stdout, "APP", log.Lshortfile)),
		elastic.SetHealthcheckInterval(30*time.Second),
	)
	_, _, err := Client.Ping(host).Do(context.Background())
	if err != nil {
		panic(fmt.Sprintf("es连接失败: %s", err))
	}
	return
}

func InsertTestData(value []byte, machineMap *sync.Map) (err error) {
	var resultDataMsg = kao.ResultDataMsg{}
	err = json.Unmarshal(value, &resultDataMsg)
	fmt.Println("resultDataMsg", resultDataMsg.ReportId)
	if err != nil {
		log2.Logger.Error("kafka消息转换失败：", err)
		return
	}
	if resultDataMsg.ReportId == "" {
		log2.Logger.Error(fmt.Sprintf("es连接失败: %s", err))
		return errors.New("报告id不能为空")
	}
	index := conf.Conf.ES.Index
	exist, err := Client.IndexExists(index).Do(context.Background())
	if err != nil {
		log2.Logger.Error(fmt.Sprintf("es连接失败: %s", err))
		return
	}
	if !exist {
		_, clientErr := Client.CreateIndex(index).Do(context.Background())
		if clientErr != nil {
			log2.Logger.Error("es创建索引", index, "失败", err)
			return
		}
	}
	if machine, ok := machineMap.Load(resultDataMsg.ReportId); !ok {
		machineMap.Store(resultDataMsg.ReportId, resultDataMsg.MachineNum+1)

	} else {
		if resultDataMsg.End == true && machine.(int64) > 1 {
			machineMap.Store(resultDataMsg.ReportId, machine.(int64)-1)
		}
		machineNum, is := machineMap.Load(resultDataMsg.ReportId)
		if is && machineNum.(int64) == 1 && resultDataMsg.End == true {
			log2.Logger.Info(resultDataMsg.ReportId, "报告结束")
			SendStopMsg(conf.Conf.GRPC.Host, resultDataMsg.ReportId)
		}
	}
	_, err = Client.Index().Index(index).BodyString(string(value)).Do(context.Background())
	if err != nil {
		log2.Logger.Error("es写入数据失败", err)
		return
	}
	return
}

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
	req.ReportID, _ = strconv.ParseInt(reportId, 10, 64)

	_, err = grpcClient.NotifyStopStress(ctx, req)
	if err != nil {
		log2.Logger.Error("发送停止任务失败", err)
		return
	}
	log2.Logger.Info(reportId, "任务结束， 消息已发送")
}
