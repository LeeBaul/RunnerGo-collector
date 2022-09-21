package server

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
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
	var machineMap = new(sync.Map)
	var requestTimeListMap = make(map[string]kao.RequestTimeList)
	var resultDataMsg = kao.ResultDataMsg{}
	var sceneTestResultDataMsg kao.SceneTestResultDataMsg
	var machineNum = int64(0)
	var index = 0
	ticker := time.NewTicker(2 * time.Second)
Loop:
	for {
		fmt.Println(111111)
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
				fmt.Println("resultDataMsg.End", resultDataMsg.End)
				machineNum = machineNum - 1
				if machineNum == 1 {
					sceneTestResultDataMsg.End = true
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
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].CustomRequestTimeLine = resultDataMsg.CustomRequestTimeLine
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].ReceivedBytes = +resultDataMsg.ReceivedBytes
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].SendBytes = +resultDataMsg.ReceivedBytes
			sceneTestResultDataMsg.Results[resultDataMsg.EventId].TotalRequestNum = +1
			if resultDataMsg.IsSucceed {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].SuccessNum = +1
			} else {
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].ErrorNum = +1
			}
			requestTimeListMap[resultDataMsg.EventId] = append(requestTimeListMap[resultDataMsg.ReportId], resultDataMsg.RequestTime)
		case <-ticker.C:
			if sceneTestResultDataMsg.Results == nil {
				continue
			}
			for eventId, requestTimeList := range requestTimeListMap {
				sort.Sort(requestTimeList)
				sceneTestResultDataMsg.Results[eventId].AvgRequestTime = sceneTestResultDataMsg.Results[eventId].TotalRequestTime / sceneTestResultDataMsg.Results[eventId].TotalRequestNum
				sceneTestResultDataMsg.Results[eventId].NinetyRequestTimeLine = kao.TimeLineCalculate(90, requestTimeList)
				sceneTestResultDataMsg.Results[eventId].NinetyFiveRequestTimeLine = kao.TimeLineCalculate(95, requestTimeList)
				sceneTestResultDataMsg.Results[eventId].NinetyNineRequestTimeLine = kao.TimeLineCalculate(99, requestTimeList)
				if sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLine != 0 {
					sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLineValue = kao.TimeLineCalculate(sceneTestResultDataMsg.Results[eventId].CustomRequestTimeLine, requestTimeList)
				}
				qps, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(len(requestTimeList))-float64(index)), 64)
				sceneTestResultDataMsg.Results[eventId].Qps = qps
			}
			sceneTestResultDataMsg.TimeStamp = time.Now().Unix()
			err := es.InsertTestData(sceneTestResultDataMsg, machineMap)
			if err != nil {
				log2.Logger.Error("es写入失败:", err)
				return
			}
			if sceneTestResultDataMsg.End {
				return
			}
		default:

		}
	}
}
