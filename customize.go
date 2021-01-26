package main

import (
	"container/list"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var gMutex sync.Mutex
var allMap, errMap map[string]int64
var reportJobs = list.New()

func initData() {
	allMap = make(map[string]int64, 0)
	errMap = make(map[string]int64, 0)
}

func lineProcess(line string) {
	//pre process
	r, _ := regexp.Compile("ReportJob")
	if !r.MatchString(line) {
		return
	}

	logData, err := UnmarshallLog(line)
	if err != nil {
		return
	}
	//post process

	//Update Global Data
	gMutex.Lock()
	if logData.CompanyCode == "" {
		errMap[""]++
	}

	if logData.Model == "ReportJob" {
		allMap[logData.CompanyCode]++
		reportJobs.PushBack(logData)
	}
	gMutex.Unlock()
}

func outPut(csvFileName string) {
	fmt.Println("开始转换 CSV ...")
	file, err := os.Create(csvFileName)

	if err != nil {
		log.Fatalln("csv 文件创建出错", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if reportJobs.Len() > 0 {
		err = writer.Write([]string{"companyCode", "model", "level", "request_time", "msg", "report_date",
			"response_time", "status"})
		if err != nil {
			log.Fatalln("CSV output error", err)
		}
	}

	for e := reportJobs.Front(); e != nil; e = e.Next() {
		itemLog := e.Value.(LogSt)
		data := itemLog.Convert()
		err = writer.Write(data)
		if err != nil {
			log.Fatalln("CSV output error", err)
		}
	}
	fmt.Printf("CSV 转换完成 \n")
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
	CompanyCode string    `json:"companycode"`
	Model       string    `json:"model"`
	Level       string    `json:"level"`
	RequestTime time.Time `json:"request_time"`
	Msg         string    `json:"msg"`
	ReportDate  string    `json:"more_info"`
	TimeConsume int       `json:"response_time"`
	Status      int       `json:"status"`
}

func (log LogSt) Convert() []string {
	on := log.RequestTime.Format("2006-01-02 15:04:05")
	timeCon := strconv.Itoa(log.TimeConsume)
	status := strconv.Itoa(log.Status)
	return []string{log.CompanyCode, log.Model, log.Level, on, log.Msg, log.ReportDate, timeCon, status}
}
