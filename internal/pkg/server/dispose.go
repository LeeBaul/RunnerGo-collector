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
			log2.Logger.Error("0812037429374927492")
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
		sceneTestResultDataMsg     = new(kao.SceneTestResultDataMsg)
		tempSceneTestResultDataMsg = new(kao.SceneTestResultDataMsg)
		resultDataMsg              = kao.ResultDataMsg{}
		requestTimeMap             = make(map[interface{}]kao.RequestTimeList)
	)
	ticker := time.NewTicker(1 * time.Second)
	index := 0
	for {
		select {
		case msg := <-partition.Messages():
			if msg.Topic != conf.Conf.Kafka.Topic {
				continue
			}
			err := json.Unmarshal(msg.Value, &resultDataMsg)
			if err != nil {
				continue
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
			if sceneTestResultDataMsg.Results == nil {
				sceneTestResultDataMsg.Results = make(map[interface{}]*kao.ApiTestResultDataMsg)
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
		case <-ticker.C:
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
				}

			}
			s, _ := json.Marshal(sceneTestResultDataMsg)
			t, _ := json.Marshal(tempSceneTestResultDataMsg)
			if string(s) == string(t) {
				index++
				break
			}
			if index >= 2 {
				break
			}
			es.InsertTestData(sceneTestResultDataMsg, conf.Conf.ES.Index)
			index = 0
			log2.Logger.Info("sceneTestResultDataMsg", sceneTestResultDataMsg)
			tempSceneTestResultDataMsg = sceneTestResultDataMsg
		}

	}
}
