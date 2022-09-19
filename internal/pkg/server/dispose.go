package server

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"kp-collector/internal/pkg/conf"
	"kp-collector/internal/pkg/dal/es"
	"kp-collector/internal/pkg/dal/kao"
	log2 "kp-collector/internal/pkg/log"
	"sort"
	"sync"
	"time"
)

func Execute(topic string, consumer sarama.Consumer) {
	if consumer == nil && topic == "" {
		return
	}
	var consumerMap = &sync.Map{}
	for {
		partitions, err := consumer.Partitions(topic)
		if err != nil {
			log2.Logger.Error("获取", conf.Conf.Kafka.Topic, "主题失败：", err)
			continue
		}
		for _, p := range partitions {
			if v, ok := consumerMap.Load(p); ok {
				continue
			} else {
				consumerMap.Store(v, 0)
			}
			partition, err := consumer.ConsumePartition(conf.Conf.Kafka.Topic, p, sarama.OffsetNewest)
			if err != nil {
				continue
			}
			go ReceiveMessage(partition, consumerMap, p)
			log2.Logger.Info("开始消费：", p, "分区消息")
		}
	}
}

func ReceiveMessage(partition sarama.PartitionConsumer, consumerMap *sync.Map, p int32) {
	defer partition.Close()
	defer consumerMap.Delete(p)
	if partition == nil || consumerMap == nil {
		return
	}
	var (
		sceneTestResultDataMsg = new(kao.SceneTestResultDataMsg)
		resultDataMsg          = kao.ResultDataMsg{}
		requestTimeMap         = make(map[string]kao.RequestTimeList)
		index                  = conf.Conf.ES.Index
		target                 = int64(0)
	)
	ticker := time.NewTicker(1 * time.Second)
Loop:
	for {
		select {
		case msg := <-partition.Messages():
			if msg.Topic != conf.Conf.Kafka.Topic {
				log2.Logger.Error("topic主题不一致")
				break Loop
			}

			err := json.Unmarshal(msg.Value, &resultDataMsg)
			if err != nil {
				log2.Logger.Error(err)
				break Loop
			}
			if resultDataMsg.MachineNum <= 0 {
				log2.Logger.Error("压力机数量为0")
				return
			}
			if resultDataMsg.ReportId == "" {
				log2.Logger.Error("报告Id为空")
				return
			}
			if target == 0 {
				target = resultDataMsg.MachineNum + 1
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
			if sceneTestResultDataMsg.SceneId == 0 {
				sceneTestResultDataMsg.SceneId = resultDataMsg.SceneId
			}
			if sceneTestResultDataMsg.SceneName == "" {
				sceneTestResultDataMsg.SceneName = resultDataMsg.SceneName
			}
			switch resultDataMsg.End {
			case true:
				target = target - 1
			case false:

				if sceneTestResultDataMsg.Results == nil {
					sceneTestResultDataMsg.Results = make(map[string]*kao.ApiTestResultDataMsg)
				}
				if sceneTestResultDataMsg.Results[resultDataMsg.EventId] == nil {
					sceneTestResultDataMsg.Results[resultDataMsg.EventId] = new(kao.ApiTestResultDataMsg)
				}
				if requestTimeMap[resultDataMsg.EventId] == nil {
					requestTimeMap[resultDataMsg.EventId] = kao.RequestTimeList{}
				}
				// 场景中各个接口响应时间
				requestTimeMap[resultDataMsg.EventId] = append(requestTimeMap[resultDataMsg.EventId], resultDataMsg.RequestTime)
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].ReceivedBytes += resultDataMsg.ReceivedBytes
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].SendBytes += resultDataMsg.SendBytes
				if resultDataMsg.IsSucceed {
					sceneTestResultDataMsg.Results[resultDataMsg.EventId].SuccessNum += 1
				} else {
					sceneTestResultDataMsg.Results[resultDataMsg.EventId].ErrorNum += 1
				}
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].TotalRequestNum += 1
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].TotalRequestTime += resultDataMsg.RequestTime
				sceneTestResultDataMsg.Results[resultDataMsg.EventId].CustomRequestTimeLineValue = resultDataMsg.CustomRequestTimeLine
			}

		case <-ticker.C:
			if sceneTestResultDataMsg.ReportId == "" {
				return
			}
			if target == 0 {
				return
			}
			if target == 1 {
				for k, v := range requestTimeMap {
					if v != nil {
						sort.Sort(v)
						sceneTestResultDataMsg.Results[k].MinRequestTime = v[0]
						sceneTestResultDataMsg.Results[k].MaxRequestTime = v[len(v)-1]
						sceneTestResultDataMsg.Results[k].AvgRequestTime = sceneTestResultDataMsg.Results[k].TotalRequestTime / sceneTestResultDataMsg.Results[k].TotalRequestNum
						sceneTestResultDataMsg.Results[k].NinetyRequestTimeLine = kao.TimeLineCalculate(90, v)
						sceneTestResultDataMsg.Results[k].NinetyFiveRequestTimeLine = kao.TimeLineCalculate(95, v)
						sceneTestResultDataMsg.Results[k].NinetyNineRequestTimeLine = kao.TimeLineCalculate(99, v)
						sceneTestResultDataMsg.Results[k].CustomRequestTimeLine = kao.TimeLineCalculate(sceneTestResultDataMsg.Results[k].CustomRequestTimeLineValue, v)
						sceneTestResultDataMsg.Results[k].Qps = float64(len(v)) - sceneTestResultDataMsg.Results[k].Qps
					}

				}
				sceneTestResultDataMsg.TimeStamp = time.Now().Unix()
				sceneTestResultDataMsg.End = true
				es.InsertTestData(sceneTestResultDataMsg, index)
				log2.Logger.Info(sceneTestResultDataMsg.ReportId, "报告结束")
				return
			}

			if target > 1 {
				for k, v := range requestTimeMap {
					if v != nil {
						sort.Sort(v)
						sceneTestResultDataMsg.Results[k].MinRequestTime = v[0]
						sceneTestResultDataMsg.Results[k].MaxRequestTime = v[len(v)-1]
						sceneTestResultDataMsg.Results[k].AvgRequestTime = sceneTestResultDataMsg.Results[k].TotalRequestTime / sceneTestResultDataMsg.Results[k].TotalRequestNum
						sceneTestResultDataMsg.Results[k].NinetyRequestTimeLine = kao.TimeLineCalculate(90, v)
						sceneTestResultDataMsg.Results[k].NinetyFiveRequestTimeLine = kao.TimeLineCalculate(95, v)
						sceneTestResultDataMsg.Results[k].NinetyNineRequestTimeLine = kao.TimeLineCalculate(99, v)
						sceneTestResultDataMsg.Results[k].CustomRequestTimeLine = kao.TimeLineCalculate(sceneTestResultDataMsg.Results[k].CustomRequestTimeLineValue, v)
						sceneTestResultDataMsg.Results[k].Qps = (float64(len(v)) - sceneTestResultDataMsg.Results[k].Qps) / 5
					}

				}
				sceneTestResultDataMsg.TimeStamp = time.Now().Unix()
				sceneTestResultDataMsg.End = false
				es.InsertTestData(sceneTestResultDataMsg, index)

			}

		}
	}
}
