package main

import (
	"container/list"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"sync"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var gMutex sync.Mutex
var allMap, errMap map[string]int64
var reportJobs = list.New()

func initData() {
	allMap = make(map[string]int64, 0)
	errMap = make(map[string]int64, 0)
}

func lineProcess(line string) {
	//pre process
	logData, err := UnmarshallLog(line)
	if err != nil {
		return
	}
	//post process

	//Update Global Data
	gMutex.Lock()
	if logData.CompanyCode == "" {
		errMap[""] ++
	}

	if logData.Model == "ReportJob" {
		allMap[logData.CompanyCode]++
		reportJobs.PushBack(logData)
	}
	gMutex.Unlock()
}

func outPut() {
	s, _ := json.Marshal(allMap)
	fmt.Println(string(s))
	s, _ = json.Marshal(errMap)
	fmt.Println(string(s))
	fmt.Println("Logs ...")
	for e := reportJobs.Front(); e != nil; e = e.Next() {
		itemLog := LogSt(e.Value.(LogSt))
		if itemLog.Status == 0 {
			fmt.Printf("%v 商家任务完成于 %v, 耗时 %f \n", itemLog.CompanyCode, itemLog.RequestTime, itemLog.TimeConsume)
		}
	}

}

func UnmarshallLog(line string) (logData LogSt, err error) {
	//pre process
	err = json.Unmarshal([]byte(line), &logData)
	if err != nil {
		log.Fatalf("Log Content: %q process erro! %q \n", line, err)
		return
	}
	//post process
	return
}

type LogSt struct {
	CompanyCode string  `json:"companycode"`
	Model       string  `json:"model"`
	Level       string  `json:"level"`
	RequestTime time.Time  `json:"request_time"`
	Msg         string  `json:"msg"`
	ReportDate  string  `json:"more_info"`
	TimeConsume float64 `json:"response_time"`
	Status      int8    `json:"status"`
}
