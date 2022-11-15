package pkg

import (
	services "RunnerGo-collector/api"
	"RunnerGo-collector/internal/pkg/conf"
	log2 "RunnerGo-collector/internal/pkg/log"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
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
	log2.Logger.Info(reportId, "   任务结束， 消息已发送")
}

func Post(url, body string) (err error) {
	strs := strings.Split(url, "?")
	url = strs[0]
	request, err := http.NewRequest("GET", url, strings.NewReader(body))
	querys := ""
	if len(strs) > 0 {
		querys = strs[1]
	}
	queryList := strings.Split(querys, "&")
	query := request.URL.Query()
	for i := 0; i < len(queryList); i++ {
		s := strings.Split(queryList[i], "=")
		query.Add(s[0], s[1])
	}
	request.URL.RawQuery = query.Encode()

	if err != nil {
		fmt.Println("http请求创建失败：   ", err)
		return
	}
	fmt.Println("url:    ", url)
	client := &http.Client{}
	resp, err := client.Do(request)
	fmt.Println("response:     ", resp.Body)
	if err != nil {
		fmt.Println("http发送请求失败：    ", err)
		return
	}
	return
}

type StopMsg struct {
	ReportId int64    `json:"report_id"`
	Machines []string `json:"machines"`
}

func SendStopStressReport(machineMap map[string]map[string]int64, reportId int64) {

	sm := StopMsg{
		ReportId: reportId,
	}
	for k, _ := range machineMap {
		sm.Machines = append(sm.Machines, k)
	}

	body, err := json.Marshal(&sm)
	if err != nil {
		log2.Logger.Error(reportId, "   ,json转换失败：  ", err.Error())
	}
	res, err := http.Post(conf.Conf.Management.Address, "application/json", strings.NewReader(string(body)))
	defer res.Body.Close()
	if err != nil {
		log2.Logger.Error("http请求建立链接失败：", err.Error())
		return
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log2.Logger.Error("http读取响应信息失败：", err.Error())
		return
	}
	if strings.Contains(string(responseBody), "\"code\":0,") {
		log2.Logger.Info(reportId, "  :报告停止任务成功： ", "请求体：", string(body), "            响应体：   ", string(responseBody))
	} else {
		log2.Logger.Error(reportId, "  :报告停止任务失败：  ", "请求体：", string(body), "          响应体：   ", string(responseBody))
	}

}
