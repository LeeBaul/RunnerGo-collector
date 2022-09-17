package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"kp-collector/internal/pkg/dal/kao"
	log2 "kp-collector/internal/pkg/log"
	"log"
	"os"
	"time"
)

var Client *elastic.Client
var Exist bool

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

func InsertTestData(sceneTestResultDataMsg *kao.SceneTestResultDataMsg, index string) {

	if Exist {
		_, err := Client.Index().Index(index).BodyJson(sceneTestResultDataMsg).Do(context.Background())
		if err != nil {
			log2.Logger.Error("es写入数据失败", err)
			return
		}
	}
	if !Exist {
		_, err := Client.CreateIndex(index).BodyJson(sceneTestResultDataMsg).Do(context.Background())
		if err != nil {
			log2.Logger.Error("es写入数据失败", err)
			return
		}
		Exist = true
	}

}

func Exists(index string) bool {
	exist, err := Client.IndexExists(index).Do(context.Background())
	if err != nil {
		panic(fmt.Sprintf("es连接失败: %s", err))
	}
	return exist
}
