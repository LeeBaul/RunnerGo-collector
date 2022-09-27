package server

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"kp-collector/internal/pkg/conf"
	"kp-collector/internal/pkg/dal/es"
	"kp-collector/internal/pkg/dal/kao"
	log2 "kp-collector/internal/pkg/log"
	"sort"
	"strconv"
	"sync"
	"time"
)

func Execute(host string) {
	var topicMap sync.Map
	var consumerMap = &sync.Map{}
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	for {
		// 创建client
		client, err := sarama.NewClient([]string{host}, saramaConfig)
		if err != nil {
			log2.Logger.Error("创建kafka客户端失败:", err)
			return
		}

		// 获取所有的topic
		topics, topicsErr := client.Topics()
		if topicsErr != nil {
			log2.Logger.Error("获取topics失败：", err)
			continue
		}
		if clientErr := client.Close(); clientErr != nil {
			log2.Logger.Error("关闭kafka客户端失败:", clientErr)
		}

		for _, topic := range topics {
			if topic == "__consumer_offsets" {
				break
			}
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
				topicMap.Store(topic, true)
				partition, partitionConsumerErr := consumer.ConsumePartition(topic, partitions[0], sarama.OffsetNewest)
				if partitionConsumerErr != nil {
					log2.Logger.Error("创建消费者失败：", partitionConsumerErr)
					continue
				}
				go ReceiveMessage(partition, consumerMap, saramaConfig, host, topic)

			}
		}
	}
}

func ReceiveMessage(partition sarama.PartitionConsumer, consumerMap *sync.Map, config *sarama.Config, host, topic string) {
	defer partition.AsyncClose()
	if partition == nil || consumerMap == nil {
		return
	}
	var requestTimeListMap = make(map[string]kao.RequestTimeList)
	var resultDataMsg = kao.ResultDataMsg{}
	var sceneTestResultDataMsg kao.SceneTestResultDataMsg
	var machineNum = int64(0)
	ticker := time.NewTicker(1 * time.Second)
Loop:
	for {
		select {
		case msg := <-partition.Messages():
			err := json.Unmarshal(msg.Value, &resultDataMsg)
			if err != nil {
				log2.Logger.Error("kafka消息转换失败：", err)
				break Loop
			}
			if resultDataMsg.ReportId == "" {
				log2.Logger.Error(fmt.Sprintf("es连接失败: %s", err))
				break Loop
			}
			if machineNum == 0 && resultDataMsg.MachineNum != 0 {
				machineNum = resultDataMsg.MachineNum + 1
			}
			if resultDataMsg.End {
				machineNum = machineNum - 1
				if machineNum == 1 {
					sceneTestResultDataMsg.End = true
					for eventId, requestTimeList := range requestTimeListMap {
						sort.Sort(requestTimeList)
						sceneTestResultDataMsg.Results[eventId].AvgRequestTime = sceneTestResultDataMsg.Results[eventId].TotalRequestTime / sceneTestResultDataMsg.Results[eventId].TotalRequestNum
						sceneTestResultDataMsg.Results[eventId].MaxRequestTime = requestTimeList[len(requestTimeList)-1]
						sceneTestResultDataMsg.Results[eventId].MinRequestTime = requestTimeList[0]
						sceneTestResultDataMsg.Results[eventId].NinetyRequestTimeLine = 90
						sceneTestResultDataMsg.Results[eventId].NinetyFiveRequestTimeLine = 95
						sceneTestResultDataMsg.Results[eventId].NinetyNineRequestTimeLine = 99
						sceneTestResultDataMsg.Results[eventId].NinetyRequestTimeLineValue = kao.TimeLineCalculate(90, requestTimeList)
						sceneTestResultDataMsg.Results[eventId].NinetyFiveRequestTimeLineValue = kao.TimeLineCalculate(95, requestTimeList)
						sceneTestResultDataMsg.Results[eventId].NinetyNineRequestTimeLineValue = kao.TimeLineCalculate(99, requestTimeList)
						sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLineValue = kao.TimeLineCalculate(sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLine, requestTimeList)
						total := float64(sceneTestResultDataMsg.Results[eventId].TotalRequestTime) / 1000
						if total <= 0 {
							total = 1
						}
						qps := fmt.Sprintf("%.2f", float64(sceneTestResultDataMsg.Results[eventId].TotalRequestNum)/total)
						sceneTestResultDataMsg.Results[eventId].Qps, _ = strconv.ParseFloat(qps, 64)
					}
					sceneTestResultDataMsg.TimeStamp = time.Now().Unix()
					if err = es.InsertTestData(sceneTestResultDataMsg); err != nil {
						log2.Logger.Error("es写入失败:", err)
						return
					}
					if err = DeleteTopic(config, host, topic); err != nil {
						return
					}
					return
				}
				if sceneTestResultDataMsg.Results == nil {
					continue
				}
				break Loop
			}
			if sceneTestResultDataMsg.SceneId == 0 {
				sceneTestResultDataMsg.SceneId = resultDataMsg.SceneId
			}
			if sceneTestResultDataMsg.SceneName == "" {
				sceneTestResultDataMsg.SceneName = resultDataMsg.SceneName
			}
			if sceneTestResultDataMsg.ReportId == "" {
				sceneTestResultDataMsg.ReportId = resultDataMsg.ReportId
			}
			if sceneTestResultDataMsg.ReportName == "" {
				sceneTestResultDataMsg.ReportName = resultDataMsg.ReportName
			}
			if sceneTestResultDataMsg.PlanId == 0 {
				sceneTestResultDataMsg.PlanId = resultDataMsg.PlanId
			}
			if sceneTestResultDataMsg.PlanName == "" {
				sceneTestResultDataMsg.PlanName = resultDataMsg.PlanName
			}
			if sceneTestResultDataMsg.Results == nil {
				sceneTestResultDataMsg.Results = make(map[string]*kao.ApiTestResultDataMsg)
			}
			if sceneTestResultDataMsg.Results[resultDataMsg.EventId] == nil {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId] = new(kao.ApiTestResultDataMsg)
			}
			if sceneTestResultDataMsg.Results[resultDataMsg.EventId].EventId == "" {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].EventId = resultDataMsg.EventId
			}
			if sceneTestResultDataMsg.Results[resultDataMsg.EventId].Name == "" {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].Name = resultDataMsg.Name
			}
			if sceneTestResultDataMsg.Results[resultDataMsg.EventId].SceneId == 0 {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].SceneId = resultDataMsg.SceneId
			}
			if sceneTestResultDataMsg.Results[resultDataMsg.EventId].PlanId == 0 {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].PlanId = resultDataMsg.PlanId
			}
			if sceneTestResultDataMsg.Results[resultDataMsg.EventId].PlanName == "" {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].PlanName = resultDataMsg.PlanName
			}
			if sceneTestResultDataMsg.Results[resultDataMsg.EventId].SceneName == "" {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].SceneName = resultDataMsg.SceneName
			}
			if resultDataMsg.PercentAge != 0 && resultDataMsg.PercentAge < 100 {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].CustomRequestTimeLine = resultDataMsg.PercentAge
			} else {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].CustomRequestTimeLine = 0
			}
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].ReceivedBytes += resultDataMsg.ReceivedBytes
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].SendBytes += resultDataMsg.SendBytes
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].TotalRequestNum += 1
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].TotalRequestTime += resultDataMsg.RequestTime
			if resultDataMsg.IsSucceed {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].SuccessNum += 1
			} else {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].ErrorNum += 1
			}
			requestTimeListMap[resultDataMsg.EventId] = append(requestTimeListMap[resultDataMsg.EventId], resultDataMsg.RequestTime)
		case <-ticker.C:
			if sceneTestResultDataMsg.Results == nil {
				continue
			}
			for eventId, requestTimeList := range requestTimeListMap {
				sort.Sort(requestTimeList)
				sceneTestResultDataMsg.Results[eventId].AvgRequestTime = sceneTestResultDataMsg.Results[eventId].TotalRequestTime / sceneTestResultDataMsg.Results[eventId].TotalRequestNum
				sceneTestResultDataMsg.Results[eventId].MaxRequestTime = requestTimeList[len(requestTimeList)-1]
				sceneTestResultDataMsg.Results[eventId].MinRequestTime = requestTimeList[0]
				sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLineValue = kao.TimeLineCalculate(sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLine, requestTimeList)
				sceneTestResultDataMsg.Results[eventId].NinetyRequestTimeLine = 90
				sceneTestResultDataMsg.Results[eventId].NinetyFiveRequestTimeLine = 95
				sceneTestResultDataMsg.Results[eventId].NinetyNineRequestTimeLine = 99
				sceneTestResultDataMsg.Results[eventId].NinetyRequestTimeLineValue = kao.TimeLineCalculate(90, requestTimeList)
				sceneTestResultDataMsg.Results[eventId].NinetyFiveRequestTimeLineValue = kao.TimeLineCalculate(95, requestTimeList)
				sceneTestResultDataMsg.Results[eventId].NinetyNineRequestTimeLineValue = kao.TimeLineCalculate(99, requestTimeList)
				total := float64(sceneTestResultDataMsg.Results[eventId].TotalRequestTime) / 1000
				if total <= 0 {
					total = 1
				}
				qps := fmt.Sprintf("%.2f", float64(sceneTestResultDataMsg.Results[eventId].TotalRequestNum)/total)
				sceneTestResultDataMsg.Results[eventId].Qps, _ = strconv.ParseFloat(qps, 64)

				//errorRate, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(sceneTestResultDataMsg.Results[eventId].ErrorNum)/float64(sceneTestResultDataMsg.Results[eventId].TotalRequestNum)), 64)
				//sceneTestResultDataMsg.Results[eventId].ErrorRate = errorRate
				//sceneTestResultDataMsg.Results[eventId].QpsList = append(sceneTestResultDataMsg.Results[eventId].QpsList, qps)
				//sceneTestResultDataMsg.Results[eventId].ErrorRateList = append(sceneTestResultDataMsg.Results[eventId].ErrorRateList, sceneTestResultDataMsg.Results[eventId].ErrorNum)
			}
			sceneTestResultDataMsg.TimeStamp = time.Now().Unix()
			err := es.InsertTestData(sceneTestResultDataMsg)
			if err != nil {
				log2.Logger.Error("es写入失败:", err)
				return
			}
		default:

		}
	}
}

func DeleteTopic(config *sarama.Config, host, topic string) (err error) {
	if config == nil || host == "" {
		return errors.New("config或host不能为空")
	}
	ca, err := sarama.NewClusterAdmin([]string{host}, config)
	if err != nil {
		log2.Logger.Error("删除topic时创建Admin失败：", err)
		return
	}
	if err = ca.DeleteTopic(topic); err != nil {
		log2.Logger.Error("删除topic"+topic+"是错误错误：", err)
		return
	}
	return
}
