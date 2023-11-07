package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	ProcessKey string = "process"
)
var dataCenter = struct {
	data    map[string]map[string]Data
	records map[string]record
	lock    sync.RWMutex
}{
	data: map[string]map[string]Data{
		ProcessKey: make(map[string]Data),
	},
	records: make(map[string]record),
}

var repeatTimes int

func ProcessMessageHandler(client mqtt.Client, msg mqtt.Message) {

	payload := &Payload[ProcessSlice]{}
	json.Unmarshal(msg.Payload(), &payload)

	if payload.IsSameMsg(ProcessKey) {
		repeatTimes += 1
		return
	}

	if payload.Timestamp < dataCenter.records[ProcessKey].Timestamp {
		// fmt.Printf("old timestamp:%d, new timestamp %s\n", dataCenter.records[ProcessKey].Mid, payload.Mid)
		fmt.Printf("\nold mid:%+v, new mid %+v\n", dataCenter.records[ProcessKey], payload)
	}

	if (payload.Data.Add == nil || len(payload.Data.Add) == 0) && (payload.Data.Del == nil || len(payload.Data.Del) == 0) && (payload.Data.Update == nil || len(payload.Data.Update) == 0) {
		fmt.Println("+++++++++++++++++++++++++++++++++++NULL DATA+++++++++++++++++++++++++++++++++++++++++")
		fmt.Printf("%s, %d, %s\n", payload.Mid, payload.Timestamp, time.Unix(payload.Timestamp/1000, 0).Format("2006-01-02 15:04:05"))
	}

	dataCenter.lock.Lock()
	fullData := dataCenter.data[ProcessKey]
	dataCenter.records[ProcessKey] = record{
		Mid:       payload.Mid,
		DeviceId:  payload.DeviceId,
		Timestamp: payload.Timestamp,
		Reply:     payload.Reply,

		Add:    payload.Data.Add.toDataSlice(),
		Del:    payload.Data.Del.toDataSlice(),
		Update: payload.Data.Update.toDataSlice(),
	}
	if payload.Data.Add != nil && len(payload.Data.Add) > 0 {
		for _, pro := range payload.Data.Add {
			_, ok := fullData[pro.ID()]
			if !ok {
				fullData[pro.ID()] = pro
			}
		}
	}

	if payload.Data.Del != nil && len(payload.Data.Del) > 0 {
		for _, pro := range payload.Data.Del {
			_, ok := fullData[pro.ID()]
			if ok {
				delete(fullData, pro.ID())
			}
		}
	}

	if payload.Data.Update != nil && len(payload.Data.Update) > 0 {
		for _, pro := range payload.Data.Update {
			p, ok := fullData[pro.ID()]
			if ok {
				fullData[pro.ID()] = p.Update(pro)
			}
		}
	}
	dataCenter.lock.Unlock()

	printProcessData()
	fmt.Printf("消息重复次数：%d", repeatTimes)
	repeatTimes = 0
	// fmt.Printf("收到主题：%s，消息：%:v\n", msg.Topic(), payload)
	// fmt.Printf("消息：%:v\n", recordes)
}
